package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// BoostService 助力相关API服务
type BoostService struct {
	client *Client
}

// GetUnusedBoostNum 获取未使用的助力数量
func (s *BoostService) GetUnusedBoostNum(ctx context.Context) (*UnusedBoostInfo, error) {
	resp, err := s.client.Get(ctx, "guild-boost/get-unused-boost-num", nil)
	if err != nil {
		return nil, err
	}

	var info UnusedBoostInfo
	if err := json.Unmarshal(resp.Data, &info); err != nil {
		return nil, fmt.Errorf("解析未使用助力数量失败: %w", err)
	}

	return &info, nil
}

// UseBoost 使用助力
func (s *BoostService) UseBoost(ctx context.Context, guildID string, count int) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if count <= 0 {
		return fmt.Errorf("助力数量必须大于0")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"count":    count,
	}

	_, err := s.client.Post(ctx, "boost/use", params)
	return err
}

// GetGuildBoosts 获取服务器助力列表
func (s *BoostService) GetGuildBoosts(ctx context.Context, guildID string, page, pageSize int) (*GuildBoostListResponse, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}
	if pageSize > 0 {
		query["page_size"] = fmt.Sprintf("%d", pageSize)
	}

	resp, err := s.client.Get(ctx, "guild-boost/list", query)
	if err != nil {
		return nil, err
	}

	var result GuildBoostListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析服务器助力列表失败: %w", err)
	}

	return &result, nil
}

// CancelBoost 取消助力
func (s *BoostService) CancelBoost(ctx context.Context, guildID string, boostID string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if boostID == "" {
		return fmt.Errorf("助力ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"boost_id": boostID,
	}

	_, err := s.client.Post(ctx, "boost/cancel", params)
	return err
}

// 数据结构定义

// UnusedBoostInfo 未使用助力信息
type UnusedBoostInfo struct {
	UnusedBoostNum int `json:"unused_boost_num"` // 未使用助力数量
}

// GuildBoost 服务器助力信息
type GuildBoost struct {
	ID        string `json:"id"`         // 助力ID
	GuildID   string `json:"guild_id"`   // 服务器ID
	UserID    string `json:"user_id"`    // 用户ID
	User      User   `json:"user"`       // 用户信息
	StartTime int64  `json:"start_time"` // 开始时间
	EndTime   int64  `json:"end_time"`   // 结束时间
	Level     int    `json:"level"`      // 助力等级
	Status    int    `json:"status"`     // 状态：1活跃，0已结束
}

// GuildBoostListResponse 服务器助力列表响应
type GuildBoostListResponse struct {
	Items []GuildBoost   `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
} 