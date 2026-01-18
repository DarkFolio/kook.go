package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// RoleService 角色相关API服务
type RoleService struct {
	client *Client
}

// GetRoleList 获取服务器角色列表
func (s *RoleService) GetRoleList(ctx context.Context, guildID string, page, pageSize int) (*ListRolesResponse, error) {
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

	resp, err := s.client.Get(ctx, "guild-role/list", query)
	if err != nil {
		return nil, err
	}

	var result ListRolesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析角色列表失败: %w", err)
	}

	return &result, nil
}

// CreateRole 创建服务器角色
func (s *RoleService) CreateRole(ctx context.Context, guildID string, name string) (*GuildRole, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
	}

	if name != "" {
		params["name"] = name
	}

	resp, err := s.client.Post(ctx, "guild-role/create", params)
	if err != nil {
		return nil, err
	}

	var roles []GuildRole
	if err := json.Unmarshal(resp.Data, &roles); err != nil {
		return nil, fmt.Errorf("解析角色信息失败: %w", err)
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("未返回角色信息")
	}

	return &roles[0], nil
}

// UpdateRole 更新服务器角色
func (s *RoleService) UpdateRole(ctx context.Context, guildID string, roleID int, params UpdateRoleParams) (*GuildRole, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}
	if roleID <= 0 {
		return nil, fmt.Errorf("角色ID不能为空")
	}

	requestParams := map[string]interface{}{
		"guild_id": guildID,
		"role_id":  roleID,
	}

	if params.Name != "" {
		requestParams["name"] = params.Name
	}
	if params.Color >= 0 {
		requestParams["color"] = params.Color
	}
	if params.Hoist >= 0 {
		requestParams["hoist"] = params.Hoist
	}
	if params.Mentionable >= 0 {
		requestParams["mentionable"] = params.Mentionable
	}
	if params.Permissions >= 0 {
		requestParams["permissions"] = params.Permissions
	}

	resp, err := s.client.Post(ctx, "guild-role/update", requestParams)
	if err != nil {
		return nil, err
	}

	var roles []GuildRole
	if err := json.Unmarshal(resp.Data, &roles); err != nil {
		return nil, fmt.Errorf("解析角色信息失败: %w", err)
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("未返回角色信息")
	}

	return &roles[0], nil
}

// DeleteRole 删除服务器角色
func (s *RoleService) DeleteRole(ctx context.Context, guildID string, roleID int) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if roleID <= 0 {
		return fmt.Errorf("角色ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"role_id":  roleID,
	}

	_, err := s.client.Post(ctx, "guild-role/delete", params)
	return err
}

// GrantRole 赋予用户角色
func (s *RoleService) GrantRole(ctx context.Context, guildID, userID string, roleID int) (*UserRoleResponse, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if roleID <= 0 {
		return nil, fmt.Errorf("角色ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"user_id":  userID,
		"role_id":  roleID,
	}

	resp, err := s.client.Post(ctx, "guild-role/grant", params)
	if err != nil {
		return nil, err
	}

	var result UserRoleResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析用户角色信息失败: %w", err)
	}

	return &result, nil
}

// RevokeRole 删除用户角色
func (s *RoleService) RevokeRole(ctx context.Context, guildID, userID string, roleID int) (*UserRoleResponse, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if roleID <= 0 {
		return nil, fmt.Errorf("角色ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"user_id":  userID,
		"role_id":  roleID,
	}

	resp, err := s.client.Post(ctx, "guild-role/revoke", params)
	if err != nil {
		return nil, err
	}

	var result UserRoleResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析用户角色信息失败: %w", err)
	}

	return &result, nil
}

// 数据结构定义

// GuildRole 服务器角色信息
type GuildRole struct {
	RoleID      int `json:"role_id"`      // 角色ID
	Name        string `json:"name"`         // 角色名称
	Color       int `json:"color"`        // 角色色值
	Position    int `json:"position"`     // 角色位置
	Hoist       int `json:"hoist"`        // 是否在用户列表排到前面
	Mentionable int `json:"mentionable"`  // 是否可以被提及
	Permissions int `json:"permissions"`  // 权限值
}

// UpdateRoleParams 更新角色参数
type UpdateRoleParams struct {
	Name        string `json:"name,omitempty"`        // 角色名称
	Color       int    `json:"color,omitempty"`       // 角色色值
	Hoist       int    `json:"hoist,omitempty"`       // 是否在用户列表排到前面
	Mentionable int    `json:"mentionable,omitempty"` // 是否可以被提及
	Permissions int    `json:"permissions,omitempty"` // 权限值
}

// ListRolesResponse 角色列表响应
type ListRolesResponse struct {
	Items []GuildRole    `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// UserRoleResponse 用户角色响应
type UserRoleResponse struct {
	UserID  string `json:"user_id"`  // 用户ID
	GuildID string `json:"guild_id"` // 服务器ID
	Roles   []int  `json:"roles"`    // 角色ID列表
} 