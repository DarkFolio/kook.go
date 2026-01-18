package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kook-go-sdk/kook"
)

func main() {
	// 获取环境变量
	token := os.Getenv("KOOK_TOKEN")
	if token == "" {
		log.Fatal("请设置环境变量 KOOK_TOKEN")
	}

	// 创建客户端
	client := kook.NewClient(token)

	// 获取当前用户信息
	user, err := client.User.GetMe(context.Background())
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
		return
	}

	fmt.Printf("机器人名称: %s#%s\n", user.Username, user.IdentifyNum)
	fmt.Printf("机器人ID: %s\n", user.ID)
	fmt.Printf("是否在线: %t\n", user.Online)

	// 创建WebSocket客户端
	ws := kook.NewWebSocketClient(client, false)

	// 注册消息事件处理器
	ws.OnEvent(kook.EventTypeTextMessage, func(event *kook.Event) {
		log.Printf("收到消息: %s", event.Content)

		// 简单的回复逻辑
		if event.Content == "ping" {
			// 发送回复消息
			params := kook.SendMessageParams{
				TargetID: event.TargetID,
				Content:  "pong",
				MsgType:  1,
			}

			_, err := client.Message.SendMessage(context.Background(), params)
			if err != nil {
				log.Printf("发送消息失败: %v", err)
			}
		}
	})

	// 连接WebSocket
	if err := ws.Connect(); err != nil {
		log.Fatalf("WebSocket连接失败: %v", err)
	}

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("正在关闭机器人...")
	ws.Close()
}
 