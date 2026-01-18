package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// UserService 用户相关API服务
type UserService struct {
	client *Client
}

// GetMe 获取当前用户信息
func (s *UserService) GetMe(ctx context.Context) (*User, error) {
	resp, err := s.client.Get(ctx, "user/me", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return &user, nil
}

// GetUser 获取指定用户信息
func (s *UserService) GetUser(ctx context.Context, userID string, guildID string) (*User, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	query := map[string]string{
		"user_id": userID,
	}

	if guildID != "" {
		query["guild_id"] = guildID
	}

	resp, err := s.client.Get(ctx, "user/view", query)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return &user, nil
}

// GetUserOnlineStatus 获取用户在线状态
func (s *UserService) GetUserOnlineStatus(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("用户ID不能为空")
	}

	query := map[string]string{
		"user_id": userID,
	}

	resp, err := s.client.Get(ctx, "user/online", query)
	if err != nil {
		return false, err
	}

	var result struct {
		Online bool `json:"online"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return false, fmt.Errorf("解析在线状态失败: %w", err)
	}

	return result.Online, nil
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, params UpdateUserParams) (*User, error) {
	// 构建请求参数
	requestParams := make(map[string]interface{})

	if params.Username != "" {
		requestParams["username"] = params.Username
	}
	if params.Avatar != "" {
		requestParams["avatar"] = params.Avatar
	}
	if params.Banner != "" {
		requestParams["banner"] = params.Banner
	}

	resp, err := s.client.Post(ctx, "user/update", requestParams)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return &user, nil
}

// UpdateUserParams 更新用户信息参数
type UpdateUserParams struct {
	Username string `json:"username,omitempty"` // 用户名
	Avatar   string `json:"avatar,omitempty"`   // 头像（base64或URL）
	Banner   string `json:"banner,omitempty"`   // 横幅图片URL
}

// BlockUser 屏蔽用户
func (s *UserService) BlockUser(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"user_id": userID,
	}

	_, err := s.client.Post(ctx, "user/block", params)
	return err
}

// UnblockUser 取消屏蔽用户
func (s *UserService) UnblockUser(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"user_id": userID,
	}

	_, err := s.client.Post(ctx, "user/unblock", params)
	return err
}

// GetBlockedUsers 获取被屏蔽的用户列表
func (s *UserService) GetBlockedUsers(ctx context.Context) ([]User, error) {
	resp, err := s.client.Get(ctx, "user/blocked", nil)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		return nil, fmt.Errorf("解析屏蔽用户列表失败: %w", err)
	}

	return users, nil
}

// SetOnline 上线机器人（仅限Webhook使用）
func (s *UserService) SetOnline(ctx context.Context) error {
	_, err := s.client.Post(ctx, "user/online", nil)
	return err
}

// SetOffline 下线机器人（仅限Webhook使用）
func (s *UserService) SetOffline(ctx context.Context) error {
	_, err := s.client.Post(ctx, "user/offline", nil)
	return err
}

// GetOnlineStatus 获取机器人在线状态
func (s *UserService) GetOnlineStatus(ctx context.Context) (*OnlineStatus, error) {
	resp, err := s.client.Get(ctx, "user/get-online-status", nil)
	if err != nil {
		return nil, err
	}

	var status OnlineStatus
	if err := json.Unmarshal(resp.Data, &status); err != nil {
		return nil, fmt.Errorf("解析在线状态失败: %w", err)
	}

	return &status, nil
}

// OnlineStatus 在线状态信息
type OnlineStatus struct {
	Online   bool     `json:"online"`     // 是否在线
	OnlineOS []string `json:"online_os"`  // 在线的平台列表
}
 