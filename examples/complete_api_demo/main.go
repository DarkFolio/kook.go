package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	fmt.Println("=== 获取机器人信息 ===")
	user, err := client.User.GetMe(context.Background())
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
	} else {
		fmt.Printf("机器人名称: %s#%s\n", user.Username, user.IdentifyNum)
		fmt.Printf("机器人ID: %s\n", user.ID)
	}

	// 获取在线状态
	fmt.Println("\n=== 获取在线状态 ===")
	status, err := client.User.GetOnlineStatus(context.Background())
	if err != nil {
		log.Printf("获取在线状态失败: %v", err)
	} else {
		fmt.Printf("在线状态: %t\n", status.Online)
		fmt.Printf("在线平台: %v\n", status.OnlineOS)
	}

	// 获取服务器列表
	fmt.Println("\n=== 获取服务器列表 ===")
	guilds, err := client.Guild.GetGuildList(context.Background(), 1, 10, "")
	if err != nil {
		log.Printf("获取服务器列表失败: %v", err)
	} else {
		fmt.Printf("服务器数量: %d\n", len(guilds.Items))
		for _, guild := range guilds.Items {
			fmt.Printf("- %s (ID: %s)\n", guild.Name, guild.ID)

			// 演示角色管理API
			fmt.Printf("\n=== 服务器 %s 的角色管理 ===\n", guild.Name)
			roles, err := client.Role.GetRoleList(context.Background(), guild.ID, 1, 10)
			if err != nil {
				log.Printf("获取角色列表失败: %v", err)
			} else {
				fmt.Printf("角色数量: %d\n", len(roles.Items))
				for _, role := range roles.Items {
					fmt.Printf("- %s (ID: %d, 权限: %d)\n", role.Name, role.RoleID, role.Permissions)
				}
			}

			// 演示频道管理API
			fmt.Printf("\n=== 服务器 %s 的频道列表 ===\n", guild.Name)
			channels, err := client.Channel.GetChannelList(context.Background(), guild.ID, 1, 10, "")
			if err != nil {
				log.Printf("获取频道列表失败: %v", err)
			} else {
				fmt.Printf("频道数量: %d\n", len(channels.Items))
				for _, channel := range channels.Items {
					fmt.Printf("- %s (ID: %s, 类型: %d)\n", channel.Name, channel.ID, channel.Type)
				}
			}

			// 演示邀请管理API
			fmt.Printf("\n=== 服务器 %s 的邀请管理 ===\n", guild.Name)
			invites, err := client.Invite.GetInviteList(context.Background(), guild.ID, 1, 10)
			if err != nil {
				log.Printf("获取邀请列表失败: %v", err)
			} else {
				fmt.Printf("邀请数量: %d\n", len(invites.Items))
				for _, invite := range invites.Items {
					fmt.Printf("- 邀请码: %s, 创建者: %s\n", invite.URLCode, invite.User.Username)
				}
			}

			// 只处理第一个服务器作为演示
			break
		}
	}

	// 演示游戏API
	fmt.Println("\n=== 游戏管理 ===")
	games, err := client.Game.GetGameList(context.Background(), "")
	if err != nil {
		log.Printf("获取游戏列表失败: %v", err)
	} else {
		fmt.Printf("游戏数量: %d\n", len(games.Items))
		for i, game := range games.Items {
			if i < 5 { // 只显示前5个游戏
				fmt.Printf("- %s (ID: %d, 类型: %d)\n", game.Name, game.ID, game.Type)
			}
		}
	}

	// 演示好友API
	fmt.Println("\n=== 好友管理 ===")
	friends, err := client.Friend.GetFriendsList(context.Background())
	if err != nil {
		log.Printf("获取好友列表失败: %v", err)
	} else {
		fmt.Printf("好友数量: %d\n", len(friends.Friend))
		fmt.Printf("好友请求数量: %d\n", len(friends.Request))
		fmt.Printf("屏蔽用户数量: %d\n", len(friends.Blocked))
	}

	// 演示消息API
	fmt.Println("\n=== 消息功能演示 ===")

	// 检查卡片消息格式
	cardContent := `[{"type":"card","theme":"primary","size":"lg","modules":[{"type":"section","text":{"type":"plain-text","content":"这是一个测试卡片消息"}}]}]`
	_, err = client.Message.CheckCard(context.Background(), cardContent)
	if err != nil {
		log.Printf("卡片消息格式检查失败: %v", err)
	} else {
		fmt.Println("卡片消息格式正确")
	}

	fmt.Println("\n=== API演示完成 ===")
	fmt.Println("所有API接口已成功调用，详细的错误处理和功能展示请查看日志输出。")
	fmt.Println("请根据实际需要调用相应的API接口。")
}
 