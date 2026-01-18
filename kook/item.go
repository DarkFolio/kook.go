package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// ItemService 物品相关API服务
type ItemService struct {
	client *Client
}

// GetItemList 获取物品列表
func (s *ItemService) GetItemList(ctx context.Context, category string) (*ItemListResponse, error) {
	query := make(map[string]string)

	if category != "" {
		query["category"] = category
	}

	resp, err := s.client.Get(ctx, "item/list", query)
	if err != nil {
		return nil, err
	}

	var result ItemListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析物品列表失败: %w", err)
	}

	return &result, nil
}

// GetBag 获取背包
func (s *ItemService) GetBag(ctx context.Context) ([]BagItem, error) {
	resp, err := s.client.Get(ctx, "item/bag", nil)
	if err != nil {
		return nil, err
	}

	var items []BagItem
	if err := json.Unmarshal(resp.Data, &items); err != nil {
		return nil, fmt.Errorf("解析背包失败: %w", err)
	}

	return items, nil
}

// UseItem 使用物品
func (s *ItemService) UseItem(ctx context.Context, userItemID int) error {
	if userItemID <= 0 {
		return fmt.Errorf("物品ID不能为空")
	}

	params := map[string]interface{}{
		"user_item_id": userItemID,
	}

	_, err := s.client.Post(ctx, "item/using", params)
	return err
}

// CancelUseItem 取消使用物品
func (s *ItemService) CancelUseItem(ctx context.Context, userItemID int) error {
	if userItemID <= 0 {
		return fmt.Errorf("物品ID不能为空")
	}

	params := map[string]interface{}{
		"user_item_id": userItemID,
	}

	_, err := s.client.Post(ctx, "item/cancel-use", params)
	return err
}

// DeleteItems 删除物品
func (s *ItemService) DeleteItems(ctx context.Context, userItemIDs []int) error {
	if len(userItemIDs) == 0 {
		return fmt.Errorf("物品ID列表不能为空")
	}

	params := map[string]interface{}{
		"user_item_ids": userItemIDs,
	}

	_, err := s.client.Post(ctx, "item/delete", params)
	return err
}

// 数据结构定义

// Item 物品信息
type Item struct {
	ID           string `json:"id"`            // 物品ID
	Status       int    `json:"status"`        // 状态
	Type         int    `json:"type"`          // 类型
	Name         string `json:"name"`          // 名称
	Price        int    `json:"price"`         // 价格（分）
	OriginPrice  int    `json:"origin_price"`  // 原价（分）
	ServiceTime  int    `json:"service_time"`  // 服务时间
	DiscountLabel string `json:"discount_label"` // 折扣标签
	IAPCode      string `json:"iap_code"`      // IAP代码
}

// BagItem 背包物品
type BagItem struct {
	UserItemID int    `json:"user_item_id"` // 用户物品ID
	ItemID     string `json:"item_id"`      // 物品ID
	Item       Item   `json:"item"`         // 物品信息
	Count      int    `json:"count"`        // 数量
	ExpiredAt  int64  `json:"expired_at"`   // 过期时间
	Using      bool   `json:"using"`        // 是否正在使用
}

// ItemListResponse 物品列表响应
type ItemListResponse struct {
	Items []Item `json:"items"`
}

// 物品分类常量
const (
	ItemCategoryAll        = "all"         // 全部
	ItemCategoryTimeLimit  = "time_limit"  // 限时物品
	ItemCategoryDecoration = "decoration"  // 装饰物品
	ItemCategoryAction     = "action"      // 动作物品
) 