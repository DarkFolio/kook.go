package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// OrderService 订单相关API服务
type OrderService struct {
	client *Client
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, params CreateOrderParams) (*Order, error) {
	if len(params.Products) == 0 {
		return nil, fmt.Errorf("商品列表不能为空")
	}

	requestParams := map[string]interface{}{
		"products": params.Products,
	}

	if params.Platform > 0 {
		requestParams["platform"] = params.Platform
	} else {
		requestParams["platform"] = 1 // 默认平台
	}

	if params.RequestPay {
		requestParams["request_pay"] = params.RequestPay
	}

	resp, err := s.client.Post(ctx, "order/create", requestParams)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp.Data, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %w", err)
	}

	return &order, nil
}

// GetOrderStatus 获取订单状态
func (s *OrderService) GetOrderStatus(ctx context.Context, orderID string) (*Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("订单ID不能为空")
	}

	query := map[string]string{
		"order_id": orderID,
	}

	resp, err := s.client.Get(ctx, "order/status", query)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp.Data, &order); err != nil {
		return nil, fmt.Errorf("解析订单信息失败: %w", err)
	}

	return &order, nil
}

// GetOrders 获取订单列表
func (s *OrderService) GetOrders(ctx context.Context, page, pageSize int) (*OrderListResponse, error) {
	query := make(map[string]string)

	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}
	if pageSize > 0 {
		query["page_size"] = fmt.Sprintf("%d", pageSize)
	}

	resp, err := s.client.Get(ctx, "order/list", query)
	if err != nil {
		return nil, err
	}

	var result OrderListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析订单列表失败: %w", err)
	}

	return &result, nil
}

// 数据结构定义

// CreateOrderParams 创建订单参数
type CreateOrderParams struct {
	Products   []OrderProduct `json:"products"`    // 商品列表
	Platform   int            `json:"platform"`    // 平台：1默认
	RequestPay bool           `json:"request_pay"` // 是否请求支付
}

// OrderProduct 订单商品
type OrderProduct struct {
	ID    int `json:"id"`    // 商品ID
	Count int `json:"count"` // 数量
}

// Order 订单信息
type Order struct {
	ID                string      `json:"id"`                  // 订单ID
	Status            int         `json:"status"`              // 订单状态
	UserID            string      `json:"user_id"`             // 用户ID
	TotalFee          int         `json:"total_fee"`           // 总费用（分）
	PayFee            int         `json:"pay_fee"`             // 支付费用（分）
	Paid              bool        `json:"paid"`                // 是否已支付
	PayTime           int64       `json:"pay_time"`            // 支付时间
	CreateTime        int64       `json:"create_time"`         // 创建时间
	Products          []Product   `json:"products"`            // 商品列表
	UsageInfo         string      `json:"usage_info"`          // 使用信息
	ItemEntitiesDesc  string      `json:"item_entities_desc"`  // 物品实体描述
	PayData           *PayData    `json:"paydata,omitempty"`   // 支付数据
}

// Product 商品信息
type Product struct {
	ID         int            `json:"id"`          // 商品ID
	ItemID     int            `json:"item_id"`     // 物品ID
	Item       ProductItem    `json:"item"`        // 物品信息
	Total      int            `json:"total"`       // 数量
	ExpireTime int64          `json:"expire_time"` // 过期时间
}

// ProductItem 商品物品信息
type ProductItem struct {
	ID              int                    `json:"id"`               // 物品ID
	Name            string                 `json:"name"`             // 名称
	Desc            string                 `json:"desc"`             // 描述
	CD              int                    `json:"cd"`               // 冷却时间
	Categories      []string               `json:"categories"`       // 分类
	Label           int                    `json:"label"`            // 标签
	LabelName       string                 `json:"label_name"`       // 标签名称
	Quality         int                    `json:"quality"`          // 品质
	Icon            string                 `json:"icon"`             // 图标
	IconThumb       string                 `json:"icon_thumb"`       // 图标缩略图
	IconExpired     string                 `json:"icon_expired"`     // 过期图标
	QualityResource QualityResource        `json:"quality_resource"` // 品质资源
	Resources       ProductItemResources   `json:"resources"`        // 资源
	Position        string                 `json:"position"`         // 位置
}

// QualityResource 品质资源
type QualityResource struct {
	Color string `json:"color"` // 颜色
	Small string `json:"small"` // 小图
	Big   string `json:"big"`   // 大图
}

// ProductItemResources 商品物品资源
type ProductItemResources struct {
	GIF            string `json:"gif"`             // GIF图片
	Height         int    `json:"height"`          // 高度
	PAG            string `json:"pag"`             // PAG文件
	Percent        int    `json:"percent"`         // 百分比
	Preview        string `json:"preview"`         // 预览图
	PreviewExpired string `json:"preview_expired"` // 预览过期图
	Time           int    `json:"time"`            // 时间
	Type           string `json:"type"`            // 类型
	WEBP           string `json:"webp"`            // WEBP图片
	Width          int    `json:"width"`           // 宽度
}

// PayData 支付数据
type PayData struct {
	ID          string `json:"id"`           // 支付ID
	PayFee      string `json:"pay_fee"`      // 支付费用
	QRCode      string `json:"qr_code"`      // 二维码
	QRCodeURL   string `json:"qr_code_url"`  // 二维码URL
	ExpiredTime int64  `json:"expired_time"` // 过期时间
	MobilePay   string `json:"mobile_pay"`   // 移动支付
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Items []Order        `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
} 