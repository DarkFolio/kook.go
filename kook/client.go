package kook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// BaseURL KOOK API 基础URL
	BaseURL = "https://www.kookapp.cn/api"
	// Version API版本
	Version = "v3"
	// UserAgent 用户代理
	UserAgent = "KOOK-Go-SDK/1.0.0"
)

// TokenType 鉴权类型
type TokenType string

const (
	// TokenTypeBot 机器人Token
	TokenTypeBot TokenType = "Bot"
	// TokenTypeBearer OAuth2 Token
	TokenTypeBearer TokenType = "Bearer"
)

// Client KOOK API客户端
type Client struct {
	httpClient  *http.Client
	token       string
	tokenType   TokenType
	baseURL     string
	logger      *logrus.Logger
	rateLimiter *GlobalRateLimiter
	retryConfig *RetryConfig

	// API服务
	User      *UserService
	Guild     *GuildService
	Channel   *ChannelService
	Message   *MessageService
	Gateway   *GatewayService
	Role      *RoleService
	Game      *GameService
	Friend    *FriendService
	Invite    *InviteService
	Asset     *AssetService
	Intimacy  *IntimacyService
	Badge     *BadgeService
	Blacklist *BlacklistService
	Emoji     *EmojiService
	Region    *RegionService
	OAuth     *OAuthService
	Live      *LiveService
	Admin     *AdminService
	Security  *SecurityService
	Voice     *VoiceService
	Item      *ItemService
	Order     *OrderService
	Coupon    *CouponService
	Boost     *BoostService
}

// ClientOption 客户端配置选项
type ClientOption func(*Client)

// WithHTTPClient 设置自定义HTTP客户端
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTokenType 设置Token类型
func WithTokenType(tokenType TokenType) ClientOption {
	return func(c *Client) {
		c.tokenType = tokenType
	}
}

// WithBaseURL 设置自定义基础URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithLogger 设置自定义日志器
func WithLogger(logger *logrus.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithRateLimiter 设置自定义速率限制器
func WithRateLimiter(rateLimiter *GlobalRateLimiter) ClientOption {
	return func(c *Client) {
		c.rateLimiter = rateLimiter
	}
}

// WithoutRateLimit 禁用速率限制
func WithoutRateLimit() ClientOption {
	return func(c *Client) {
		c.rateLimiter = nil
	}
}

// WithRetryConfig 设置自定义重试配置
func WithRetryConfig(config *RetryConfig) ClientOption {
	return func(c *Client) {
		c.retryConfig = config
	}
}

// WithoutRetry 禁用重试
func WithoutRetry() ClientOption {
	return func(c *Client) {
		c.retryConfig = &RetryConfig{MaxRetries: 0}
	}
}

// NewClient 创建新的KOOK客户端
func NewClient(token string, options ...ClientOption) *Client {
	if token == "" {
		panic("token不能为空")
	}

	// 默认HTTP客户端
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 默认日志器
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	client := &Client{
		httpClient:  httpClient,
		token:       token,
		tokenType:   TokenTypeBot,
		baseURL:     BaseURL,
		logger:      logger,
		rateLimiter: NewGlobalRateLimiter(),
		retryConfig: DefaultRetryConfig(),
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	// 初始化API服务
	client.User = &UserService{client: client}
	client.Guild = &GuildService{client: client}
	client.Channel = &ChannelService{client: client}
	client.Message = &MessageService{client: client}
	client.Gateway = &GatewayService{client: client}
	client.Role = &RoleService{client: client}
	client.Game = &GameService{client: client}
	client.Friend = &FriendService{client: client}
	client.Invite = &InviteService{client: client}
	client.Asset = &AssetService{client: client}
	client.Intimacy = &IntimacyService{client: client}
	client.Badge = &BadgeService{client: client}
	client.Blacklist = &BlacklistService{client: client}
	client.Emoji = &EmojiService{client: client}
	client.Region = &RegionService{client: client}
	client.OAuth = &OAuthService{client: client}
	client.Live = &LiveService{client: client}
	client.Admin = &AdminService{client: client}
	client.Security = &SecurityService{client: client}
	client.Voice = &VoiceService{client: client}
	client.Item = &ItemService{client: client}
	client.Order = &OrderService{client: client}
	client.Coupon = &CouponService{client: client}
	client.Boost = &BoostService{client: client}

	return client
}

// buildURL 构建完整的API URL
func (c *Client) buildURL(endpoint string) string {
	endpoint = strings.TrimPrefix(endpoint, "/")
	return fmt.Sprintf("%s/%s/%s", c.baseURL, Version, endpoint)
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(ctx context.Context, method, endpoint string, params map[string]interface{}, query map[string]string) (*Response, error) {
	// 使用重试机制执行请求
	return DoWithRetry(ctx, func(ctx context.Context) (*Response, error) {
		return c.doSingleRequest(ctx, method, endpoint, params, query)
	}, c.retryConfig, c.logger)
}

// doSingleRequest 执行单次HTTP请求
func (c *Client) doSingleRequest(ctx context.Context, method, endpoint string, params map[string]interface{}, query map[string]string) (*Response, error) {
	// 应用速率限制
	if c.rateLimiter != nil {
		c.rateLimiter.Wait(endpoint)
	}

	requestURL := c.buildURL(endpoint)

	// 添加查询参数
	if len(query) > 0 {
		u, err := url.Parse(requestURL)
		if err != nil {
			return nil, fmt.Errorf("解析URL失败: %w", err)
		}

		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		requestURL = u.String()
	}

	var body io.Reader
	if params != nil {
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("序列化请求参数失败: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
		c.logger.WithField("params", string(jsonData)).Debugf("请求参数")
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.tokenType, c.token))
	req.Header.Set("User-Agent", UserAgent)
	if method == "POST" && params != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept-Language", "zh-cn")

	c.logger.WithFields(logrus.Fields{
		"method":  method,
		"url":     requestURL,
		"headers": req.Header,
	}).Debugf("发送API请求")

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.WithError(err).Errorf("请求失败")
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithError(err).Errorf("读取响应失败")
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"status": resp.StatusCode,
		"body":   string(respBody),
	}).Debugf("收到API响应")

	// 解析响应
	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		c.logger.WithError(err).Errorf("解析响应失败")
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查API错误
	if response.Code != 0 {
		err := NewKOOKError(response.Code, response.Message).
			WithContext(method, endpoint)

		// 从响应头中提取请求ID
		if requestID := resp.Header.Get("X-Request-Id"); requestID != "" {
			err = err.WithRequestID(requestID)
		}

		// 从响应头中提取重试延迟
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			if seconds, parseErr := time.ParseDuration(retryAfter + "s"); parseErr == nil {
				err = err.WithRetryAfter(seconds)
			}
		}

		err.HTTPStatus = resp.StatusCode

		c.logger.WithError(err).Errorf("API返回错误")
		return &response, err
	}

	c.logger.Infof("API请求成功: %s %s", method, requestURL)
	return &response, nil
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, endpoint string, query map[string]string) (*Response, error) {
	return c.doRequest(ctx, "GET", endpoint, nil, query)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, endpoint string, params map[string]interface{}) (*Response, error) {
	return c.doRequest(ctx, "POST", endpoint, params, nil)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, endpoint string, params map[string]interface{}) (*Response, error) {
	return c.doRequest(ctx, "PUT", endpoint, params, nil)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, endpoint string, params map[string]interface{}) (*Response, error) {
	return c.doRequest(ctx, "DELETE", endpoint, params, nil)
}

// Response API响应结构
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}
