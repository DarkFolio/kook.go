package kook

import (
	"compress/zlib"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// WebhookHandler Webhook处理器
type WebhookHandler struct {
	client       *Client
	encryptKey   string
	verifyToken  string
	eventHandlers map[int][]EventHandler
}

// WebhookMessage Webhook消息结构
type WebhookMessage struct {
	S         int             `json:"s"`          // 信令类型
	D         json.RawMessage `json:"d"`          // 数据
	SN        int             `json:"sn"`         // 序号
	Challenge string          `json:"challenge"`  // 验证挑战
}

// NewWebhookHandler 创建新的Webhook处理器
func NewWebhookHandler(client *Client, encryptKey, verifyToken string) *WebhookHandler {
	return &WebhookHandler{
		client:        client,
		encryptKey:    encryptKey,
		verifyToken:   verifyToken,
		eventHandlers: make(map[int][]EventHandler),
	}
}

// OnEvent 注册事件处理器
func (wh *WebhookHandler) OnEvent(eventType int, handler EventHandler) {
	wh.eventHandlers[eventType] = append(wh.eventHandlers[eventType], handler)
}

// HandleRequest 处理HTTP请求
func (wh *WebhookHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 验证请求方法
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		wh.client.logger.WithError(err).Error("读取请求体失败")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 验证签名
	if !wh.verifySignature(r, body) {
		wh.client.logger.Error("Webhook签名验证失败")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 解压数据（如果需要）
	if r.Header.Get("Content-Encoding") == "gzip" || r.Header.Get("Content-Encoding") == "deflate" {
		body, err = wh.decompress(body)
		if err != nil {
			wh.client.logger.WithError(err).Error("解压数据失败")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	wh.client.logger.Debugf("收到Webhook消息: %s", string(body))

	// 解析消息
	var msg WebhookMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		wh.client.logger.WithError(err).Error("解析Webhook消息失败")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 处理消息
	if err := wh.handleMessage(&msg); err != nil {
		wh.client.logger.WithError(err).Error("处理Webhook消息失败")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// 如果是挑战请求，返回challenge
	if msg.Challenge != "" {
		response := map[string]string{"challenge": msg.Challenge}
		json.NewEncoder(w).Encode(response)
	} else {
		w.Write([]byte(`{"code": 0}`))
	}
}

// verifySignature 验证签名
func (wh *WebhookHandler) verifySignature(r *http.Request, body []byte) bool {
	if wh.verifyToken == "" {
		return true // 如果没有设置验证token，跳过验证
	}

	// 获取时间戳和签名
	timestamp := r.Header.Get("X-Kook-Request-Timestamp")
	nonce := r.Header.Get("X-Kook-Request-Nonce")
	signature := r.Header.Get("X-Kook-Signature")

	if timestamp == "" || nonce == "" || signature == "" {
		return false
	}

	// 验证时间戳（5分钟内有效）
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	
	if time.Now().Unix()-ts > 300 {
		wh.client.logger.Warn("Webhook时间戳过期")
		return false
	}

	// 计算签名
	data := wh.verifyToken + timestamp + nonce + string(body)
	h := hmac.New(sha256.New, []byte(wh.verifyToken))
	h.Write([]byte(data))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 比较签名
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// handleMessage 处理Webhook消息
func (wh *WebhookHandler) handleMessage(msg *WebhookMessage) error {
	// 如果是挑战请求，直接返回
	if msg.Challenge != "" {
		wh.client.logger.Info("收到Webhook验证挑战")
		return nil
	}

	// 处理事件
	if msg.S == SignalEvent {
		return wh.handleEvent(msg)
	}

	return nil
}

// handleEvent 处理事件
func (wh *WebhookHandler) handleEvent(msg *WebhookMessage) error {
	var event Event
	if err := json.Unmarshal(msg.D, &event); err != nil {
		return fmt.Errorf("解析事件失败: %w", err)
	}

	wh.client.logger.Debugf("收到Webhook事件: 类型=%d, 内容=%s", event.Type, event.Content)

	// 调用事件处理器
	handlers := wh.eventHandlers[event.Type]
	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					wh.client.logger.Errorf("事件处理器发生panic: %v", r)
				}
			}()
			h(&event)
		}(handler)
	}

	return nil
}

// decompress 解压数据
func (wh *WebhookHandler) decompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

// StartWebhookServer 启动Webhook服务器
func (wh *WebhookHandler) StartWebhookServer(addr, path string) error {
	http.HandleFunc(path, wh.HandleRequest)
	
	wh.client.logger.Infof("启动Webhook服务器: %s%s", addr, path)
	return http.ListenAndServe(addr, nil)
} 