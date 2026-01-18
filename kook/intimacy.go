package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// IntimacyService 亲密度相关API服务
type IntimacyService struct {
	client *Client
}

// GetIntimacy 获取用户亲密度
func (s *IntimacyService) GetIntimacy(ctx context.Context, userID string) (*Intimacy, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	query := map[string]string{
		"user_id": userID,
	}

	resp, err := s.client.Get(ctx, "intimacy/index", query)
	if err != nil {
		return nil, err
	}

	var intimacy Intimacy
	if err := json.Unmarshal(resp.Data, &intimacy); err != nil {
		return nil, fmt.Errorf("解析亲密度信息失败: %w", err)
	}

	return &intimacy, nil
}

// UpdateIntimacy 更新用户亲密度
func (s *IntimacyService) UpdateIntimacy(ctx context.Context, userID string, score int, socialInfo string, imgID string) (*Intimacy, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"user_id": userID,
	}

	if score >= 0 {
		params["score"] = score
	}
	if socialInfo != "" {
		params["social_info"] = socialInfo
	}
	if imgID != "" {
		params["img_id"] = imgID
	}

	resp, err := s.client.Post(ctx, "intimacy/update", params)
	if err != nil {
		return nil, err
	}

	var intimacy Intimacy
	if err := json.Unmarshal(resp.Data, &intimacy); err != nil {
		return nil, fmt.Errorf("解析亲密度信息失败: %w", err)
	}

	return &intimacy, nil
}

// 数据结构定义

// Intimacy 亲密度信息
type Intimacy struct {
	UserID     string `json:"user_id"`     // 用户ID
	Score      int    `json:"score"`       // 亲密度分数
	SocialInfo string `json:"social_info"` // 社交信息
	LastRead   int64  `json:"last_read"`   // 最后阅读时间
	LastModify int64  `json:"last_modify"` // 最后修改时间
	ImgID      string `json:"img_id"`      // 图片ID
	ImgURL     string `json:"img_url"`     // 图片URL
} 