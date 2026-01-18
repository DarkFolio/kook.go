package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// EmojiService 表情包相关API服务
type EmojiService struct {
	client *Client
}

// GetEmojiList 获取服务器表情列表
func (s *EmojiService) GetEmojiList(ctx context.Context, guildID string, page, pageSize int) (*EmojiListResponse, error) {
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

	resp, err := s.client.Get(ctx, "emoji/list", query)
	if err != nil {
		return nil, err
	}

	var result EmojiListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析表情列表失败: %w", err)
	}

	return &result, nil
}

// CreateEmoji 创建表情
func (s *EmojiService) CreateEmoji(ctx context.Context, name, guildID string, emoji interface{}) (*Emoji, error) {
	if name == "" {
		return nil, fmt.Errorf("表情名称不能为空")
	}
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	params := map[string]interface{}{
		"name":     name,
		"guild_id": guildID,
		"emoji":    emoji, // 可以是文件或URL
	}

	resp, err := s.client.Post(ctx, "emoji/create", params)
	if err != nil {
		return nil, err
	}

	var result Emoji
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析表情信息失败: %w", err)
	}

	return &result, nil
}

// UpdateEmoji 更新表情
func (s *EmojiService) UpdateEmoji(ctx context.Context, id, name string) (*Emoji, error) {
	if id == "" {
		return nil, fmt.Errorf("表情ID不能为空")
	}

	params := map[string]interface{}{
		"id": id,
	}

	if name != "" {
		params["name"] = name
	}

	resp, err := s.client.Post(ctx, "emoji/update", params)
	if err != nil {
		return nil, err
	}

	var result Emoji
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析表情信息失败: %w", err)
	}

	return &result, nil
}

// DeleteEmoji 删除表情
func (s *EmojiService) DeleteEmoji(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("表情ID不能为空")
	}

	params := map[string]interface{}{
		"id": id,
	}

	_, err := s.client.Post(ctx, "emoji/delete", params)
	return err
}

// 数据结构定义

// Emoji 表情信息
type Emoji struct {
	ID     string `json:"id"`     // 表情ID
	Name   string `json:"name"`   // 表情名称
	URL    string `json:"url"`    // 表情URL
	UserID string `json:"user_id"` // 创建者ID
}

// EmojiListResponse 表情列表响应
type EmojiListResponse struct {
	Items []Emoji        `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
} 