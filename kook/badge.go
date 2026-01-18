package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// BadgeService 徽章相关API服务
type BadgeService struct {
	client *Client
}

// GetGuildBadges 获取服务器徽章列表
func (s *BadgeService) GetGuildBadges(ctx context.Context, guildID string) ([]Badge, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	resp, err := s.client.Get(ctx, "badge/guild", query)
	if err != nil {
		return nil, err
	}

	var badges []Badge
	if err := json.Unmarshal(resp.Data, &badges); err != nil {
		return nil, fmt.Errorf("解析徽章列表失败: %w", err)
	}

	return badges, nil
}

// 数据结构定义

// Badge 徽章信息
type Badge struct {
	ID          string `json:"id"`          // 徽章ID
	Name        string `json:"name"`        // 徽章名称
	Description string `json:"description"` // 徽章描述
	Icon        string `json:"icon"`        // 徽章图标URL
	Type        int    `json:"type"`        // 徽章类型
	Level       int    `json:"level"`       // 徽章等级
	Unlocked    bool   `json:"unlocked"`    // 是否已解锁
} 