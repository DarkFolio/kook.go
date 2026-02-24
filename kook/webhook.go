package kook

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// WebhookHandler Webhook处理器
type WebhookHandler struct {
	client        *Client
	encryptKey    string
	verifyToken   string
	eventHandlers map[int][]EventHandler
	mu            sync.RWMutex
}

// WebhookMessage Webhook消息结构
type WebhookMessage struct {
	S  int             `json:"s"`  // 信令类型
	D  json.RawMessage `json:"d"`  // 数据
	SN int             `json:"sn"` // 序号
}

type encryptedWebhookMessage struct {
	Encrypt string `json:"encrypt"`
}

type webhookPayloadMeta struct {
	ChannelType string `json:"channel_type"`
	VerifyToken string `json:"verify_token"`
	Challenge   string `json:"challenge"`
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
	wh.mu.Lock()
	defer wh.mu.Unlock()
	wh.eventHandlers[eventType] = append(wh.eventHandlers[eventType], handler)
}

// HandleRequest 处理HTTP请求
func (wh *WebhookHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		wh.client.logger.WithError(err).Error("读取请求体失败")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	body, err = decodeRequestBody(body, r.Header.Get("Content-Encoding"))
	if err != nil {
		wh.client.logger.WithError(err).Error("解码Webhook请求体失败")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	body, err = wh.tryDecryptBody(body)
	if err != nil {
		wh.client.logger.WithError(err).Error("解密Webhook请求体失败")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	wh.client.logger.Debugf("收到Webhook消息: %s", string(body))

	var msg WebhookMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		wh.client.logger.WithError(err).Error("解析Webhook消息失败")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	challenge, err := wh.handleMessage(&msg)
	if err != nil {
		wh.client.logger.WithError(err).Error("处理Webhook消息失败")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if challenge != "" {
		_ = json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
		return
	}
	_, _ = w.Write([]byte(`{"code":0}`))
}

func decodeRequestBody(body []byte, encoding string) ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "", "identity":
		return body, nil
	case "gzip":
		r, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return io.ReadAll(r)
	case "deflate":
		r, err := zlib.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return io.ReadAll(r)
	default:
		return nil, fmt.Errorf("不支持的Content-Encoding: %s", encoding)
	}
}

func (wh *WebhookHandler) tryDecryptBody(body []byte) ([]byte, error) {
	var encrypted encryptedWebhookMessage
	if err := json.Unmarshal(body, &encrypted); err != nil {
		return body, nil
	}
	if encrypted.Encrypt == "" {
		return body, nil
	}
	return decryptWebhookPayload(encrypted.Encrypt, wh.encryptKey)
}

func decryptWebhookPayload(encrypted, encryptKey string) ([]byte, error) {
	if encryptKey == "" {
		return nil, fmt.Errorf("Webhook消息已加密但未配置encryptKey")
	}

	keyBytes := []byte(encryptKey)
	if len(keyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		keyBytes = padded
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	payload, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("解码encrypt失败: %w", err)
	}
	if len(payload) <= aes.BlockSize {
		return nil, fmt.Errorf("encrypt内容长度异常")
	}

	iv := payload[:aes.BlockSize]
	cipherBase64 := payload[aes.BlockSize:]
	cipherText, err := base64.StdEncoding.DecodeString(string(cipherBase64))
	if err != nil {
		return nil, fmt.Errorf("解码cipher失败: %w", err)
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("cipher长度不是BlockSize整数倍")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plainText, cipherText)

	plainText, err = pkcs7Unpad(plainText, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("无效的PKCS7数据")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > blockSize || pad > len(data) {
		return nil, fmt.Errorf("无效的PKCS7填充")
	}
	for i := len(data) - pad; i < len(data); i++ {
		if int(data[i]) != pad {
			return nil, fmt.Errorf("无效的PKCS7填充字节")
		}
	}
	return data[:len(data)-pad], nil
}

// handleMessage 处理Webhook消息
func (wh *WebhookHandler) handleMessage(msg *WebhookMessage) (string, error) {
	if msg.S != SignalEvent {
		return "", nil
	}

	var meta webhookPayloadMeta
	if err := json.Unmarshal(msg.D, &meta); err != nil {
		return "", fmt.Errorf("解析Webhook元数据失败: %w", err)
	}

	if wh.verifyToken != "" && meta.VerifyToken != wh.verifyToken {
		return "", fmt.Errorf("Webhook verify_token 不匹配")
	}

	if strings.EqualFold(meta.ChannelType, "WEBHOOK_CHALLENGE") && meta.Challenge != "" {
		wh.client.logger.Info("收到Webhook验证挑战")
		return meta.Challenge, nil
	}

	return "", wh.handleEvent(msg)
}

// handleEvent 处理事件
func (wh *WebhookHandler) handleEvent(msg *WebhookMessage) error {
	var event Event
	if err := json.Unmarshal(msg.D, &event); err != nil {
		return fmt.Errorf("解析事件失败: %w", err)
	}

	wh.client.logger.Debugf("收到Webhook事件: 类型=%d, 内容=%s", event.Type, event.Content)

	wh.mu.RLock()
	handlers := append([]EventHandler(nil), wh.eventHandlers[event.Type]...)
	wh.mu.RUnlock()

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

// StartWebhookServer 启动Webhook服务器
func (wh *WebhookHandler) StartWebhookServer(addr, path string) error {
	http.HandleFunc(path, wh.HandleRequest)

	wh.client.logger.Infof("启动Webhook服务器: %s%s", addr, path)
	return http.ListenAndServe(addr, nil)
}
