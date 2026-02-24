package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// MessageService 消息相关API服务
type MessageService struct {
	client *Client
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(ctx context.Context, params SendMessageParams) (*Message, error) {
	var endpoint string
	requestParams := make(map[string]interface{})

	scope, err := normalizeMessageScope(params.Type)
	if err != nil {
		return nil, err
	}

	// 根据类型判断是私聊还是频道消息
	if scope == "private" {
		endpoint = "direct-message/create"
		if params.TargetID == "" && params.ChatCode == "" {
			return nil, fmt.Errorf("私聊消息必须提供目标ID或会话Code")
		}
		if params.TargetID != "" {
			requestParams["target_id"] = params.TargetID
		}
		if params.ChatCode != "" {
			requestParams["chat_code"] = params.ChatCode
		}
	} else {
		endpoint = "message/create"
		if params.TargetID == "" {
			return nil, fmt.Errorf("频道消息目标ID不能为空")
		}
		requestParams["target_id"] = params.TargetID
	}

	// 设置消息内容和类型
	if params.Content == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}
	requestParams["content"] = params.Content

	msgType := params.MsgType
	if msgType <= 0 {
		if scope == "private" {
			msgType = MessageTypeText
		} else {
			// 频道消息官方默认是 9 (kmarkdown)
			msgType = MessageTypeKMD
		}
	}
	if msgType == MessageTypeCard {
		if err := validateCardContent(params.Content); err != nil {
			return nil, err
		}
	}
	requestParams["type"] = msgType

	// 设置可选参数
	if params.Quote != "" {
		requestParams["quote"] = params.Quote
	}
	if params.Nonce != "" {
		requestParams["nonce"] = params.Nonce
	}
	if params.TempTargetID != "" {
		requestParams["temp_target_id"] = params.TempTargetID
	}
	if params.TemplateID != "" {
		requestParams["template_id"] = params.TemplateID
	}
	if params.ReplyMsgID != "" {
		requestParams["reply_msg_id"] = params.ReplyMsgID
	}

	resp, err := s.client.Post(ctx, endpoint, requestParams)
	if err != nil {
		return nil, err
	}

	var created struct {
		MsgID        string `json:"msg_id"`
		MsgTimestamp int64  `json:"msg_timestamp"`
		Nonce        string `json:"nonce"`
	}
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return nil, fmt.Errorf("解析消息失败: %w", err)
	}

	return &Message{
		ID:       created.MsgID,
		Type:     msgType,
		Content:  params.Content,
		CreateAt: created.MsgTimestamp,
	}, nil
}

// GetMessageList 获取消息列表
func (s *MessageService) GetMessageList(ctx context.Context, targetID string, params GetMessageListParams) (*ListMessagesResponse, error) {
	var endpoint string
	query := make(map[string]string)

	scope, err := normalizeMessageScope(params.Type)
	if err != nil {
		return nil, err
	}

	// 根据类型选择端点
	if scope == "private" {
		endpoint = "direct-message/list"
		if targetID == "" && params.ChatCode == "" {
			return nil, fmt.Errorf("私聊消息列表必须提供目标ID或会话Code")
		}
		if targetID != "" {
			query["target_id"] = targetID
		}
		if params.ChatCode != "" {
			query["chat_code"] = params.ChatCode
		}
	} else {
		endpoint = "message/list"
		if targetID == "" {
			return nil, fmt.Errorf("频道消息列表目标ID不能为空")
		}
		query["target_id"] = targetID
	}

	// 添加查询参数
	if params.MsgID != "" {
		query["msg_id"] = params.MsgID
	}
	if params.Pin > 0 {
		query["pin"] = strconv.Itoa(params.Pin)
	}
	if params.Flag != "" {
		query["flag"] = params.Flag
	}
	if params.PageSize > 0 && params.PageSize <= 100 {
		query["page_size"] = strconv.Itoa(params.PageSize)
	}

	resp, err := s.client.Get(ctx, endpoint, query)
	if err != nil {
		return nil, err
	}

	var result ListMessagesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析消息列表失败: %w", err)
	}

	return &result, nil
}

// GetMessage 获取消息详情
func (s *MessageService) GetMessage(ctx context.Context, msgID string) (*Message, error) {
	if msgID == "" {
		return nil, fmt.Errorf("消息ID不能为空")
	}

	query := map[string]string{
		"msg_id": msgID,
	}

	resp, err := s.client.Get(ctx, "message/view", query)
	if err != nil {
		return nil, err
	}

	var message Message
	if err := json.Unmarshal(resp.Data, &message); err != nil {
		return nil, fmt.Errorf("解析消息失败: %w", err)
	}

	return &message, nil
}

