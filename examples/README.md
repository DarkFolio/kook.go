# KOOK Go SDK 示例程序

本目录包含了KOOK Go SDK的完整示例程序，展示了如何使用SDK的各种功能。

## 环境准备

在运行任何示例之前，请确保：

1. 已安装Go 1.21或更高版本
2. 已获取KOOK机器人Token
3. 设置环境变量：

```bash
# Windows PowerShell
$env:KOOK_TOKEN="你的机器人令牌"

# Linux/macOS
export KOOK_TOKEN="你的机器人令牌"
```

## 示例程序列表

### 1. simple_bot - 基础机器人

**文件位置**: `examples/simple_bot/main.go`

**功能介绍**:
- 展示最基本的机器人功能
- WebSocket连接和消息监听
- 简单的ping-pong回复机制
- 适合新手入门学习

**主要特性**:
- 获取机器人基本信息
- 建立WebSocket连接
- 监听文本消息事件
- 自动回复"ping"消息为"pong"
- 优雅关闭处理

**运行方法**:
```bash
cd examples/simple_bot
go run main.go
```

**使用说明**:
1. 机器人启动后会显示基本信息
2. 在KOOK频道中发送"ping"，机器人会回复"pong"
3. 按Ctrl+C优雅关闭机器人

### 2. advanced_bot - 高级机器人

**文件位置**: `examples/advanced_bot/main.go`

**功能介绍**:
- 完整的命令系统机器人
- 展示多种API的综合使用
- 包含丰富的交互功能
- 适合实际项目参考

**支持的命令**:
- `!help` - 显示帮助信息
- `!roles` - 查看服务器角色列表
- `!emojis` - 查看服务器表情列表
- `!blacklist` - 查看服务器屏蔽用户
- `!pin` - 置顶当前消息
- `!game <游戏名>` - 设置游戏动态
- `!music <歌手> <歌名>` - 设置音乐动态
- `!regions` - 查看可用区域
- `!invites` - 查看服务器邀请
- `!badges` - 查看服务器徽章
- `!nickname <昵称>` - 修改你的昵称
- `!upload` - 上传文件示例

**运行方法**:
```bash
cd examples/advanced_bot
go run main.go
```

**使用说明**:
1. 在KOOK频道中使用以上命令
2. 部分命令需要相应权限才能执行
3. 机器人会根据命令执行相应操作并返回结果

### 3. api_usage - API使用演示

**文件位置**: `examples/api_usage/main.go`

**功能介绍**:
- 演示各种API的基本使用方法
- 展示用户、服务器、消息API的调用
- 适合学习API调用方式
- 一次性执行，不保持连接

**演示内容**:
- 用户API：获取当前用户信息
- 服务器API：获取服务器列表和成员信息
- 消息API：发送消息和获取消息列表（需要配置频道ID）

**运行方法**:
```bash
cd examples/api_usage
go run main.go
```

**注意事项**:
- 需要将代码中的`YOUR_CHANNEL_ID_HERE`替换为实际的频道ID才能演示消息API
- 程序执行完毕后会自动退出

### 4. webhook_bot - Webhook机器人

**文件位置**: `examples/webhook_bot/main.go`

**功能介绍**:
- 演示Webhook模式的机器人
- HTTP服务器接收KOOK事件
- 适合服务器部署使用
- 相比WebSocket更稳定

**主要特性**:
- 启动HTTP服务器监听Webhook
- 处理KOOK发送的事件
- 简单的hello回复功能
- 支持验证Token验证

**环境变量**:
```bash
# 必需
export KOOK_TOKEN="你的机器人令牌"

# 可选（用于验证请求来源）
export KOOK_VERIFY_TOKEN="你的验证令牌"
```

**运行方法**:
```bash
cd examples/webhook_bot
go run main.go
```

**配置说明**:
1. 程序会在8080端口启动HTTP服务器
2. Webhook接收地址：`http://你的服务器:8080/webhook`
3. 需要在KOOK开发者后台配置Webhook URL

### 5. complete_api_demo - 完整API演示

**文件位置**: `examples/complete_api_demo/main.go`

**功能介绍**:
- 演示SDK支持的所有主要API
- 完整的功能覆盖测试
- 适合测试SDK完整性
- 展示API的实际返回数据

**演示的API模块**:
- 用户信息和在线状态
- 服务器列表和详细信息
- 角色管理API
- 频道管理API
- 邀请管理API
- 游戏和动态API
- 好友管理API
- 消息功能API

**运行方法**:
```bash
cd examples/complete_api_demo
go run main.go
```

**输出内容**:
- 机器人基本信息
- 所有服务器的详细信息
- 各种API的调用结果
- 错误处理演示

## 通用说明

### 错误处理

所有示例都包含完善的错误处理：
```go
if err != nil {
    log.Printf("操作失败: %v", err)
    // 根据需要处理错误
}
```

### 日志输出

示例使用Go标准库的`log`包输出日志：
- 正常信息使用`log.Printf`
- 严重错误使用`log.Fatal`
- 调试信息使用`log.Println`

### 权限要求

部分功能需要机器人具有相应权限：
- 消息管理：发送消息、置顶消息
- 服务器管理：查看成员、管理角色
- 频道管理：查看频道列表、管理频道

### 开发建议

1. **从simple_bot开始**：如果是初学者，建议从simple_bot开始学习
2. **参考advanced_bot**：开发实际项目时可以参考advanced_bot的架构
3. **使用complete_api_demo测试**：用来验证Token是否有效和权限是否充足
4. **选择合适的模式**：WebSocket适合实时性要求高的场景，Webhook适合服务器部署

### 常见问题

**Q: 机器人无法接收消息？**
A: 检查机器人是否已加入服务器，是否有查看频道的权限

**Q: API调用返回403错误？**
A: 检查机器人Token是否正确，是否有执行该操作的权限

**Q: WebSocket连接失败？**
A: 检查网络连接，确认Token有效性

**Q: Webhook收不到事件？**
A: 确认Webhook URL配置正确，服务器端口已开放

### 安全注意事项

1. **保护Token安全**：不要在代码中硬编码Token
2. **使用环境变量**：通过环境变量传递敏感信息
3. **验证Webhook**：使用verify_token验证Webhook请求来源
4. **错误处理**：避免在日志中泄露敏感信息

## 扩展开发

基于这些示例，你可以：

1. **添加更多命令**：在advanced_bot基础上扩展功能
2. **集成数据库**：存储用户数据和配置信息
3. **添加定时任务**：定期执行某些操作
4. **接入其他服务**：与外部API集成
5. **部署到云服务**：使用Docker容器化部署

## 贡献

如果你有新的示例想法或发现了bug，欢迎：
1. 提交Issue报告问题
2. 提交Pull Request贡献代码
3. 改进现有示例的文档

## 相关链接

- [KOOK官方文档](https://developer.kookapp.cn/doc/reference)
- [Go语言官网](https://golang.org/)
- [项目主页](../README.md) 