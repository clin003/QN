/*
 * @Author: baicai_way
 * @Date: 2022-01-01
 */
package logging

import (
	"sync"

	"gitee.com/lyhuilin/QN/bot"

	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/log/lager"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func init() {
	instance = &logging{}
	bot.RegisterModule(instance)
}

type logging struct {
}

func (m *logging) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "internal.logging",
		Instance: instance,
	}
}

func (a *logging) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
	log.Infof("【QN】模块初始化=%+v", a.MiraiGoModule().ID)
}

func (m *logging) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *logging) Serve(b *bot.Bot) {
	// 注册服务函数部分
	registerLog(b)
}

func (m *logging) Start(b *bot.Bot) {
	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *logging) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}

var instance *logging

// var logger = utils.GetModuleLogger("internal.logging")

func logGroupMessage(msg *message.GroupMessage) {
	// logger.
	// 	WithField("from", "GroupMessage").
	// 	WithField("MessageID", msg.Id).
	// 	WithField("MessageIID", msg.InternalId).
	// 	WithField("GroupCode", msg.GroupCode).
	// 	WithField("SenderID", msg.Sender.Uin).
	// 	Info(msg.ToString())
	// log.Info(msg.ToString(), lager.Data{
	// 	"from":       "GroupMessage",
	// 	"MessageID":  msg.Id,
	// 	"MessageIID": msg.InternalId,
	// 	"GroupCode":  msg.GroupCode,
	// 	"SenderID":   msg.Sender.Uin,
	// })
	log.Info("群消息", lager.Data{
		"from":       "GroupMessage",
		"MessageID":  msg.Id,
		"MessageIID": msg.InternalId,
		"GroupCode":  msg.GroupCode,
		"SenderID":   msg.Sender.Uin,
		"msg":        msg.ToString(),
	})
}

func logPrivateMessage(msg *message.PrivateMessage) {
	// logger.
	// 	WithField("from", "PrivateMessage").
	// 	WithField("MessageID", msg.Id).
	// 	WithField("MessageIID", msg.InternalId).
	// 	WithField("SenderID", msg.Sender.Uin).
	// 	WithField("Target", msg.Target).
	// 	Info(msg.ToString())
	// log.Info(msg.ToString(), lager.Data{
	// 	"from":       "PrivateMessage",
	// 	"MessageID":  msg.Id,
	// 	"MessageIID": msg.InternalId,
	// 	"SenderID":   msg.Sender.Uin,
	// 	"Target":     msg.Target,
	// })
	log.Info("私聊消息", lager.Data{
		"from":       "PrivateMessage",
		"MessageID":  msg.Id,
		"MessageIID": msg.InternalId,
		"SenderID":   msg.Sender.Uin,
		"Target":     msg.Target,
		"msg":        msg.ToString(),
	})
}

func logFriendMessageRecallEvent(event *client.FriendMessageRecalledEvent) {
	// logger.
	// 	WithField("from", "FriendsMessageRecall").
	// 	WithField("MessageID", event.MessageId).
	// 	WithField("SenderID", event.FriendUin).
	// 	Info("friend message recall")
	// log.Info("friend message recall", lager.Data{
	// 	"from":      "FriendsMessageRecall",
	// 	"MessageID": event.MessageId,
	// 	"SenderID":  event.FriendUin,
	// })
	log.Info("好友消息回音", lager.Data{
		"from":      "FriendsMessageRecall",
		"MessageID": event.MessageId,
		"SenderID":  event.FriendUin,
	})
}

func logGroupMessageRecallEvent(event *client.GroupMessageRecalledEvent) {
	// logger.
	// 	WithField("from", "GroupMessageRecall").
	// 	WithField("MessageID", event.MessageId).
	// 	WithField("GroupCode", event.GroupCode).
	// 	WithField("SenderID", event.AuthorUin).
	// 	WithField("OperatorID", event.OperatorUin).
	// 	Info("group message recall")
	// log.Info("group message recall", lager.Data{
	// 	"from":       "GroupMessageRecall",
	// 	"MessageID":  event.MessageId,
	// 	"GroupCode":  event.GroupCode,
	// 	"SenderID":   event.AuthorUin,
	// 	"OperatorID": event.OperatorUin,
	// })
	log.Info("群消息回音", lager.Data{
		"from":       "GroupMessageRecall",
		"MessageID":  event.MessageId,
		"GroupCode":  event.GroupCode,
		"SenderID":   event.AuthorUin,
		"OperatorID": event.OperatorUin,
	})
}

func logGroupMuteEvent(event *client.GroupMuteEvent) {
	// logger.
	// 	WithField("from", "GroupMute").
	// 	WithField("GroupCode", event.GroupCode).
	// 	WithField("OperatorID", event.OperatorUin).
	// 	WithField("TargetID", event.TargetUin).
	// 	WithField("MuteTime", event.Time).
	// 	Info("group mute")
	// log.Info("group mute", lager.Data{
	// 	"from":       "GroupMute",
	// 	"GroupCode":  event.GroupCode,
	// 	"OperatorID": event.OperatorUin,
	// 	"TargetID":   event.TargetUin,
	// 	"MuteTime":   event.Time,
	// })
	log.Info("群事件", lager.Data{
		"from":       "GroupMute",
		"GroupCode":  event.GroupCode,
		"OperatorID": event.OperatorUin,
		"TargetID":   event.TargetUin,
		"MuteTime":   event.Time,
	})
}

func logDisconnect(event *client.ClientDisconnectedEvent) {
	// logger.
	// 	WithField("from", "Disconnected").
	// 	WithField("reason", event.Message).
	// 	Warn("bot disconnected")
	// log.Warn("bot disconnected", lager.Data{
	// 	"from":   "Disconnected",
	// 	"reason": event.Message,
	// })
	log.Warn("断开连接", lager.Data{
		"from":   "Disconnected",
		"reason": event.Message,
	})
}

func registerLog(b *bot.Bot) {
	b.GroupMessageRecalledEvent.Subscribe(func(qqClient *client.QQClient, event *client.GroupMessageRecalledEvent) {
		logGroupMessageRecallEvent(event)
	})

	b.GroupMessageEvent.Subscribe(func(qqClient *client.QQClient, groupMessage *message.GroupMessage) {
		logGroupMessage(groupMessage)
	})

	b.GroupMuteEvent.Subscribe(func(qqClient *client.QQClient, event *client.GroupMuteEvent) {
		logGroupMuteEvent(event)
	})

	b.PrivateMessageEvent.Subscribe(func(qqClient *client.QQClient, privateMessage *message.PrivateMessage) {
		logPrivateMessage(privateMessage)
	})

	b.FriendMessageRecalledEvent.Subscribe(func(qqClient *client.QQClient, event *client.FriendMessageRecalledEvent) {
		logFriendMessageRecallEvent(event)
	})

	b.DisconnectedEvent.Subscribe(func(qqClient *client.QQClient, event *client.ClientDisconnectedEvent) {
		logDisconnect(event)
	})

}
