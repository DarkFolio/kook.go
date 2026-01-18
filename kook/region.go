package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// RegionService 服务器区域相关API服务
type RegionService struct {
	client *Client
}

// GetRegionList 获取可用区域列表
func (s *RegionService) GetRegionList(ctx context.Context) ([]Region, error) {
	resp, err := s.client.Get(ctx, "guild/regions", nil)
	if err != nil {
		return nil, err
	}

	var regions []Region
	if err := json.Unmarshal(resp.Data, &regions); err != nil {
		return nil, fmt.Errorf("解析区域列表失败: %w", err)
	}

	return regions, nil
}

// 数据结构定义

// Region 服务器区域信息
type Region struct {
	ID       string `json:"id"`       // 区域ID
	Name     string `json:"name"`     // 区域名称
	Crowding int    `json:"crowding"` // 拥挤程度（百分比）
} 