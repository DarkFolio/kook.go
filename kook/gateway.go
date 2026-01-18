package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// GatewayService 网关相关API服务
type GatewayService struct {
	client *Client
}

// GetGateway 获取网关连接信息
func (s *GatewayService) GetGateway(ctx context.Context, compress int) (*Gateway, error) {
	query := make(map[string]string)
	if compress >= 0 {
		query["compress"] = fmt.Sprintf("%d", compress)
	}

	resp, err := s.client.Get(ctx, "gateway/index", query)
	if err != nil {
		return nil, err
	}

	var gateway Gateway
	if err := json.Unmarshal(resp.Data, &gateway); err != nil {
		return nil, fmt.Errorf("解析网关信息失败: %w", err)
	}

	return &gateway, nil
}

// GetVoiceGateway 获取语音网关连接信息
func (s *GatewayService) GetVoiceGateway(ctx context.Context, channelID string) (*VoiceGateway, error) {
	if channelID == "" {
		return nil, fmt.Errorf("频道ID不能为空")
	}

	query := map[string]string{
		"channel_id": channelID,
	}

	resp, err := s.client.Get(ctx, "gateway/voice", query)
	if err != nil {
		return nil, err
	}

	var voiceGateway VoiceGateway
	if err := json.Unmarshal(resp.Data, &voiceGateway); err != nil {
		return nil, fmt.Errorf("解析语音网关信息失败: %w", err)
	}

	return &voiceGateway, nil
}
 