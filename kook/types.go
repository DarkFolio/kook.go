package kook

import (
	"encoding/json"
	"time"
)

// User 用户信息
type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	IdentifyNum    string `json:"identify_num"`
	Online         bool   `json:"online"`
	Bot            bool   `json:"bot"`
	Status         int    `json:"status"`
	Avatar         string `json:"avatar"`
	VipAvatar      string `json:"vip_avatar"`
	Banner         string `json:"banner"`
	Nickname       string `json:"nickname"`
	Roles          []int  `json:"roles"`
	IsVip          bool   `json:"is_vip"`
	VipAmp         bool   `json:"vip_amp"`
	InvitedCount   int    `json:"invited_count"`
	TagInfo        TagInfo `json:"tag_info"`
}

// TagInfo 标签信息
type TagInfo struct {
	Color string `json:"color"`
	Text  string `json:"text"`
}

// Guild 服务器信息
type Guild struct {
	ID                      string          `json:"id"`
	Name                    string          `json:"name"`
	Topic                   string          `json:"topic"`
	UserID                  string          `json:"user_id"`
	Icon                    string          `json:"icon"`
	NotifyType              int             `json:"notify_type"`
	Region                  string          `json:"region"`
	EnableOpen              bool            `json:"enable_open"`
	OpenID                  string          `json:"open_id"`
	DefaultChannelID        string          `json:"default_channel_id"`
	WelcomeChannelID        string          `json:"welcome_channel_id"`
	Roles                   []Role          `json:"roles"`
	Channels                []Channel       `json:"channels"`
	MaxPersons              int             `json:"max_persons"`
	Level                   int             `json:"level"`
	BoostNum                int             `json:"boost_num"`
	BufferBoostNum          int             `json:"buffer_boost_num"`
	Banner                  string          `json:"banner"`
	Features                []GuildFeature  `json:"features"`
	Emojis                  []Emoji         `json:"emojis"`
}

// GuildFeature 服务器功能特性
type GuildFeature struct {
	Feature     string `json:"feature"`
	Description string `json:"description"`
}

// Role 角色信息
type Role struct {
	RoleID      int    `json:"role_id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Hoist       int    `json:"hoist"`
	Mentionable int    `json:"mentionable"`
	Permissions int    `json:"permissions"`
}

// Channel 频道信息
type Channel struct {
	ID                     string                 `json:"id"`
	Name                   string                 `json:"name"`
	UserID                 string                 `json:"user_id"`
	GuildID                string                 `json:"guild_id"`
	Topic                  string                 `json:"topic"`
	IsCategory             bool                   `json:"is_category"`
	ParentID               string                 `json:"parent_id"`
	Level                  int                    `json:"level"`
	SlowMode               int                    `json:"slow_mode"`
	Type                   int                    `json:"type"`
	PermissionOverwrites   []PermissionOverwrite  `json:"permission_overwrites"`
	PermissionUsers        []PermissionUser       `json:"permission_users"`
	PermissionSync         int                    `json:"permission_sync"`
	HasPassword            bool                   `json:"has_password"`
	LimitAmount            int                    `json:"limit_amount"`
	VoiceQuality           int                    `json:"voice_quality"`
}

// PermissionOverwrite 权限覆写
type PermissionOverwrite struct {
	RoleID int `json:"role_id"`
	Allow  int `json:"allow"`
	Deny   int `json:"deny"`
}

// PermissionUser 用户权限
type PermissionUser struct {
	User  User `json:"user"`
	Allow int  `json:"allow"`
	Deny  int  `json:"deny"`
}

// Message 消息信息
type Message struct {
	ID               string        `json:"id"`
	Type             int           `json:"type"`
	Content          string        `json:"content"`
	Mention          []string      `json:"mention"`
	MentionAll       bool          `json:"mention_all"`
	MentionRoles     []string      `json:"mention_roles"`
	MentionHere      bool          `json:"mention_here"`
	Embeds           []interface{} `json:"embeds"`
	Attachments      []Attachment  `json:"attachments"`
	CreateAt         int64         `json:"create_at"`
	UpdatedAt        int64         `json:"updated_at"`
	Reactions        []Reaction    `json:"reactions"`
	Author           User          `json:"author"`
	ImageName        string        `json:"image_name"`
	ReadStatus       bool          `json:"read_status"`
	Quote            *Quote        `json:"quote"`
	MentionInfo      MentionInfo   `json:"mention_info"`
}

// Attachment 附件信息
type Attachment struct {
	Type string `json:"type"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

// Reaction 反应信息
type Reaction struct {
	Emoji Emoji   `json:"emoji"`
	Count int     `json:"count"`
	Me    bool    `json:"me"`
}

// Quote 引用消息
type Quote struct {
	ID               string        `json:"id"`
	Type             int           `json:"type"`
	Content          string        `json:"content"`
	CreateAt         int64         `json:"create_at"`
	Author           User          `json:"author"`
	RonCreate        bool          `json:"ron_create"`
}

// MentionInfo 提及信息
type MentionInfo struct {
	MentionPart    []MentionPart `json:"mention_part"`
	MentionRolePart []MentionRolePart `json:"mention_role_part"`
}

// MentionPart 提及用户信息
type MentionPart struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}

