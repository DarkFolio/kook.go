package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// LiveService 直播相关API服务
type LiveService struct {
	client *Client
}

// StartLive 开始直播
func (s *LiveService) StartLive(ctx context.Context, channelID, title string) (*LiveInfo, error) {
	if channelID == "" {
		return nil, fmt.Errorf("频道ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
	}

	if title != "" {
		params["title"] = title
	}

	resp, err := s.client.Post(ctx, "live/start", params)
	if err != nil {
		return nil, err
	}

	var liveInfo LiveInfo
	if err := json.Unmarshal(resp.Data, &liveInfo); err != nil {
		return nil, fmt.Errorf("解析直播信息失败: %w", err)
	}

	return &liveInfo, nil
}

// StopLive 停止直播
func (s *LiveService) StopLive(ctx context.Context, channelID string) error {
	if channelID == "" {
		return fmt.Errorf("频道ID不能为空")
	}

	params := map[string]interface{}{
		"channel_id": channelID,
	}

	_, err := s.client.Post(ctx, "live/stop", params)
	return err
}

// GetLiveInfo 获取直播信息
func (s *LiveService) GetLiveInfo(ctx context.Context, channelID string) (*LiveInfo, error) {
	if channelID == "" {
		return nil, fmt.Errorf("频道ID不能为空")
	}

	query := map[string]string{
		"channel_id": channelID,
	}

	resp, err := s.client.Get(ctx, "live/info", query)
	if err != nil {
		return nil, err
	}

	var liveInfo LiveInfo
	if err := json.Unmarshal(resp.Data, &liveInfo); err != nil {
		return nil, fmt.Errorf("解析直播信息失败: %w", err)
	}

	return &liveInfo, nil
}

// 数据结构定义

// LiveInfo 直播信息
type LiveInfo struct {
	ChannelID   string `json:"channel_id"`   // 频道ID
	Title       string `json:"title"`        // 直播标题
	Status      int    `json:"status"`       // 直播状态：0未开始，1直播中，2已结束
	ViewerCount int    `json:"viewer_count"` // 观看人数
	StartTime   int64  `json:"start_time"`   // 开始时间
	EndTime     int64  `json:"end_time"`     // 结束时间
	StreamURL   string `json:"stream_url"`   // 推流地址
	PlayURL     string `json:"play_url"`     // 播放地址
} 