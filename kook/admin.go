package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// AdminService 管理员相关API服务
type AdminService struct {
	client *Client
}

// GetAuditLog 获取审计日志
func (s *AdminService) GetAuditLog(ctx context.Context, guildID string, userID string, targetID string, actionType int, page, pageSize int) (*AuditLogResponse, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	if userID != "" {
		query["user_id"] = userID
	}
	if targetID != "" {
		query["target_id"] = targetID
	}
	if actionType > 0 {
		query["action_type"] = fmt.Sprintf("%d", actionType)
	}
	if page > 0 {
		query["page"] = fmt.Sprintf("%d", page)
	}
	if pageSize > 0 {
		query["page_size"] = fmt.Sprintf("%d", pageSize)
	}

	resp, err := s.client.Get(ctx, "guild/audit-log", query)
	if err != nil {
		return nil, err
	}

	var result AuditLogResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析审计日志失败: %w", err)
	}

	return &result, nil
}

// BanUser 封禁用户
func (s *AdminService) BanUser(ctx context.Context, guildID, userID string, reason string, delMsgDays int) error {
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

	if reason != "" {
		params["reason"] = reason
	}
	if delMsgDays > 0 {
		params["del_msg_days"] = delMsgDays
	}

	_, err := s.client.Post(ctx, "guild/ban", params)
	return err
}

// UnbanUser 解封用户
func (s *AdminService) UnbanUser(ctx context.Context, guildID, userID string) error {
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

	_, err := s.client.Post(ctx, "guild/unban", params)
	return err
}

// GetBannedUsers 获取被封禁的用户列表
func (s *AdminService) GetBannedUsers(ctx context.Context, guildID string, page, pageSize int) (*BannedUsersResponse, error) {
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

	resp, err := s.client.Get(ctx, "guild/ban-list", query)
	if err != nil {
		return nil, err
	}

	var result BannedUsersResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析封禁用户列表失败: %w", err)
	}

	return &result, nil
}

// 数据结构定义

// AuditLogEntry 审计日志条目
type AuditLogEntry struct {
	ID         string                 `json:"id"`          // 日志ID
	UserID     string                 `json:"user_id"`     // 操作者ID
	User       User                   `json:"user"`        // 操作者信息
	TargetID   string                 `json:"target_id"`   // 目标ID
	ActionType int                    `json:"action_type"` // 操作类型
	Reason     string                 `json:"reason"`      // 原因
	Options    map[string]interface{} `json:"options"`     // 选项
	CreatedAt  int64                  `json:"created_at"`  // 创建时间
}

// AuditLogResponse 审计日志响应
type AuditLogResponse struct {
	Items []AuditLogEntry `json:"items"`
	Meta  PaginationMeta  `json:"meta"`
	Sort  map[string]int  `json:"sort"`
}

// BannedUser 被封禁的用户
type BannedUser struct {
	User      User   `json:"user"`       // 用户信息
	Reason    string `json:"reason"`     // 封禁原因
	BannedAt  int64  `json:"banned_at"`  // 封禁时间
	BannedBy  string `json:"banned_by"`  // 封禁者ID
}

// BannedUsersResponse 被封禁用户列表响应
type BannedUsersResponse struct {
	Items []BannedUser   `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// 审计日志操作类型常量
const (
	AuditLogActionGuildUpdate      = 1  // 服务器更新
	AuditLogActionChannelCreate    = 10 // 频道创建
	AuditLogActionChannelUpdate    = 11 // 频道更新
	AuditLogActionChannelDelete    = 12 // 频道删除
	AuditLogActionRoleCreate       = 30 // 角色创建
	AuditLogActionRoleUpdate       = 31 // 角色更新
	AuditLogActionRoleDelete       = 32 // 角色删除
	AuditLogActionMemberKick       = 20 // 踢出成员
	AuditLogActionMemberBan        = 22 // 封禁成员
	AuditLogActionMemberUnban      = 23 // 解封成员
	AuditLogActionMemberUpdate     = 24 // 成员更新
	AuditLogActionMemberRoleUpdate = 25 // 成员角色更新
	AuditLogActionMessageDelete    = 72 // 消息删除
) 