// MentionRolePart 提及角色信息
type MentionRolePart struct {
	RoleID int    `json:"role_id"`
	Name   string `json:"name"`
	Color  int    `json:"color"`
}



// GuildMember 服务器成员信息
type GuildMember struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Nickname       string    `json:"nickname"`
	IdentifyNum    string    `json:"identify_num"`
	Online         bool      `json:"online"`
	Bot            bool      `json:"bot"`
	Status         int       `json:"status"`
	Avatar         string    `json:"avatar"`
	VipAvatar      string    `json:"vip_avatar"`
	Roles          []int     `json:"roles"`
	JoinedAt       int64     `json:"joined_at"`
	ActiveTime     int64     `json:"active_time"`
	IsVip          bool      `json:"is_vip"`
	VipAmp         bool      `json:"vip_amp"`
}

// Gateway 网关信息
type Gateway struct {
	URL string `json:"url"`
}

// VoiceGateway 语音网关信息
type VoiceGateway struct {
	GatewayURL   string `json:"gateway_url"`
	IosVoiceSDK  int    `json:"ios_voice_sdk"`
	PCVoiceSDK   int    `json:"pc_voice_sdk"`
}

// Event 事件信息
type Event struct {
	ChannelType string      `json:"channel_type"`
	Type        int         `json:"type"`
	TargetID    string      `json:"target_id"`
	AuthorID    string      `json:"author_id"`
	Content     string      `json:"content"`
	MsgID       string      `json:"msg_id"`
	MsgTimestamp int64      `json:"msg_timestamp"`
	Nonce       string      `json:"nonce"`
	Extra       interface{} `json:"extra"`
}





// PaginationMeta 分页信息
type PaginationMeta struct {
	Page      int `json:"page"`
	PageTotal int `json:"page_total"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
}

// ListResponse 列表响应
type ListResponse struct {
	Items []interface{}  `json:"items"`
	Meta  PaginationMeta `json:"meta"`
	Sort  map[string]int `json:"sort"`
}

// Time 自定义时间类型，用于处理KOOK API的时间戳
type Time struct {
	time.Time
}

// UnmarshalJSON 实现JSON反序列化
func (t *Time) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t.Time = time.Unix(timestamp/1000, (timestamp%1000)*1000000)
	return nil
}

// MarshalJSON 实现JSON序列化
func (t Time) MarshalJSON() ([]byte, error) {
	timestamp := t.Unix()*1000 + int64(t.Nanosecond()/1000000)
	return json.Marshal(timestamp)
} 