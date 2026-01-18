# KOOK Go SDK

基于 [KOOK 开发者文档](https://developer.kookapp.cn/doc/reference) 的 Go 语言 SDK，提供完整的 API v3 接口封装。

## 相关文档

- [KOOK Go SDK API文档](https://blog.lastwhisper.net/article/16)
- [KOOK API 错误代码文档](https://blog.lastwhisper.net/article/17)

## 功能特性

- **完整的 KOOK API v3 支持** - 覆盖24个主要API服务
- **自动处理 Token 身份验证** - 支持 Bot 和 Bearer 两种认证方式
- **完善的错误处理机制** - 增强的错误类型和详细错误信息
- **WebSocket 实时连接** - 支持自动重连、心跳监控和断线恢复
- **智能速率限制管理** - 全局和端点级别的令牌桶算法
- **自动请求重试机制** - 支持指数退避和可配置重试策略
- **类型安全的 API 调用** - 完整的类型定义和参数验证
- **全面的单元测试** - 核心功能测试覆盖
- **详细的文档和示例代码** - 5个不同场景的完整示例

## 安装

```bash
go get github.com/yourusername/kook-go-sdk
```

## 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/kook-go-sdk/kook"
)

func main() {
    // 使用 Bot Token 创建客户端
    client := kook.NewClient("你的机器人令牌")
    
    // 获取机器人信息
    user, err := client.User.GetMe()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("机器人名称: %s#%s\n", user.Username, user.IdentifyNum)
    fmt.Printf("机器人ID: %s\n", user.ID)
}
```

### 消息操作

```go
// 向频道发送消息
message, err := client.Message.SendMessage(kook.SendMessageParams{
    TargetID: "频道ID",
    Content:  "你好，KOOK！",
    MsgType:  1, // 文本消息
})
if err != nil {
    log.Printf("发送消息失败: %v", err)
    return
}

// 获取消息列表
messages, err := client.Message.GetMessageList("频道ID", kook.GetMessageListParams{
    PageSize: 50,
})
if err != nil {
    log.Printf("获取消息列表失败: %v", err)
    return
}
```

### 服务器和频道管理

```go
// 获取服务器列表
guilds, err := client.Guild.GetGuildList(1, 10, "")
if err != nil {
    log.Printf("获取服务器列表失败: %v", err)
    return
}

// 获取服务器的频道列表
channels, err := client.Channel.GetChannelList("服务器ID", 1, 10, "")
if err != nil {
    log.Printf("获取频道列表失败: %v", err)
    return
}

// 创建新频道
channel, err := client.Channel.CreateChannel(kook.CreateChannelParams{
    GuildID: "服务器ID",
    Name:    "新频道",
    Type:    1, // 文字频道
})
```

### WebSocket 实时事件

```go
// 创建 WebSocket 客户端
wsClient := kook.NewWebSocketClient(client, false) // false = 不压缩

// 注册事件处理器
wsClient.OnEvent(kook.EventTypeTextMessage, func(event *kook.Event) {
    fmt.Printf("收到消息: %s\n", event.Content)
})

wsClient.OnEvent(kook.EventTypeUserJoinedGuild, func(event *kook.Event) {
    fmt.Printf("用户加入服务器: %s\n", event.AuthorID)
})

// 连接到 WebSocket
err := wsClient.Connect()
if err != nil {
    log.Fatal("WebSocket 连接失败:", err)
}

// 保持连接
select {} // 永久阻塞
```

### 角色管理

```go
// 获取角色列表
roles, err := client.Role.GetRoleList("服务器ID", 1, 10)
if err != nil {
    log.Printf("获取角色列表失败: %v", err)
    return
}

// 创建新角色
role, err := client.Role.CreateRole("服务器ID", kook.CreateRoleParams{
    Name:        "新角色",
    Color:       0xFF0000, // 红色
    Permissions: 1024,     // 基础权限
})

// 给用户分配角色
err = client.Role.GrantRole("服务器ID", "用户ID", role.RoleID)
```

### 资源上传

```go
// 上传图片
asset, err := client.Asset.CreateAsset("路径/到/图片.png")
if err != nil {
    log.Printf("上传资源失败: %v", err)
    return
}

// 在消息中使用上传的资源URL
message, err := client.Message.SendMessage(kook.SendMessageParams{
    TargetID: "频道ID",
    Content:  asset.URL,
    MsgType:  2, // 图片消息
})
```

## 高级配置

### 生产环境配置

```go
// 生产环境推荐配置
client := kook.NewClient("你的机器人令牌",
    // 自定义HTTP客户端
    kook.WithHTTPClient(&http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
        },
    }),
    // 自定义重试配置
    kook.WithRetryConfig(&kook.RetryConfig{
        MaxRetries:    5,
        InitialDelay:  1 * time.Second,
        MaxDelay:      30 * time.Second,
        BackoffFactor: 2.0,
    }),
    // 自定义速率限制
    kook.WithRateLimiter(kook.NewGlobalRateLimiter()),
    // 自定义日志器
    kook.WithLogger(customLogger),
)
```

### WebSocket 高级配置

```go
// 创建高可用的WebSocket客户端
wsClient := kook.NewWebSocketClient(client, true) // 启用压缩

