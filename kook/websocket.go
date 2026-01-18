package kook

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// EventHandler 事件处理器函数类型
type EventHandler func(*Event)

// WebSocketClient WebSocket客户端
type WebSocketClient struct {
	client          *Client
	conn            *websocket.Conn
	eventHandlers   map[int][]EventHandler
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	compress        bool
	sn              int
	sessionID       string
	heartbeatTicker *time.Ticker
	gatewayURL      string
	reconnectCount  int
	maxReconnects   int
	reconnectDelay  time.Duration
	isConnected     bool
	connMu          sync.RWMutex
}

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	S  int             `json:"s"`  // 信令类型
	D  json.RawMessage `json:"d"`  // 数据
	SN int             `json:"sn"` // 序号
}

// HelloMessage Hello消息
type HelloMessage struct {
	Code      int    `json:"code"`
	SessionID string `json:"session_id"`
}

// PingMessage Ping消息
type PingMessage struct {
	SN int `json:"sn"`
}

// PongMessage Pong消息
type PongMessage struct {
	SN int `json:"sn"`
}

// ResumeMessage Resume消息
type ResumeMessage struct {
	SessionID string `json:"session_id"`
	SN        int    `json:"sn"`
}

// 信令类型常量
const (
	SignalEvent     = 0 // 事件
	SignalHello     = 1 // 服务端发送，客户端接收，代表连接成功
	SignalPing      = 2 // 双向：服务端ping客户端，客户端也可以ping服务端
	SignalPong      = 3 // 双向：ping的响应
	SignalResume    = 4 // 客户端发送，服务端接收，代表重连
	SignalReconnect = 5 // 服务端发送，客户端接收，代表需要重连
	SignalResumeAck = 6 // 服务端发送，客户端接收，代表重连成功
)

// NewWebSocketClient 创建新的WebSocket客户端
func NewWebSocketClient(client *Client, compress bool) *WebSocketClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &WebSocketClient{
		client:         client,
		eventHandlers:  make(map[int][]EventHandler),
		ctx:            ctx,
		cancel:         cancel,
		compress:       compress,
		maxReconnects:  10,
		reconnectDelay: 5 * time.Second,
	}
}

// OnEvent 注册事件处理器
func (ws *WebSocketClient) OnEvent(eventType int, handler EventHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.eventHandlers[eventType] = append(ws.eventHandlers[eventType], handler)
}

// Connect 连接到WebSocket网关
func (ws *WebSocketClient) Connect() error {
	return ws.connectWithRetry()
}

// connectWithRetry 带重试的连接
func (ws *WebSocketClient) connectWithRetry() error {
	for attempts := 0; attempts <= ws.maxReconnects; attempts++ {
		err := ws.doConnect()
		if err == nil {
			ws.reconnectCount = 0
			return nil
		}

		ws.client.logger.WithError(err).Errorf("WebSocket连接失败，尝试 %d/%d", attempts+1, ws.maxReconnects+1)

		if attempts < ws.maxReconnects {
			select {
			case <-ws.ctx.Done():
				return ws.ctx.Err()
			case <-time.After(ws.reconnectDelay * time.Duration(attempts+1)):
				// 指数退避
			}
		}
	}

	return fmt.Errorf("WebSocket连接失败，已达到最大重试次数")
}

// doConnect 执行实际连接
func (ws *WebSocketClient) doConnect() error {
	// 获取网关信息
	compress := 0
	if ws.compress {
		compress = 1
	}

	gateway, err := ws.client.Gateway.GetGateway(compress)
	if err != nil {
		return fmt.Errorf("获取网关信息失败: %w", err)
	}

	ws.gatewayURL = gateway.URL

	// 创建WebSocket连接
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("%s %s", ws.client.tokenType, ws.client.token))

	ws.client.logger.Infof("连接到WebSocket网关: %s", gateway.URL)

	conn, _, err := websocket.DefaultDialer.Dial(gateway.URL, header)
	if err != nil {
		return fmt.Errorf("WebSocket连接失败: %w", err)
	}

	ws.connMu.Lock()
	ws.conn = conn
	ws.isConnected = true
	ws.connMu.Unlock()

	ws.client.logger.Info("WebSocket连接成功")

	// 启动消息处理协程
	go ws.handleMessages()

	return nil
}

