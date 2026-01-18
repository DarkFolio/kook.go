package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// GuildService 服务器相关API服务
type GuildService struct {
	client *Client
}

// GetGuildList 获取当前用户的服务器列表
func (s *GuildService) GetGuildList(ctx context.Context, page, pageSize int, sort string) (*ListGuildsResponse, error) {
	query := make(map[string]string)

	if page > 0 {
		query["page"] = strconv.Itoa(page)
	}
	if pageSize > 0 && pageSize <= 50 {
		query["page_size"] = strconv.Itoa(pageSize)
	}
	if sort != "" {
		query["sort"] = sort
	}

	resp, err := s.client.Get(ctx, "guild/list", query)
	if err != nil {
		return nil, err
	}

	var result ListGuildsResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析服务器列表失败: %w", err)
	}

	return &result, nil
}

// GetGuildInfo 获取服务器信息
func (s *GuildService) GetGuildInfo(ctx context.Context, guildID string) (*Guild, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	resp, err := s.client.Get(ctx, "guild/view", query)
	if err != nil {
		return nil, err
	}

	var guild Guild
	if err := json.Unmarshal(resp.Data, &guild); err != nil {
		return nil, fmt.Errorf("解析服务器信息失败: %w", err)
	}

	return &guild, nil
}

// CreateGuild 创建服务器
func (s *GuildService) CreateGuild(ctx context.Context, params CreateGuildParams) (*Guild, error) {
	if params.Name == "" {
		return nil, fmt.Errorf("服务器名称不能为空")
	}

	requestParams := map[string]interface{}{
		"name": params.Name,
	}

	if params.Icon != "" {
		requestParams["icon"] = params.Icon
	}
	if params.Region != "" {
		requestParams["region"] = params.Region
	} else {
		requestParams["region"] = "beijing" // 默认区域
	}
	if params.TemplateID > 0 {
		requestParams["template_id"] = params.TemplateID
	}

	resp, err := s.client.Post(ctx, "guild/create", requestParams)
	if err != nil {
		return nil, err
	}

	var guild Guild
	if err := json.Unmarshal(resp.Data, &guild); err != nil {
		return nil, fmt.Errorf("解析服务器信息失败: %w", err)
	}

	return &guild, nil
}

// UpdateGuild 更新服务器信息
func (s *GuildService) UpdateGuild(ctx context.Context, guildID string, params UpdateGuildParams) (*Guild, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	requestParams := map[string]interface{}{
		"guild_id": guildID,
	}

	if params.Name != "" {
		requestParams["name"] = params.Name
	}
	if params.Region != "" {
		requestParams["region"] = params.Region
	}
	if params.DefaultChannelID != "" {
		requestParams["default_channel_id"] = params.DefaultChannelID
	}
	if params.WelcomeChannelID != "" {
		requestParams["welcome_channel_id"] = params.WelcomeChannelID
	}
	if params.NotifyType >= 0 {
		requestParams["notify_type"] = params.NotifyType
	}
	if params.EnableOpen != nil {
		if *params.EnableOpen {
			requestParams["enable_open"] = 1
		} else {
			requestParams["enable_open"] = 0
		}
	}

	resp, err := s.client.Post(ctx, "guild/update", requestParams)
	if err != nil {
		return nil, err
	}

	var guild Guild
	if err := json.Unmarshal(resp.Data, &guild); err != nil {
		return nil, fmt.Errorf("解析服务器信息失败: %w", err)
	}

	return &guild, nil
}

// DeleteGuild 删除服务器
func (s *GuildService) DeleteGuild(ctx context.Context, guildID string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
	}

	_, err := s.client.Post(ctx, "guild/delete", params)
	return err
}

// LeaveGuild 离开服务器
func (s *GuildService) LeaveGuild(ctx context.Context, guildID string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
	}

	_, err := s.client.Post(ctx, "guild/leave", params)
	return err
}

// JoinGuild 加入服务器
func (s *GuildService) JoinGuild(ctx context.Context, params JoinGuildParams) (*JoinGuildResponse, error) {
	if params.Code == "" && params.ID == "" {
		return nil, fmt.Errorf("邀请码或服务器ID不能都为空")
	}

	query := make(map[string]string)
	if params.Code != "" {
		query["code"] = params.Code
	}
	if params.ID != "" {
		query["id"] = params.ID
	}

	resp, err := s.client.Get(ctx, "guild/join", query)
	if err != nil {
		return nil, err
	}

	var result JoinGuildResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析加入服务器响应失败: %w", err)
	}

	return &result, nil
}

// GetGuildMembers 获取服务器成员列表
func (s *GuildService) GetGuildMembers(ctx context.Context, guildID string, page, pageSize int, sort string) (*ListGuildMembersResponse, error) {
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
	if sort != "" {
		query["sort"] = sort
	}

	resp, err := s.client.Get(ctx, "guild/user-list", query)
	if err != nil {
		return nil, err
	}

	var result ListGuildMembersResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析服务器成员列表失败: %w", err)
	}

	return &result, nil
}

// GetGuildMember 获取服务器成员信息
func (s *GuildService) GetGuildMember(ctx context.Context, guildID, userID string) (*GuildMember, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
		"user_id":  userID,
	}

	resp, err := s.client.Get(ctx, "guild/user", query)
	if err != nil {
		return nil, err
	}

	var member GuildMember
	if err := json.Unmarshal(resp.Data, &member); err != nil {
		return nil, fmt.Errorf("解析服务器成员信息失败: %w", err)
	}

	return &member, nil
}

