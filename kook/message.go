package kook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// MessageService 消息相关API服务
type MessageService struct {
	client *Client
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(ctx context.Context, params SendMessageParams) (*Message, error) {
	var endpoint string
	requestParams := make(map[string]interface{})

	// 根据消息类型选择端点
	if params.TargetID == "" {
		return nil, fmt.Errorf("目标ID不能为空")
	}

	// 根据类型判断是私聊还是频道消息
	if params.Type == "private" {
		endpoint = "direct-message/create"
		requestParams["target_id"] = params.TargetID
	} else {
		endpoint = "message/create"
		requestParams["target_id"] = params.TargetID
	}

	// 设置消息内容和类型
	if params.Content == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}
	requestParams["content"] = params.Content

	if params.MsgType > 0 {
		requestParams["type"] = params.MsgType
	} else {
		requestParams["type"] = 1 // 默认文本消息
	}

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

	resp, err := s.client.Post(ctx, endpoint, requestParams)
	if err != nil {
		return nil, err
	}

	var message Message
	if err := json.Unmarshal(resp.Data, &message); err != nil {
		return nil, fmt.Errorf("解析消息失败: %w", err)
	}

	return &message, nil
}

// GetMessageList 获取消息列表
func (s *MessageService) GetMessageList(ctx context.Context, targetID string, params GetMessageListParams) (*ListMessagesResponse, error) {
	if targetID == "" {
		return nil, fmt.Errorf("目标ID不能为空")
	}

	var endpoint string
	query := map[string]string{
		"target_id": targetID,
	}

	// 根据类型选择端点
	if params.Type == "private" {
		endpoint = "direct-message/list"
	} else {
		endpoint = "message/list"
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

	var message Message
	if err := json.Unmarshal(resp.Data, &message); err != nil {
		return nil, fmt.Errorf("解析消息失败: %w", err)
	}

	return &message, nil
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
	Type         string `json:"type,omitempty"`          // 消息类型：private, channel
	TargetID     string `json:"target_id"`               // 目标ID（频道ID或用户ID）
	Content      string `json:"content"`                 // 消息内容
	MsgType      int    `json:"msg_type,omitempty"`      // 消息类型（1文本，2图片等）
	Quote        string `json:"quote,omitempty"`         // 引用消息ID
	Nonce        string `json:"nonce,omitempty"`         // 随机字符串，防重复
	TempTargetID string `json:"temp_target_id,omitempty"` // 临时目标ID
}

// GetMessageListParams 获取消息列表参数
type GetMessageListParams struct {
	Type     string `json:"type,omitempty"`      // 消息类型：private, channel
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
func (s *MessageService) PinMessage(ctx context.Context, msgID string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
	}

	_, err := s.client.Post(ctx, "message/pin", params)
	return err
}

// UnpinMessage 取消置顶消息
func (s *MessageService) UnpinMessage(ctx context.Context, msgID string) error {
	if msgID == "" {
		return fmt.Errorf("消息ID不能为空")
	}

	params := map[string]interface{}{
		"msg_id": msgID,
	}

	_, err := s.client.Post(ctx, "message/unpin", params)
	return err
}
 