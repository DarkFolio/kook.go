package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// FriendService 好友相关API服务
type FriendService struct {
	client *Client
}

// SendFriendRequest 发送好友请求
func (s *FriendService) SendFriendRequest(ctx context.Context, params SendFriendRequestParams) error {
	if params.UserCode == "" {
		return fmt.Errorf("用户识别码不能为空")
	}

	requestParams := map[string]interface{}{
		"user_code": params.UserCode,
		"from":      params.From,
	}

	if params.From == 2 && params.GuildID != "" {
		requestParams["guild_id"] = params.GuildID
	}

	_, err := s.client.Post(ctx, "friend/request", requestParams)
	return err
}

// GetFriendsList 获取好友列表
func (s *FriendService) GetFriendsList(ctx context.Context) (*FriendsListResponse, error) {
	resp, err := s.client.Get(ctx, "friends", nil)
	if err != nil {
		return nil, err
	}

	var result FriendsListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析好友列表失败: %w", err)
	}

	return &result, nil
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"user_id": userID,
	}

	_, err := s.client.Post(ctx, "friend/delete", params)
	return err
}

// HandleFriendRequest 处理好友请求
func (s *FriendService) HandleFriendRequest(ctx context.Context, requestID string, accept bool) error {
	if requestID == "" {
		return fmt.Errorf("请求ID不能为空")
	}

	params := map[string]interface{}{
		"id":     requestID,
		"accept": accept,
	}

	_, err := s.client.Post(ctx, "friend/handle-request", params)
	return err
}

// AcceptFriendRequest 接受好友请求
func (s *FriendService) AcceptFriendRequest(ctx context.Context, requestID string) error {
	return s.HandleFriendRequest(ctx, requestID, true)
}

// RejectFriendRequest 拒绝好友请求
func (s *FriendService) RejectFriendRequest(ctx context.Context, requestID string) error {
	return s.HandleFriendRequest(ctx, requestID, false)
}

// 数据结构定义

// SendFriendRequestParams 发送好友请求参数
type SendFriendRequestParams struct {
	UserCode string `json:"user_code"`         // 用户识别码，格式: username#identify_num
	From     int    `json:"from"`              // 请求来源：0直接添加，1普通添加，2从服务器添加
	GuildID  string `json:"guild_id,omitempty"` // 服务器ID（当from=2时必填）
}

// FriendRequest 好友请求信息
type FriendRequest struct {
	ID      string `json:"id"`       // 请求ID
	UserID  string `json:"user_id"`  // 用户ID
	User    User   `json:"user"`     // 用户信息
	Status  int    `json:"status"`   // 请求状态
	Time    int64  `json:"time"`     // 请求时间
	Message string `json:"message"`  // 请求消息
}

// FriendsListResponse 好友列表响应
type FriendsListResponse struct {
	Request []FriendRequest `json:"request"` // 好友请求列表
	Friend  []User          `json:"friend"`  // 好友列表
	Blocked []User          `json:"blocked"` // 被屏蔽用户列表
}

// 好友请求来源常量
const (
	FriendRequestFromDirect = 0 // 直接添加
	FriendRequestFromNormal = 1 // 普通添加
	FriendRequestFromGuild  = 2 // 从服务器添加
) 