package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// CouponService 优惠券相关API服务
type CouponService struct {
	client *Client
}

// ExchangeCoupon 兑换优惠券
func (s *CouponService) ExchangeCoupon(ctx context.Context, code string) (*CouponExchangeResult, error) {
	if code == "" {
		return nil, fmt.Errorf("优惠券代码不能为空")
	}

	params := map[string]interface{}{
		"code": code,
	}

	resp, err := s.client.Post(ctx, "coupon/exchange", params)
	if err != nil {
		return nil, err
	}

	var result CouponExchangeResult
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析兑换结果失败: %w", err)
	}

	return &result, nil
}

// GetCoupons 获取优惠券列表
func (s *CouponService) GetCoupons(ctx context.Context, page, pageSize int) (*CouponListResponse, error) {
	query := make(map[string]string)

	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}
	if pageSize > 0 {
		query["page_size"] = fmt.Sprintf("%d", pageSize)
	}

	resp, err := s.client.Get(ctx, "coupon/list", query)
	if err != nil {
		return nil, err
	}

	var result CouponListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析优惠券列表失败: %w", err)
	}

	return &result, nil
}

// UseCoupon 使用优惠券
func (s *CouponService) UseCoupon(ctx context.Context, couponID string, orderID string) error {
	if couponID == "" {
		return fmt.Errorf("优惠券ID不能为空")
	}
	if orderID == "" {
		return fmt.Errorf("订单ID不能为空")
	}

	params := map[string]interface{}{
		"coupon_id": couponID,
		"order_id":  orderID,
	}

	_, err := s.client.Post(ctx, "coupon/use", params)
	return err
}

// 数据结构定义

// Coupon 优惠券信息
type Coupon struct {
	ID          string `json:"id"`           // 优惠券ID
	Code        string `json:"code"`         // 优惠券代码
	Name        string `json:"name"`         // 名称
	Description string `json:"description"`  // 描述
	Type        int    `json:"type"`         // 类型：1折扣，2满减
	Value       int    `json:"value"`        // 值（分）
	MinAmount   int    `json:"min_amount"`   // 最小金额（分）
	ExpiredAt   int64  `json:"expired_at"`   // 过期时间
	Used        bool   `json:"used"`         // 是否已使用
	UsedAt      int64  `json:"used_at"`      // 使用时间
	CreatedAt   int64  `json:"created_at"`   // 创建时间
}

// CouponExchangeResult 优惠券兑换结果
type CouponExchangeResult struct {
	Success bool   `json:"success"`         // 是否成功
	Message string `json:"message"`         // 消息
	Coupon  Coupon `json:"coupon,omitempty"` // 优惠券信息
	Items   []Item `json:"items,omitempty"`  // 物品列表
}

// CouponListResponse 优惠券列表响应
type CouponListResponse struct {
	Items []Coupon       `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// 优惠券类型常量
const (
	CouponTypeDiscount = 1 // 折扣券
	CouponTypeReduction = 2 // 满减券
) 