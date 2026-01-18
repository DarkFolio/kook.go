package kook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ErrorCode KOOK API 错误代码
type ErrorCode int

// 常见错误代码
const (
	ErrorCodeOK                  ErrorCode = 0
	ErrorCodeBadRequest          ErrorCode = 40000
	ErrorCodeUnauthorized        ErrorCode = 40100
	ErrorCodeForbidden           ErrorCode = 40300
	ErrorCodeNotFound            ErrorCode = 40400
	ErrorCodeMethodNotAllowed    ErrorCode = 40500
	ErrorCodeTooManyRequests     ErrorCode = 42900
	ErrorCodeInternalServerError ErrorCode = 50000
	ErrorCodeBadGateway          ErrorCode = 50200
	ErrorCodeServiceUnavailable  ErrorCode = 50300
	ErrorCodeGatewayTimeout      ErrorCode = 50400
)

// ErrorCodeMap 错误代码映射
var ErrorCodeMap = map[ErrorCode]string{
	ErrorCodeOK:                  "请求成功",
	ErrorCodeBadRequest:          "请求参数错误",
	ErrorCodeUnauthorized:        "认证失败，Token无效",
	ErrorCodeForbidden:           "权限不足",
	ErrorCodeNotFound:            "资源不存在",
	ErrorCodeMethodNotAllowed:    "请求方法不允许",
	ErrorCodeTooManyRequests:     "请求过于频繁",
	ErrorCodeInternalServerError: "服务器内部错误",
	ErrorCodeBadGateway:          "网关错误",
	ErrorCodeServiceUnavailable:  "服务不可用",
	ErrorCodeGatewayTimeout:      "网关超时",
}

// KOOKError KOOK API 错误
type KOOKError struct {
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	RequestID  string        `json:"request_id,omitempty"`
	Timestamp  time.Time     `json:"timestamp"`
	HTTPStatus int           `json:"http_status,omitempty"`
	Endpoint   string        `json:"endpoint,omitempty"`
	Method     string        `json:"method,omitempty"`
	RetryAfter time.Duration `json:"retry_after,omitempty"`
	Details    interface{}   `json:"details,omitempty"`
}

// Error 实现 error 接口
func (e *KOOKError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("KOOK API错误 [%d]: %s", e.Code, e.Message)
	}

	if desc, exists := ErrorCodeMap[ErrorCode(e.Code)]; exists {
		return fmt.Sprintf("KOOK API错误 [%d]: %s", e.Code, desc)
	}

	return fmt.Sprintf("KOOK API错误 [%d]: 未知错误", e.Code)
}

// IsRetryable 判断错误是否可重试
func (e *KOOKError) IsRetryable() bool {
	// 5xx 服务器错误通常可重试
	if e.Code >= 500 && e.Code < 600 {
		return true
	}

	// 429 速率限制错误可重试
	if e.Code == 429 {
		return true
	}

	// 网络超时相关错误可重试
	if e.HTTPStatus == http.StatusRequestTimeout ||
		e.HTTPStatus == http.StatusTooManyRequests ||
		e.HTTPStatus == http.StatusBadGateway ||
		e.HTTPStatus == http.StatusServiceUnavailable ||
		e.HTTPStatus == http.StatusGatewayTimeout {
		return true
	}

	return false
}

// IsRateLimited 判断是否为速率限制错误
func (e *KOOKError) IsRateLimited() bool {
	return e.Code == 429 || e.HTTPStatus == http.StatusTooManyRequests
}

// IsAuthError 判断是否为认证错误
func (e *KOOKError) IsAuthError() bool {
	return e.Code == 40100 || e.HTTPStatus == http.StatusUnauthorized
}

// IsPermissionError 判断是否为权限错误
func (e *KOOKError) IsPermissionError() bool {
	return e.Code == 40300 || e.HTTPStatus == http.StatusForbidden
}

// IsNotFoundError 判断是否为资源不存在错误
func (e *KOOKError) IsNotFoundError() bool {
	return e.Code == 40400 || e.HTTPStatus == http.StatusNotFound
}

// IsServerError 判断是否为服务器错误
func (e *KOOKError) IsServerError() bool {
	return e.Code >= 500 || e.HTTPStatus >= 500
}

// WithContext 添加错误上下文
func (e *KOOKError) WithContext(method, endpoint string) *KOOKError {
	newErr := *e
	newErr.Method = method
	newErr.Endpoint = endpoint
	newErr.Timestamp = time.Now()
	return &newErr
}

// WithRequestID 添加请求ID
func (e *KOOKError) WithRequestID(requestID string) *KOOKError {
	newErr := *e
	newErr.RequestID = requestID
	return &newErr
}

// WithRetryAfter 添加重试延迟信息
func (e *KOOKError) WithRetryAfter(retryAfter time.Duration) *KOOKError {
	newErr := *e
	newErr.RetryAfter = retryAfter
	return &newErr
}

// WithDetails 添加错误详情
func (e *KOOKError) WithDetails(details interface{}) *KOOKError {
	newErr := *e
	newErr.Details = details
	return &newErr
}

// ValidationError 参数验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error 实现 error 接口
func (e *ValidationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("参数验证失败 [%s]: %s (值: %s)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("参数验证失败 [%s]: %s", e.Field, e.Message)
}

// NewKOOKError 创建 KOOK 错误
func NewKOOKError(code int, message string) *KOOKError {
	return &KOOKError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewKOOKErrorFromResponse 从 HTTP 响应创建 KOOK 错误
func NewKOOKErrorFromResponse(resp *http.Response, body []byte) *KOOKError {
	err := &KOOKError{
		HTTPStatus: resp.StatusCode,
		Timestamp:  time.Now(),
	}

	// 尝试解析 JSON 错误响应
	var apiResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if json.Unmarshal(body, &apiResp) == nil && apiResp.Code != 0 {
		err.Code = apiResp.Code
		err.Message = apiResp.Message
	} else {
		// 使用 HTTP 状态码
		err.Code = resp.StatusCode
		err.Message = http.StatusText(resp.StatusCode)
	}

	// 提取请求ID
	if requestID := resp.Header.Get("X-Request-Id"); requestID != "" {
		err.RequestID = requestID
	}

	// 提取重试延迟
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
			err.RetryAfter = time.Duration(seconds) * time.Second
		}
	}

	return err
}

// NewValidationError 创建参数验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewValidationErrorWithValue 创建带值的参数验证错误
func NewValidationErrorWithValue(field, message, value string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// IsKOOKError 检查是否为 KOOK 错误
func IsKOOKError(err error) (*KOOKError, bool) {
	if kookErr, ok := err.(*KOOKError); ok {
		return kookErr, true
	}
	return nil, false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) (*ValidationError, bool) {
	if validationErr, ok := err.(*ValidationError); ok {
		return validationErr, true
	}
	return nil, false
}

// 为了保持向后兼容，保留原有的 APIError 类型
type APIError = KOOKError

// IsAPIError 检查是否为 API 错误（向后兼容）
func IsAPIError(err error) (*APIError, bool) {
	if kookErr, ok := IsKOOKError(err); ok {
		// KOOKError 和 APIError 是同一个类型（别名），所以可以直接返回
		return kookErr, true
	}
	return nil, false
}