// GetDirectMessage 获取私聊消息详情
func (s *MessageService) GetDirectMessage(ctx context.Context, chatCode, msgID string) (*Message, error) {
	if chatCode == "" {
		return nil, fmt.Errorf("私聊会话Code不能为空")
	}
	if msgID == "" {
		return nil, fmt.Errorf("消息ID不能为空")
	}

	query := map[string]string{
		"chat_code": chatCode,
		"msg_id":    msgID,
	}

	resp, err := s.client.Get(ctx, "direct-message/view", query)
	if err != nil {
		return nil, err
	}

	var message Message
	if err := json.Unmarshal(resp.Data, &message); err != nil {
		return nil, fmt.Errorf("解析私聊消息失败: %w", err)
	}

	return &message, nil
}

// UpdateMessage 更新消息
func (s *MessageService) UpdateMessage(ctx context.Context, msgID, content string, quote string, tempTargetID string) (*Message, error) {
	if msgID == "" {
		return nil, fmt.Errorf("消息ID不能为空")
	}
	if content == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}

	params := map[string]interface{}{
		"msg_id":  msgID,
		"content": content,
	}

	if quote != "" {
		params["quote"] = quote
	}
	if tempTargetID != "" {
		params["temp_target_id"] = tempTargetID
	}

	resp, err := s.client.Post(ctx, "message/update", params)
	if err != nil {
		return nil, err
	}

	// 官方接口通常返回空数组 []，这里兼容两种返回体
	var message Message
	if len(resp.Data) > 0 && string(resp.Data) != "[]" && string(resp.Data) != "null" {
		if err := json.Unmarshal(resp.Data, &message); err != nil {
			return nil, fmt.Errorf("解析消息失败: %w", err)
		}
		return &message, nil
	}

	return &Message{
		ID:      msgID,
		Content: content,
	}, nil
}

// UpdateDirectMessage 更新私聊消息（支持 KMarkdown 与 CardMessage）
func (s *MessageService) UpdateDirectMessage(ctx context.Context, msgID, content, quote string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if content == "" {
		return fmt.Errorf("消息内容不能为空")
	}

	params := map[string]interface{}{
		"msg_id":  msgID,
		"content": content,
	}
	if quote != "" {
		params["quote"] = quote
	}

	_, err := s.client.Post(ctx, "direct-message/update", params)
	return err
}

// DeleteMessage 删除消息
func (s *MessageService) DeleteMessage(ctx context.Context, msgID string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
	}

	_, err := s.client.Post(ctx, "message/delete", params)
	return err
}

// DeleteDirectMessage 删除私聊消息
func (s *MessageService) DeleteDirectMessage(ctx context.Context, msgID string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
	}

	_, err := s.client.Post(ctx, "direct-message/delete", params)
	return err
}

// AddReaction 添加回应
func (s *MessageService) AddReaction(ctx context.Context, msgID, emoji string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if emoji == "" {
		return fmt.Errorf("表情不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
		"emoji":  emoji,
	}

	_, err := s.client.Post(ctx, "message/add-reaction", params)
	return err
}

// AddDirectReaction 为私聊消息添加回应
func (s *MessageService) AddDirectReaction(ctx context.Context, msgID, emoji string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if emoji == "" {
		return fmt.Errorf("表情不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
		"emoji":  emoji,
	}

	_, err := s.client.Post(ctx, "direct-message/add-reaction", params)
	return err
}

// DeleteReaction 删除回应
func (s *MessageService) DeleteReaction(ctx context.Context, msgID, emoji, userID string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if emoji == "" {
		return fmt.Errorf("表情不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
		"emoji":  emoji,
	}

	if userID != "" {
		params["user_id"] = userID
	}

	_, err := s.client.Post(ctx, "message/delete-reaction", params)
	return err
}

// DeleteDirectReaction 删除私聊消息回应
func (s *MessageService) DeleteDirectReaction(ctx context.Context, msgID, emoji, userID string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if emoji == "" {
		return fmt.Errorf("表情不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
		"emoji":  emoji,
	}
	if userID != "" {
		params["user_id"] = userID
	}

	_, err := s.client.Post(ctx, "direct-message/delete-reaction", params)
	return err
}

// GetReactionUserList 获取回应用户列表
func (s *MessageService) GetReactionUserList(ctx context.Context, msgID, emoji string) ([]User, error) {
	if msgID == "" {
		return nil, fmt.Errorf("消息ID不能为空")
	}
	if emoji == "" {
		return nil, fmt.Errorf("表情不能为空")
	}

	query := map[string]string{
		"msg_id": msgID,
		"emoji":  emoji,
	}

	resp, err := s.client.Get(ctx, "message/reaction-list", query)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		return nil, fmt.Errorf("解析用户列表失败: %w", err)
	}

	return users, nil
}

// GetDirectReactionUserList 获取私聊消息回应用户列表
func (s *MessageService) GetDirectReactionUserList(ctx context.Context, msgID, emoji string) ([]User, error) {
	if msgID == "" {
		return nil, fmt.Errorf("消息ID不能为空")
	}
	if emoji == "" {
		return nil, fmt.Errorf("表情不能为空")
	}

	query := map[string]string{
		"msg_id": msgID,
		"emoji":  emoji,
	}

	resp, err := s.client.Get(ctx, "direct-message/reaction-list", query)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		return nil, fmt.Errorf("解析用户列表失败: %w", err)
	}

	return users, nil
}

// CheckCard 检查卡片消息格式
func (s *MessageService) CheckCard(ctx context.Context, content string) (*CheckCardResponse, error) {
	if content == "" {
		return nil, fmt.Errorf("卡片内容不能为空")
	}

	params := map[string]interface{}{
		"content": content,
	}

	resp, err := s.client.Post(ctx, "message/check-card", params)
	if err != nil {
		return nil, err
	}

	var result CheckCardResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析检查结果失败: %w", err)
	}

	return &result, nil
}

// SendMessageParams 发送消息参数
type SendMessageParams struct {
	Type         string `json:"type,omitempty"`           // 作用域：private 或 channel（默认channel）
	TargetID     string `json:"target_id,omitempty"`      // 目标ID（频道ID或用户ID）
	ChatCode     string `json:"chat_code,omitempty"`      // 私聊会话Code（私聊可选，和TargetID二选一）
	Content      string `json:"content"`                  // 消息内容
	MsgType      int    `json:"msg_type,omitempty"`       // 消息类型（1文本，9KMarkdown，10卡片）
	Quote        string `json:"quote,omitempty"`          // 引用消息ID
	Nonce        string `json:"nonce,omitempty"`          // 随机字符串，防重复
	TempTargetID string `json:"temp_target_id,omitempty"` // 临时目标ID
	TemplateID   string `json:"template_id,omitempty"`    // 模板ID
	ReplyMsgID   string `json:"reply_msg_id,omitempty"`   // 回复消息ID（用于配额折扣）
}

// SendCardMessage 发送卡片消息
func (s *MessageService) SendCardMessage(ctx context.Context, params SendMessageParams) (*Message, error) {
	if err := validateCardContent(params.Content); err != nil {
		return nil, err
	}
	params.MsgType = MessageTypeCard
	return s.SendMessage(ctx, params)
}

func normalizeMessageScope(scope string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "", "channel", "guild":
		return "channel", nil
	case "private", "direct", "dm":
		return "private", nil
	default:
		return "", fmt.Errorf("无效消息作用域: %s（可选: channel/private）", scope)
	}
}