// Close 关闭WebSocket连接
func (ws *WebSocketClient) Close() error {
	ws.cancel()

	if ws.heartbeatTicker != nil {
		ws.heartbeatTicker.Stop()
	}

	if ws.conn != nil {
		return ws.conn.Close()
	}

	return nil
}

// handleMessages 处理WebSocket消息
func (ws *WebSocketClient) handleMessages() {
	defer func() {
		if r := recover(); r != nil {
			ws.client.logger.Errorf("WebSocket消息处理发生panic: %v", r)
		}

		// 标记连接已断开
		ws.connMu.Lock()
		ws.isConnected = false
		ws.connMu.Unlock()

		// 尝试重连
		ws.attemptReconnect()
	}()

	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			ws.connMu.RLock()
			conn := ws.conn
			ws.connMu.RUnlock()

			if conn == nil {
				ws.client.logger.Error("WebSocket连接为空")
				return
			}

			_, data, err := conn.ReadMessage()
			if err != nil {
				ws.client.logger.WithError(err).Error("读取WebSocket消息失败")
				return
			}

			// 如果启用了压缩，需要解压
			if ws.compress {
				data, err = ws.decompress(data)
				if err != nil {
					ws.client.logger.WithError(err).Error("解压消息失败")
					continue
				}
			}

			ws.client.logger.Debugf("收到WebSocket消息: %s", string(data))

			var msg WebSocketMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				ws.client.logger.WithError(err).Error("解析WebSocket消息失败")
				continue
			}

			if err := ws.handleMessage(&msg); err != nil {
				ws.client.logger.WithError(err).Error("处理WebSocket消息失败")
			}
		}
	}
}

// attemptReconnect 尝试重连
func (ws *WebSocketClient) attemptReconnect() {
	if ws.reconnectCount >= ws.maxReconnects {
		ws.client.logger.Error("已达到最大重连次数，停止重连")
		return
	}

	ws.reconnectCount++
	ws.client.logger.Infof("开始第 %d 次重连尝试", ws.reconnectCount)

	// 等待一段时间后重连
	time.Sleep(ws.reconnectDelay * time.Duration(ws.reconnectCount))

	err := ws.doConnect()
	if err != nil {
		ws.client.logger.WithError(err).Errorf("重连失败")
		// 递归尝试重连
		go ws.attemptReconnect()
	} else {
		ws.client.logger.Info("重连成功")
		ws.reconnectCount = 0
	}
}

// IsConnected 检查连接状态
func (ws *WebSocketClient) IsConnected() bool {
	ws.connMu.RLock()
	defer ws.connMu.RUnlock()
	return ws.isConnected
}

// handleMessage 处理单个WebSocket消息
func (ws *WebSocketClient) handleMessage(msg *WebSocketMessage) error {
	switch msg.S {
	case SignalEvent:
		// 处理事件
		return ws.handleEvent(msg)
	case SignalHello:
		// 处理Hello消息
		return ws.handleHello(msg)
	case SignalPing:
		// 处理Ping消息
		return ws.handlePing(msg)
	case SignalReconnect:
		// 处理重连消息
		return ws.handleReconnect(msg)
	case SignalResumeAck:
		// 处理重连确认消息
		return ws.handleResumeAck(msg)
	case SignalPong:
		// 处理Pong消息
		var pong PongMessage
		if msg.D != nil {
			if err := json.Unmarshal(msg.D, &pong); err != nil {
				ws.client.logger.WithError(err).Debug("解析Pong消息失败，可能是空的Pong")
			} else {
				ws.client.logger.Debugf("收到Pong响应，SN: %d", pong.SN)
			}
		} else {
			ws.client.logger.Debug("收到Pong响应")
		}
		return nil
	default:
		ws.client.logger.Warnf("收到未知信令类型: %d", msg.S)
	}

	return nil
}

