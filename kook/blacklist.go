package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// BlacklistService 屏蔽/黑名单相关API服务
type BlacklistService struct {
	client *Client
}

// GetBlacklistUsers 获取屏蔽用户列表
func (s *BlacklistService) GetBlacklistUsers(ctx context.Context, guildID string, page, pageSize int) (*BlacklistResponse, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	if page > 0 {
		query["page"] = strconv.Itoa(page)
	}
	if pageSize > 0 && pageSize <= 50 {
		query["page_size"] = strconv.Itoa(pageSize)
	}

	resp, err := s.client.Get(ctx, "blacklist/list", query)
	if err != nil {
		return nil, err
	}

	var result BlacklistResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析屏蔽用户列表失败: %w", err)
	}

	return &result, nil
}

// CreateBlacklistUser 屏蔽用户
func (s *BlacklistService) CreateBlacklistUser(ctx context.Context, guildID, userID string, remark string, delMsgDays int) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"user_id":  userID,
	}

	if remark != "" {
		params["remark"] = remark
	}
	if delMsgDays > 0 {
		params["del_msg_days"] = delMsgDays
	}

	_, err := s.client.Post(ctx, "blacklist/create", params)
	return err
}

// DeleteBlacklistUser 取消屏蔽用户
func (s *BlacklistService) DeleteBlacklistUser(ctx context.Context, guildID, userID string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"user_id":  userID,
	}

	_, err := s.client.Post(ctx, "blacklist/delete", params)
	return err
}

// 数据结构定义

// BlacklistUser 屏蔽用户信息
type BlacklistUser struct {
	User       User   `json:"user"`        // 用户信息
	Remark     string `json:"remark"`      // 屏蔽备注
	UserID     string `json:"user_id"`     // 用户ID
	CreatedAt  int64  `json:"created_at"`  // 屏蔽时间
	UpdatedAt  int64  `json:"updated_at"`  // 更新时间
}

// BlacklistResponse 屏蔽用户列表响应
type BlacklistResponse struct {
	Items []BlacklistUser `json:"items"`
	Meta  PaginationMeta  `json:"meta"`
	Sort  map[string]int  `json:"sort"`
} 