func validateCardContent(content string) error {
	var cards []json.RawMessage
	if err := json.Unmarshal([]byte(content), &cards); err != nil {
		return fmt.Errorf("卡片消息 content 必须是 JSON 数组字符串: %w", err)
	}
	if len(cards) == 0 {
		return fmt.Errorf("卡片消息至少包含一个 card")
	}
	return nil
}

// GetMessageListParams 获取消息列表参数
type GetMessageListParams struct {
	Type     string `json:"type,omitempty"`      // 消息类型：private, channel
	ChatCode string `json:"chat_code,omitempty"` // 私聊会话Code（私聊可选）
	MsgID    string `json:"msg_id,omitempty"`    // 参考消息ID
	Pin      int    `json:"pin,omitempty"`       // 只看置顶消息：0否，1是
	Flag     string `json:"flag,omitempty"`      // 查询模式：before, around, after
	PageSize int    `json:"page_size,omitempty"` // 返回数量，默认50，最大100
}

// ListMessagesResponse 消息列表响应
type ListMessagesResponse struct {
	Items []Message `json:"items"`
}

// CheckCardResponse 检查卡片响应
type CheckCardResponse struct {
	Mention struct {
		Mentions     []string      `json:"mentions"`
		MentionRoles []string      `json:"mentionRoles"`
		MentionAll   bool          `json:"mentionAll"`
		MentionHere  bool          `json:"mentionHere"`
		MentionPart  []interface{} `json:"mentionPart"`
		NavChannels  []interface{} `json:"navChannels"`
		ChannelPart  []interface{} `json:"channelPart"`
		GuildEmojis  []interface{} `json:"guildEmojis"`
	} `json:"mention"`
	Content string `json:"content"`
}

// PinMessage 置顶消息
func (s *MessageService) PinMessage(ctx context.Context, msgID string, targetID ...string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if len(targetID) == 0 || strings.TrimSpace(targetID[0]) == "" {
		return fmt.Errorf("置顶消息需要目标ID（频道ID或私聊目标ID）")
	}

	params := map[string]interface{}{
		"msg_id":    msgID,
		"target_id": strings.TrimSpace(targetID[0]),
	}

	_, err := s.client.Post(ctx, "message/pin", params)
	return err
}

// UnpinMessage 取消置顶消息
func (s *MessageService) UnpinMessage(ctx context.Context, msgID string, targetID ...string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}
	if len(targetID) == 0 || strings.TrimSpace(targetID[0]) == "" {
		return fmt.Errorf("取消置顶消息需要目标ID（频道ID或私聊目标ID）")
	}

	params := map[string]interface{}{
		"msg_id":    msgID,
		"target_id": strings.TrimSpace(targetID[0]),
	}

	_, err := s.client.Post(ctx, "message/unpin", params)
	return err
}