// 监控连接状态
go func() {
    for {
        if !wsClient.IsConnected() {
            log.Println("WebSocket连接已断开，正在重连...")
        }
        time.Sleep(30 * time.Second)
    }
}()
```

### 错误处理最佳实践

```go
user, err := client.User.GetMe()
if err != nil {
    if kookErr, ok := kook.IsKOOKError(err); ok {
        switch {
        case kookErr.IsAuthError():
            log.Fatal("认证失败，请检查Token")
        case kookErr.IsRateLimited():
            log.Printf("请求被限流，请等待 %v", kookErr.RetryAfter)
        case kookErr.IsServerError():
            log.Printf("服务器错误，将自动重试")
        default:
            log.Printf("API错误: %s", kookErr.Error())
        }
    } else {
        log.Printf("网络错误: %v", err)
    }
    return
}
```



## API 覆盖范围

此 SDK 提供对 KOOK API v3 的全面覆盖：

### 核心服务
- **UserService**: 用户信息、在线状态管理
- **MessageService**: 消息发送、管理、表情回应
- **GuildService**: 服务器管理、成员操作
- **ChannelService**: 频道管理、权限设置
- **RoleService**: 角色管理、权限分配

### 扩展服务
- **GatewayService**: WebSocket 网关管理
- **GameService**: 游戏状态和活动管理
- **FriendService**: 好友系统操作
- **InviteService**: 邀请管理
- **AssetService**: 媒体上传和管理
- **IntimacyService**: 亲密度系统
- **BadgeService**: 徽章系统
- **BlacklistService**: 黑名单管理
- **EmojiService**: 自定义表情管理
- **RegionService**: 服务器区域信息
- **OAuthService**: OAuth2 认证
- **LiveService**: 直播功能
- **AdminService**: 管理员功能
- **SecurityService**: 安全设置
- **VoiceService**: 语音频道操作
- **ItemService**: 物品系统
- **OrderService**: 订单管理
- **CouponService**: 优惠券系统
- **BoostService**: 服务器助力系统

## 错误处理

SDK 提供完善的错误处理机制：

```go
user, err := client.User.GetMe()
if err != nil {
    // 检查是否为 API 错误
    if apiErr, ok := kook.IsAPIError(err); ok {
        fmt.Printf("API 错误 %d: %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("网络错误: %v\n", err)
    }
    return
}
```

## 项目结构

```
kook-go-sdk/
├── kook/                 # 核心 SDK 代码
│   ├── client.go         # HTTP 客户端实现
│   ├── types.go          # 数据类型定义
│   ├── user.go           # 用户服务
│   ├── message.go        # 消息服务
│   ├── guild.go          # 服务器服务
│   ├── channel.go        # 频道服务
│   ├── websocket.go      # WebSocket 客户端
│   ├── webhook.go        # Webhook 处理器
│   └── ...               # 其他服务实现
├── examples/             # 使用示例
│   ├── simple_bot/       # 基础机器人示例
│   ├── advanced_bot/     # 高级机器人（WebSocket）
│   ├── api_usage/        # API 使用示例
│   └── webhook_bot/      # Webhook 机器人示例
├── docs/                 # 文档
├── go.mod                # Go 模块文件
└── README.md            # 本文件
```

## 环境变量

设置以下环境变量进行测试：

```bash
export KOOK_TOKEN="你的机器人令牌"
```

## 测试

运行示例程序测试 SDK：

```bash
# 测试基础功能
go run examples/simple_bot/main.go

# 测试 API 操作
go run examples/api_usage/main.go

# 测试 WebSocket 连接
go run examples/advanced_bot/main.go
```

## 贡献

1. Fork 此仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 系统要求

- Go 1.21 或更高版本
- 有效的 KOOK 机器人令牌

## 依赖项

- `github.com/gorilla/websocket` - WebSocket 客户端
- `github.com/sirupsen/logrus` - 结构化日志记录

## 许可证

本项目基于 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 免责声明

这是 KOOK API 的非官方 SDK。与 KOOK 或其开发者没有关联、认可或赞助关系。 
屎山代码，碰到问题请直接开喷，本人对项目状态敏感，打个喷嚏就知道哪里坏了，会第一时间修复并更新。
