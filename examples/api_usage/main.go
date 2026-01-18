package main

import (
	"context"
	"log"
	"os"

	"kook-go-sdk/kook"
)

func main() {
	// 从环境变量获取Token
	token := os.Getenv("KOOK_TOKEN")
	if token == "" {
		log.Fatal("请设置环境变量 KOOK_TOKEN")
	}

	// 创建客户端
	client := kook.NewClient(token)

	// 演示用户API
	demonstrateUserAPI(client)

	// 演示服务器API
	demonstrateGuildAPI(client)

	// 演示消息API
	demonstrateMessageAPI(client)
}

func demonstrateUserAPI(client *kook.Client) {
	log.Println("=== 用户API演示 ===")

	// 获取当前用户信息
	user, err := client.User.GetMe(context.Background())
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		return
	}

	log.Printf("当前用户: %s#%s", user.Username, user.IdentifyNum)
	log.Printf("用户ID: %s", user.ID)
	log.Printf("是否为机器人: %v", user.Bot)
}

func demonstrateGuildAPI(client *kook.Client) {
	log.Println("=== 服务器API演示 ===")

	// 获取服务器列表
	guilds, err := client.Guild.GetGuildList(context.Background(), 1, 10, "")
	if err != nil {
		log.Printf("获取服务器列表失败: %v", err)
		return
	}

	log.Printf("服务器总数: %d", guilds.Meta.Total)
	for _, guild := range guilds.Items {
		log.Printf("服务器: %s (ID: %s)", guild.Name, guild.ID)

		// 获取服务器成员列表
		members, err := client.Guild.GetGuildMembers(context.Background(), guild.ID, 1, 5, "")
		if err != nil {
			log.Printf("获取服务器成员失败: %v", err)
			continue
		}

		log.Printf("  成员数量: %d", members.Meta.Total)
		for _, member := range members.Items {
			log.Printf("  成员: %s", member.Username)
		}
	}
}

func demonstrateMessageAPI(client *kook.Client) {
	log.Println("=== 消息API演示 ===")

	// 这里需要一个实际的频道ID来演示
	// 在实际使用中，你应该从事件或其他API获取频道ID
	channelID := "YOUR_CHANNEL_ID_HERE"

	if channelID == "YOUR_CHANNEL_ID_HERE" {
		log.Println("跳过消息API演示 - 需要设置实际的频道ID")
		return
	}

	// 发送消息
	params := kook.SendMessageParams{
		TargetID: channelID,
		Content:  "Hello from KOOK Go SDK! 🚀",
		MsgType:  1,
	}

	message, err := client.Message.SendMessage(context.Background(), params)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
		return
	}

	log.Printf("消息发送成功: %s", message.ID)

	// 获取消息列表
	listParams := kook.GetMessageListParams{
		PageSize: 10,
	}

	messages, err := client.Message.GetMessageList(context.Background(), channelID, listParams)
	if err != nil {
		log.Printf("获取消息列表失败: %v", err)
		return
	}

	log.Printf("获取到 %d 条消息", len(messages.Items))
	for _, msg := range messages.Items {
		log.Printf("消息: %s (作者: %s)", msg.Content, msg.Author.Username)
	}
}
 