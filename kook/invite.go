package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// InviteService 邀请相关API服务
type InviteService struct {
	client *Client
}

// GetInviteList 获取邀请列表
func (s *InviteService) GetInviteList(ctx context.Context, guildID string, page, pageSize int) (*ListInvitesResponse, error) {
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

	resp, err := s.client.Get(ctx, "invite/list", query)
	if err != nil {
		return nil, err
	}

	var result ListInvitesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析邀请列表失败: %w", err)
	}

	return &result, nil
}

// CreateInvite 创建邀请
func (s *InviteService) CreateInvite(ctx context.Context, params CreateInviteParams) (*Invite, error) {
	requestParams := make(map[string]interface{})

	if params.GuildID != "" {
		requestParams["guild_id"] = params.GuildID
	}
	if params.ChannelID != "" {
		requestParams["channel_id"] = params.ChannelID
	}
	if params.Duration > 0 {
		requestParams["duration"] = params.Duration
	}
	if params.Setting > 0 {
		requestParams["setting"] = params.Setting
	}

	resp, err := s.client.Post(ctx, "invite/create", requestParams)
	if err != nil {
		return nil, err
	}

	var invite Invite
	if err := json.Unmarshal(resp.Data, &invite); err != nil {
		return nil, fmt.Errorf("解析邀请信息失败: %w", err)
	}

	return &invite, nil
}

// DeleteInvite 删除邀请
func (s *InviteService) DeleteInvite(ctx context.Context, urlCode string) error {
	if urlCode == "" {
		return fmt.Errorf("邀请码不能为空")
	}

	params := map[string]interface{}{
		"url_code": urlCode,
	}

	_, err := s.client.Post(ctx, "invite/delete", params)
	return err
}

// 数据结构定义

// Invite 邀请信息
type Invite struct {
	GuildID     string `json:"guild_id"`     // 服务器ID
	ChannelID   string `json:"channel_id"`   // 频道ID
	URLCode     string `json:"url_code"`     // 邀请码
	URL         string `json:"url"`          // 邀请链接
	User        User   `json:"user"`         // 创建者信息
	CreatedAt   int64  `json:"created_at"`   // 创建时间
	UpdatedAt   int64  `json:"updated_at"`   // 更新时间
	ExpiredAt   int64  `json:"expired_at"`   // 过期时间
	Duration    int    `json:"duration"`     // 有效期（秒）
	Setting     int    `json:"setting"`      // 设置
	RemainTimes int    `json:"remain_times"` // 剩余使用次数
}

// CreateInviteParams 创建邀请参数
type CreateInviteParams struct {
	GuildID   string `json:"guild_id,omitempty"`   // 服务器ID
	ChannelID string `json:"channel_id,omitempty"` // 频道ID
	Duration  int    `json:"duration,omitempty"`   // 有效期（秒）：0永久，1800半小时，3600一小时，21600六小时，43200十二小时，86400一天，604800七天
	Setting   int    `json:"setting,omitempty"`    // 设置：次数限制，默认-1无限制
}

// ListInvitesResponse 邀请列表响应
type ListInvitesResponse struct {
	Items []Invite       `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// 邀请有效期常量
const (
	InviteDurationForever   = 0      // 永久
	InviteDurationHalfHour  = 1800   // 半小时
	InviteDurationOneHour   = 3600   // 一小时
	InviteDurationSixHours  = 21600  // 六小时
	InviteDurationTwelveHours = 43200 // 十二小时
	InviteDurationOneDay    = 86400  // 一天
	InviteDurationOneWeek   = 604800 // 七天
) 