// handleEvent 处理事件消息
func (ws *WebSocketClient) handleEvent(msg *WebSocketMessage) error {
	var event Event
	if err := json.Unmarshal(msg.D, &event); err != nil {
		return fmt.Errorf("解析事件失败: %w", err)
	}

	ws.sn = msg.SN
	ws.client.logger.Debugf("收到事件: 类型=%d, 内容=%s", event.Type, event.Content)

	// 调用事件处理器
	ws.mu.RLock()
	handlers := ws.eventHandlers[event.Type]
	ws.mu.RUnlock()

	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					ws.client.logger.Errorf("事件处理器发生panic: %v", r)
				}
			}()
			h(&event)
		}(handler)
	}

	return nil
}

// handleHello 处理Hello消息
func (ws *WebSocketClient) handleHello(msg *WebSocketMessage) error {
	var hello HelloMessage
	if err := json.Unmarshal(msg.D, &hello); err != nil {
		return fmt.Errorf("解析Hello消息失败: %w", err)
	}

	ws.sessionID = hello.SessionID
	ws.client.logger.Infof("WebSocket会话建立成功: %s", ws.sessionID)

	// 启动心跳
	ws.startHeartbeat()

	return nil
}

// handlePing 处理Ping消息
func (ws *WebSocketClient) handlePing(msg *WebSocketMessage) error {
	var ping PingMessage
	if err := json.Unmarshal(msg.D, &ping); err != nil {
		return fmt.Errorf("解析Ping消息失败: %w", err)
	}

	// 发送Pong响应
	pong := WebSocketMessage{
		S: SignalPong,
		D: nil,
	}

	pongData, _ := json.Marshal(PongMessage{SN: ping.SN})
	pong.D = pongData

	return ws.sendMessage(&pong)
}

// handleReconnect 处理重连消息
func (ws *WebSocketClient) handleReconnect(msg *WebSocketMessage) error {
	ws.client.logger.Warn("服务器要求重连")

	// 发送Resume消息
	resume := WebSocketMessage{
		S: SignalResume,
	}

	resumeData, _ := json.Marshal(ResumeMessage{
		SessionID: ws.sessionID,
		SN:        ws.sn,
	})
	resume.D = resumeData

	return ws.sendMessage(&resume)
}

// handleResumeAck 处理重连确认消息
func (ws *WebSocketClient) handleResumeAck(msg *WebSocketMessage) error {
	ws.client.logger.Info("重连成功")
	return nil
}

// startHeartbeat 启动心跳
func (ws *WebSocketClient) startHeartbeat() {
	// 每30秒发送一次心跳
	ws.heartbeatTicker = time.NewTicker(30 * time.Second)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ws.client.logger.Errorf("心跳处理发生panic: %v", r)
			}
		}()

		consecutiveFailures := 0
		const maxFailures = 3

		for {
			select {
			case <-ws.ctx.Done():
				return
			case <-ws.heartbeatTicker.C:
				ping := WebSocketMessage{
					S: SignalPing,
				}

				pingData, _ := json.Marshal(PingMessage{SN: ws.sn})
				ping.D = pingData

				if err := ws.sendMessage(&ping); err != nil {
					consecutiveFailures++
					ws.client.logger.WithError(err).Errorf("发送心跳失败 (%d/%d)", consecutiveFailures, maxFailures)

					if consecutiveFailures >= maxFailures {
						ws.client.logger.Error("连续心跳失败，触发重连")
						ws.connMu.Lock()
						ws.isConnected = false
						if ws.conn != nil {
							ws.conn.Close()
						}
						ws.connMu.Unlock()
						go ws.attemptReconnect()
						return
					}
				} else {
					if consecutiveFailures > 0 {
						ws.client.logger.Info("心跳恢复正常")
					}
					consecutiveFailures = 0
				}
			}
		}
	}()
}

// sendMessage 发送WebSocket消息
func (ws *WebSocketClient) sendMessage(msg *WebSocketMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	ws.client.logger.Debugf("发送WebSocket消息: %s", string(data))

	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

// decompress 解压数据
func (ws *WebSocketClient) decompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}
