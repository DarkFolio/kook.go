package kook

import (
	"context"
	"encoding/json"
	"fmt"
)

// GameService 游戏/动态相关API服务
type GameService struct {
	client *Client
}

// GetGameList 获取游戏列表
func (s *GameService) GetGameList(ctx context.Context, gameType string) (*ListGamesResponse, error) {
	query := make(map[string]string)
	if gameType != "" {
		query["type"] = gameType
	}

	resp, err := s.client.Get(ctx, "game", query)
	if err != nil {
		return nil, err
	}

	var result ListGamesResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("解析游戏列表失败: %w", err)
	}

	return &result, nil
}

// CreateGame 添加游戏
func (s *GameService) CreateGame(ctx context.Context, name, icon string) (*Game, error) {
	if name == "" {
		return nil, fmt.Errorf("游戏名称不能为空")
	}

	params := map[string]interface{}{
		"name": name,
	}

	if icon != "" {
		params["icon"] = icon
	}

	resp, err := s.client.Post(ctx, "game/create", params)
	if err != nil {
		return nil, err
	}

	var game Game
	if err := json.Unmarshal(resp.Data, &game); err != nil {
		return nil, fmt.Errorf("解析游戏信息失败: %w", err)
	}

	return &game, nil
}

// UpdateGame 更新游戏
func (s *GameService) UpdateGame(ctx context.Context, id int, name, icon string) (*Game, error) {
	if id <= 0 {
		return nil, fmt.Errorf("游戏ID不能为空")
	}

	params := map[string]interface{}{
		"id": id,
	}

	if name != "" {
		params["name"] = name
	}
	if icon != "" {
		params["icon"] = icon
	}

	resp, err := s.client.Post(ctx, "game/update", params)
	if err != nil {
		return nil, err
	}

	var game Game
	if err := json.Unmarshal(resp.Data, &game); err != nil {
		return nil, fmt.Errorf("解析游戏信息失败: %w", err)
	}

	return &game, nil
}

// DeleteGame 删除游戏
func (s *GameService) DeleteGame(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("游戏ID不能为空")
	}

	params := map[string]interface{}{
		"id": id,
	}

	_, err := s.client.Post(ctx, "game/delete", params)
	return err
}

// AddGameActivity 添加游戏活动记录（开始玩游戏）
func (s *GameService) AddGameActivity(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("游戏ID不能为空")
	}

	params := map[string]interface{}{
		"id":        id,
		"data_type": 1, // 1表示游戏
	}

	_, err := s.client.Post(ctx, "game/activity", params)
	return err
}

// AddMusicActivity 添加音乐活动记录（开始听音乐）
func (s *GameService) AddMusicActivity(ctx context.Context, params MusicActivityParams) error {
	if params.Singer == "" {
		return fmt.Errorf("歌手名不能为空")
	}
	if params.MusicName == "" {
		return fmt.Errorf("歌曲名不能为空")
	}

	requestParams := map[string]interface{}{
		"data_type":  2, // 2表示音乐
		"singer":     params.Singer,
		"music_name": params.MusicName,
	}

	if params.Software != "" {
		requestParams["software"] = params.Software
	} else {
		requestParams["software"] = "cloudmusic" // 默认网易云音乐
	}

	_, err := s.client.Post(ctx, "game/activity", requestParams)
	return err
}

// DeleteActivity 删除活动记录（结束玩游戏/听音乐）
func (s *GameService) DeleteActivity(ctx context.Context, dataType int) error {
	if dataType != 1 && dataType != 2 {
		return fmt.Errorf("数据类型必须为1（游戏）或2（音乐）")
	}

	params := map[string]interface{}{
		"data_type": dataType,
	}

	_, err := s.client.Post(ctx, "game/delete-activity", params)
	return err
}

// DeleteGameActivity 删除游戏活动记录（结束玩游戏）
func (s *GameService) DeleteGameActivity(ctx context.Context) error {
	return s.DeleteActivity(ctx, 1)
}

// DeleteMusicActivity 删除音乐活动记录（结束听音乐）
func (s *GameService) DeleteMusicActivity(ctx context.Context) error {
	return s.DeleteActivity(ctx, 2)
}

// 数据结构定义

// Game 游戏信息
type Game struct {
	ID          int      `json:"id"`           // 游戏ID
	Name        string   `json:"name"`         // 游戏名称
	Type        int      `json:"type"`         // 游戏类型：0游戏，1VUP，2进程
	Options     string   `json:"options"`      // 进程额外信息
	KmhookAdmin bool     `json:"kmhook_admin"` // 是否以管理员权限启动KOOK
	ProcessName []string `json:"process_name"` // 进程名称列表
	ProductName []string `json:"product_name"` // 产品名称列表
	Icon        string   `json:"icon"`         // 游戏图标URL
}

// MusicActivityParams 音乐活动参数
type MusicActivityParams struct {
	Software  string `json:"software"`   // 软件名：cloudmusic, qqmusic, kugou
	Singer    string `json:"singer"`     // 歌手名
	MusicName string `json:"music_name"` // 歌曲名
}

// ListGamesResponse 游戏列表响应
type ListGamesResponse struct {
	Items []Game         `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// 游戏类型常量
const (
	GameTypeAll    = "0" // 全部
	GameTypeUser   = "1" // 用户创建
	GameTypeSystem = "2" // 系统创建
)

// 音乐软件常量
const (
	SoftwareCloudMusic = "cloudmusic" // 网易云音乐
	SoftwareQQMusic    = "qqmusic"    // QQ音乐
	SoftwareKugou      = "kugou"      // 酷狗音乐
) 