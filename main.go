package main

import (
	"log"
	"os"

	"kook-go-sdk/kook"
)

func main() {
	// 从环境变量获取Token
	token := os.Getenv("KOOK_BOT_TOKEN")
	if token == "" {
		log.Println("请设置环境变量 KOOK_BOT_TOKEN")
		log.Println("示例: export KOOK_BOT_TOKEN=Bot_your_token_here")
		return
	}

	// 创建客户端
	client := kook.NewClient(token)

	// 获取机器人信息
	user, err := client.User.GetMe()
	if err != nil {
		log.Fatalf("获取机器人信息失败: %v", err)
	}

	log.Printf("KOOK Go SDK 测试成功!")
	log.Printf("机器人名称: %s#%s", user.Username, user.IdentifyNum)
	log.Printf("机器人ID: %s", user.ID)
	log.Printf("是否在线: %v", user.Online)
	log.Printf("是否为机器人: %v", user.Bot)

	// 获取服务器列表
	guilds, err := client.Guild.GetGuildList(1, 5, "")
	if err != nil {
		log.Printf("获取服务器列表失败: %v", err)
	} else {
		log.Printf("服务器列表 (前5个):")
		for i, guild := range guilds.Items {
			log.Printf("  %d. %s (ID: %s)", i+1, guild.Name, guild.ID)
		}
		log.Printf("总服务器数: %d", guilds.Meta.Total)
	}

	log.Println("\nSDK 功能测试完成！")
	log.Println("查看 examples/ 目录了解更多用法")
	log.Println("查看 docs/api.md 了解完整API文档")
}