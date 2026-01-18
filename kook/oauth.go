package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// OAuthService OAuth相关API服务
type OAuthService struct {
	client *Client
}

// GetOAuthToken 获取OAuth Token
func (s *OAuthService) GetOAuthToken(ctx context.Context, grantType, clientID, clientSecret, code, redirectURI string) (*OAuthTokenResponse, error) {
	if grantType == "" {
		return nil, fmt.Errorf("授权类型不能为空")
	}
	if clientID == "" {
		return nil, fmt.Errorf("客户端ID不能为空")
	}

	params := map[string]interface{}{
		"grant_type":    grantType,
		"client_id":     clientID,
		"client_secret": clientSecret,
	}

	if code != "" {
		params["code"] = code
	}
	if redirectURI != "" {
		params["redirect_uri"] = redirectURI
	}

	resp, err := s.client.Post(ctx, "oauth2/token", params)
	if err != nil {
		return nil, err
	}

	var result OAuthTokenResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析OAuth Token失败: %w", err)
	}

	return &result, nil
}

// 数据结构定义

// OAuthTokenResponse OAuth Token响应
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	TokenType    string `json:"token_type"`    // 令牌类型
	ExpiresIn    int    `json:"expires_in"`    // 过期时间（秒）
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	Scope        string `json:"scope"`         // 权限范围
} 