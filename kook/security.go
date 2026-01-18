package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// SecurityService 安全设置相关API服务
type SecurityService struct {
	client *Client
}

// GetSecuritySettings 获取服务器安全设置
func (s *SecurityService) GetSecuritySettings(ctx context.Context, guildID string) (*SecuritySettings, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	resp, err := s.client.Get(ctx, "guild-security/settings", query)
	if err != nil {
		return nil, err
	}

	var settings SecuritySettings
	if err := json.Unmarshal(resp.Data, &settings); err != nil {
		return nil, fmt.Errorf("解析安全设置失败: %w", err)
	}

	return &settings, nil
}

// UpdateSecuritySetting 更新安全设置
func (s *SecurityService) UpdateSecuritySetting(ctx context.Context, guildID, settingID string, enabled bool) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}
	if settingID == "" {
		return fmt.Errorf("设置ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"id":       settingID,
		"switch":   enabled,
	}

	_, err := s.client.Post(ctx, "guild-security/update", params)
	return err
}

// GetVerificationLevel 获取验证等级设置
func (s *SecurityService) GetVerificationLevel(ctx context.Context, guildID string) (*VerificationLevel, error) {
	if guildID == "" {
		return nil, fmt.Errorf("服务器ID不能为空")
	}

	query := map[string]string{
		"guild_id": guildID,
	}

	resp, err := s.client.Get(ctx, "guild/verification-level", query)
	if err != nil {
		return nil, err
	}

	var level VerificationLevel
	if err := json.Unmarshal(resp.Data, &level); err != nil {
		return nil, fmt.Errorf("解析验证等级失败: %w", err)
	}

	return &level, nil
}

// UpdateVerificationLevel 更新验证等级
func (s *SecurityService) UpdateVerificationLevel(ctx context.Context, guildID string, level int) error {
	if guildID == "" {
		return fmt.Errorf("服务器ID不能为空")
	}

	params := map[string]interface{}{
		"guild_id": guildID,
		"level":    level,
	}

	_, err := s.client.Post(ctx, "guild/verification-level", params)
	return err
}

// 数据结构定义

// SecuritySettings 安全设置
type SecuritySettings struct {
	GuildID  string           `json:"guild_id"`  // 服务器ID
	Settings []SecurityRule   `json:"settings"`  // 安全规则列表
}

// SecurityRule 安全规则
type SecurityRule struct {
	ID          string `json:"id"`          // 规则ID
	Name        string `json:"name"`        // 规则名称
	Description string `json:"description"` // 规则描述
	Enabled     bool   `json:"enabled"`     // 是否启用
	Type        int    `json:"type"`        // 规则类型
}

// VerificationLevel 验证等级
type VerificationLevel struct {
	GuildID string `json:"guild_id"` // 服务器ID
	Level   int    `json:"level"`    // 验证等级：0无限制，1低，2中，3高，4极高
}

// 验证等级常量
const (
	VerificationLevelNone   = 0 // 无限制
	VerificationLevelLow    = 1 // 低（需要验证邮箱）
	VerificationLevelMedium = 2 // 中（需要在服务器待满5分钟）
	VerificationLevelHigh   = 3 // 高（需要在KOOK注册超过10分钟）
	VerificationLevelVeryHigh = 4 // 极高（需要绑定手机号）
) 