// KickGuildMember 踢出服务器成员
func (s *GuildService) KickGuildMember(ctx context.Context, guildID, userID string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"user_id":  userID,
	}

	_, err := s.client.Post(ctx, "guild/kickout", params)
	return err
}

// UpdateGuildMemberNickname 修改服务器成员昵称
func (s *GuildService) UpdateGuildMemberNickname(ctx context.Context, guildID, userID, nickname string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"nickname": nickname,
	}

	if userID != "" {
		params["user_id"] = userID
	}

	_, err := s.client.Post(ctx, "guild/nickname", params)
	return err
}

// GetRegions 获取可用的服务器区域列表
func (s *GuildService) GetRegions(ctx context.Context) (*ListRegionsResponse, error) {
	resp, err := s.client.Get(ctx, "guild/regions", nil)
	if err != nil {
		return nil, err
	}

	var result ListRegionsResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析区域列表失败: %w", err)
	}

	return &result, nil
}

// UpdateNickname 修改用户昵称
func (s *GuildService) UpdateNickname(ctx context.Context, guildID, userID, nickname string) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"user_id":  userID,
	}

	if nickname != "" {
		params["nickname"] = nickname
	}

	_, err := s.client.Post(ctx, "guild/nickname", params)
	return err
}

// UpdateGuildSettings 更新服务器设置
func (s *GuildService) UpdateGuildSettings(ctx context.Context, params UpdateGuildParams) (*Guild, error) {
	if params.GuildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	requestParams := map[string]interface{}{
		"guild_id": params.GuildID,
	}

	if params.Name != "" {
		requestParams["name"] = params.Name
	}
	if params.Region != "" {
		requestParams["region"] = params.Region
	}
	if params.DefaultChannelID != "" {
		requestParams["default_channel_id"] = params.DefaultChannelID
	}
	if params.WelcomeChannelID != "" {
		requestParams["welcome_channel_id"] = params.WelcomeChannelID
	}
	if params.EnableOpen != nil {
		if *params.EnableOpen {
			requestParams["enable_open"] = 1
		} else {
			requestParams["enable_open"] = 0
		}
	}
	if params.Icon != "" {
		requestParams["icon"] = params.Icon
	}
	if params.Banner != "" {
		requestParams["banner"] = params.Banner
	}

	resp, err := s.client.Post(ctx, "guild/update", requestParams)
	if err != nil {
		return nil, err
	}

	var guild Guild
	if err := json.Unmarshal(resp.Data, &guild); err != nil {
		return nil, fmt.Errorf("解析服务器信息失败: %w", err)
	}

	return &guild, nil
}

// GetGuildBoostInfo 获取服务器助力信息
func (s *GuildService) GetGuildBoostInfo(ctx context.Context, guildID string) (*GuildBoostInfo, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	resp, err := s.client.Get(ctx, "guild-boost/info", query)
	if err != nil {
		return nil, err
	}

	var boostInfo GuildBoostInfo
	if err := json.Unmarshal(resp.Data, &boostInfo); err != nil {
		return nil, fmt.Errorf("解析助力信息失败: %w", err)
	}

	return &boostInfo, nil
}

// 数据结构定义

// ListGuildsResponse 服务器列表响应
type ListGuildsResponse struct {
	Items []Guild        `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// CreateGuildParams 创建服务器参数
type CreateGuildParams struct {
	Name       string `json:"name"`        // 服务器名称
	Icon       string `json:"icon"`        // 服务器图标
	Region     string `json:"region"`      // 服务器区域
	TemplateID int    `json:"template_id"` // 模板ID
}

// UpdateGuildParams 更新服务器参数
type UpdateGuildParams struct {
	GuildID          string `json:"guild_id,omitempty"`               // 服务器ID
	Name             string `json:"name,omitempty"`               // 服务器名称
	Region           string `json:"region,omitempty"`             // 服务器区域
	DefaultChannelID string `json:"default_channel_id,omitempty"` // 默认频道ID
	WelcomeChannelID string `json:"welcome_channel_id,omitempty"` // 欢迎频道ID
	NotifyType       int    `json:"notify_type,omitempty"`        // 通知类型
	EnableOpen       *bool  `json:"enable_open,omitempty"`        // 是否开启公开
	Icon             string `json:"icon,omitempty"`             // 服务器图标
	Banner           string `json:"banner,omitempty"`             // 服务器横幅
}

// JoinGuildParams 加入服务器参数
type JoinGuildParams struct {
	Code string `json:"code,omitempty"` // 邀请码
	ID   string `json:"id,omitempty"`   // 服务器ID
}

// JoinGuildResponse 加入服务器响应
type JoinGuildResponse struct {
	Joined bool   `json:"joined"` // 是否已加入
	Guild  *Guild `json:"guild"`  // 服务器信息
}

// ListGuildMembersResponse 服务器成员列表响应
type ListGuildMembersResponse struct {
	Items []GuildMember  `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// ListRegionsResponse 区域列表响应
type ListRegionsResponse struct {
	Items []Region       `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// GuildBoostInfo 服务器助力信息
type GuildBoostInfo struct {
	BoostNum       int `json:"boost_num"`       // 助力数量
	BufferBoostNum int `json:"buffer_boost_num"` // 缓冲助力数量
	Level          int `json:"level"`           // 服务器等级
}
 