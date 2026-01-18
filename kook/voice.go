package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// VoiceService 语音相关API服务
type VoiceService struct {
	client *Client
}

// JoinVoiceChannel 加入语音频道
func (s *VoiceService) JoinVoiceChannel(ctx context.Context, channelID string) (*VoiceConnectionInfo, error) {
	if channelID == "" {
		return nil, fmt.Errorf("频道ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
	}

	resp, err := s.client.Post(ctx, "voice/join", params)
	if err != nil {
		return nil, err
	}

	var connInfo VoiceConnectionInfo
	if err := json.Unmarshal(resp.Data, &connInfo); err != nil {
		return nil, fmt.Errorf("解析语音连接信息失败: %w", err)
	}

	return &connInfo, nil
}

// LeaveVoiceChannel 离开语音频道
func (s *VoiceService) LeaveVoiceChannel(ctx context.Context, channelID string) error {
	if channelID == "" {
		return fmt.Errorf("频道ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
	}

	_, err := s.client.Post(ctx, "voice/leave", params)
	return err
}

// GetVoiceChannelUsers 获取语音频道用户列表
func (s *VoiceService) GetVoiceChannelUsers(ctx context.Context, channelID string) ([]VoiceUser, error) {
	if channelID == "" {
		return nil, fmt.Errorf("频道ID不能为空")
	}

	query := map[string]string{
		"channel_id": channelID,
	}

	resp, err := s.client.Get(ctx, "voice/users", query)
	if err != nil {
		return nil, err
	}

	var users []VoiceUser
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		return nil, fmt.Errorf("解析语音频道用户列表失败: %w", err)
	}

	return users, nil
}

// MuteUser 静音用户
func (s *VoiceService) MuteUser(ctx context.Context, channelID, userID string) error {
	if channelID == "" {
		return fmt.Errorf("频道ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
		"user_id":    userID,
	}

	_, err := s.client.Post(ctx, "voice/mute", params)
	return err
}

// UnmuteUser 取消静音用户
func (s *VoiceService) UnmuteUser(ctx context.Context, channelID, userID string) error {
	if channelID == "" {
		return fmt.Errorf("频道ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
		"user_id":    userID,
	}

	_, err := s.client.Post(ctx, "voice/unmute", params)
	return err
}

// DeafenUser 闭麦用户
func (s *VoiceService) DeafenUser(ctx context.Context, channelID, userID string) error {
	if channelID == "" {
		return fmt.Errorf("频道ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
		"user_id":    userID,
	}

	_, err := s.client.Post(ctx, "voice/deafen", params)
	return err
}

// UndeafenUser 取消闭麦用户
func (s *VoiceService) UndeafenUser(ctx context.Context, channelID, userID string) error {
	if channelID == "" {
		return fmt.Errorf("频道ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
		"user_id":    userID,
	}

	_, err := s.client.Post(ctx, "voice/undeafen", params)
	return err
}

// 数据结构定义

// VoiceConnectionInfo 语音连接信息
type VoiceConnectionInfo struct {
	GatewayURL string `json:"gateway_url"` // 语音网关URL
	Token      string `json:"token"`       // 语音令牌
	Endpoint   string `json:"endpoint"`    // 连接端点
	SessionID  string `json:"session_id"`  // 会话ID
}

// VoiceUser 语音频道用户
type VoiceUser struct {
	User        User `json:"user"`         // 用户信息
	Muted       bool `json:"muted"`        // 是否被静音
	Deafened    bool `json:"deafened"`     // 是否被闭麦
	SelfMuted   bool `json:"self_muted"`   // 是否自我静音
	SelfDeafened bool `json:"self_deafened"` // 是否自我闭麦
	Speaking    bool `json:"speaking"`     // 是否正在说话
} 