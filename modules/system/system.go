/*
 * @Author: baicai_way
 * @Date: 2022-01-01
 */
package system

import (
	"fmt"

	"sync"

	"gitee.com/lyhuilin/QN/bot"
	"gitee.com/lyhuilin/log"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &system{}

type system struct {
}

func (a *system) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "module.system",
		Instance: instance,
	}
}

func (a *system) Init() {
	log.Infof("【QN】模块初始化=%+v", a.MiraiGoModule().ID)
}

func (a *system) PostInit() {
}
func (a *system) Serve(b *bot.Bot) {

	// b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
	// 	fmt.Printf("message=%+v\n", msg.Sender.Nickname)
	// 	if botObj := bot.Instance; botObj == nil {
	// 		// 机器人已下线，直接结束回复流程
	// 		fmt.Println("【收到消息】机器人已下线，直接结束回复流程")
	// 		return
	// 	}

	// 	groupInfo := c.FindGroup(msg.GroupCode)
	// 	if groupInfo == nil {
	// 		// QQ 群信息获取失败，结束流程
	// 		fmt.Println("【收到消息】QQ 群信息获取失败，直接结束回复流程")
	// 		return
	// 	}
	// 	groupMemberInfo := groupInfo.FindMember(c.Uin)
	// 	botName := ""
	// 	if groupMemberInfo == nil {
	// 		// QQ 群我的数据获取失败，直接赋值昵称
	// 		botName = b.Nickname
	// 	} else {
	// 		botName = groupMemberInfo.DisplayName()
	// 	}

	// 	fmt.Printf("群昵称=%+v\n", botName)

	// 	fmt.Println("【收到消息】" + msg.ToString())

	// })

	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
		// fmt.Printf("message=%+v\n", msg.ToString())
		if botObj := bot.Instance; botObj == nil {
			// 机器人已下线，直接结束回复流程
			fmt.Println("【收到消息】机器人已下线，直接结束回复流程")
			return
		}

		if msg.ToString() == "\\help" {
			m := message.NewSendingMessage().Append(message.NewText(`【指令菜单】
\help (帮助菜单)
\info (系统状态)
`))
			c.SendPrivateMessage(msg.Sender.Uin, m)
			return
		}

		if msg.ToString() == "\\info" {
			m := message.NewSendingMessage().Append(message.NewText("系统状态: 正常\n机器人状态: 在线"))
			c.SendPrivateMessage(msg.Sender.Uin, m)
			return
		}

	})
}

func (a *system) Start(bot *bot.Bot) {
}

func (a *system) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
