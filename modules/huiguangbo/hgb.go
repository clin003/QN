/*
 * @Author: baicai_way
 * @Date: 2022-01-01
 */
package huiguangbo

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	"gitee.com/lyhuilin/QN/bot"
	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/util"

	// "gitee.com/lyhuilin/open_api/model/feedmsg"
	"gitee.com/lyhuilin/model/feedmsg"

	"github.com/Mrs4s/MiraiGo/message"
	"gopkg.in/yaml.v2"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &hgb{}

var hgbConf = HGBConf{}

var robot *bot.Bot
var guildIDChannelID string
var isReady bool
var once sync.Once

type hgb struct {
}

func (a *hgb) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "module.huiguangbo",
		Instance: instance,
	}
}
func getPathConf() (retText string) {
	path := viper.GetString("module.huiguangbo.path")

	if path == "" {
		path = "./conf/huiguangbo/huiguangbo.yaml"
	}
	return path
}
func (a *hgb) Init() {
	log.Infof("【QN】模块初始化=%+v", a.MiraiGoModule().ID)
	path := getPathConf()

	bytes, err := util.ReadFile(path)
	if err != nil {
		log.Errorf(err, "读取配置文件 %s 出错了", path)
	}
	err = yaml.Unmarshal(bytes, &hgbConf)
	if err != nil {
		log.Errorf(err, "加载配置文件 %s 出错了", path)
	}
	wsServerUrl := viper.GetString("module.huiguangbo.server_url")
	channel := viper.GetString("module.huiguangbo.server_token")

	// go wsClientStart(wsServerUrl, channel)
	go initWsServer(wsServerUrl, channel)
}

func (a *hgb) PostInit() {
}
func (a *hgb) Serve(b *bot.Bot) {
	fmt.Println("慧广播 Serve")
}

func (a *hgb) Start(bot *bot.Bot) {
	robot = bot
	go InitHGBConf()
	fmt.Println("慧广播 Start")
}

func (a *hgb) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("慧广播 Stop")
}

// 将richMsg消息转化为SendingMessage
func richMsgToSendingMessage(groupCode int64, richMsg feedmsg.FeedRichMsgModel) (retMsg *message.SendingMessage, err error) {
	m := message.NewSendingMessage()
	if richMsg.Msgtype == "rich" {
		if len(richMsg.Text.Content) > 0 {
			m.Append(message.NewText(richMsg.Text.Content))
		}

		if len(richMsg.Image.PicURL) > 0 && strings.HasPrefix(richMsg.Image.PicURL, "http") {
			imgBin, err := util.GetUrlToByte(richMsg.Image.PicURL)
			if err != nil {
				log.Errorf(err, "下载图片文件(%s)出错(GetUrlToByte)", richMsg.Image.PicURL)
			} else {
				gm, err := robot.UploadGroupImage(groupCode, bytes.NewReader(imgBin))
				if err != nil {
					log.Errorf(err, "上传图片文件(%d)出错啦(UploadGroupImage)", groupCode)
				} else {
					m.Append(gm)
				}
			}
		}
	}

	if len(m.Elements) > 0 {
		retMsg = m
		return retMsg, nil
	}
	if err == nil {
		err = fmt.Errorf("no msg(空消息):%v", richMsg)
	}
	return nil, err
}
func sendMsg(richMsg feedmsg.FeedRichMsgModel) {
	if !isReady {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("run time panic(sendMsg): %v", err)
			return
		}
	}()
	log.Infof("收到广播消息，开始处理(%s)", richMsg.ToString())
	go sendMsgGuild(richMsg)

	if !robot.Online.Load() {
		log.Warnf("机器人(%d:%s)离线，请重新登录(重新打开程序)", robot.Uin, robot.Nickname)
	}
	// isConverMsg := false
	// 处理(格式化)待发布消息
	for _, v := range hgbConf.GroupList {
		if !v.IsFeed {
			continue
		}

		if hgbConf.SenderSleep <= 100*time.Microsecond {
			vID := v.Id
			go func(groupID int64) {
				groupCode := groupID
				msg, err := richMsgToSendingMessage(groupCode, richMsg)
				if err != nil {
					log.Errorf(err, "消息处理失败(%d): %s", groupCode, richMsg.ToString())
					return
				}

				// 广播消息
				sendResult := robot.SendGroupMessage(groupCode, msg)
				if sendResult != nil {
					log.Infof("群(%d) 广播模式 已启用,发送消息 (ID: %d InternalId: %d ) ", groupCode, sendResult.Id, sendResult.InternalId) //, sendResult.ToString()
				} else {
					log.Infof("群(%d) 广播模式 已启用,发送消息 失败 :%s", groupCode, richMsg.ToString())
				}
			}(vID)
		} else {
			groupCode := v.Id
			msg, err := richMsgToSendingMessage(groupCode, richMsg)
			if err != nil {
				log.Errorf(err, "消息处理失败(%d): %s", groupCode, richMsg.ToString())
				continue
			}

			// 广播消息
			sendResult := robot.SendGroupMessage(groupCode, msg)
			if sendResult != nil {
				log.Infof("群(%d) 广播模式 已启用,发送消息 (ID: %d InternalId: %d ) ", groupCode, sendResult.Id, sendResult.InternalId) //, sendResult.ToString()
			} else {
				log.Infof("群(%d) 广播模式 已启用,发送消息 失败 :%s", groupCode, richMsg.ToString())
			}

			time.Sleep(hgbConf.SenderSleep)
		}

		// groupCode := v.Id
		// msg, err := richMsgToSendingMessage(groupCode, richMsg)
		// if err != nil {
		// 	log.Errorf(err, "消息处理失败(%d): %s", groupCode, richMsg.ToString())
		// 	continue
		// }

		// // 广播消息
		// sendResult := robot.SendGroupMessage(v.Id, msg)
		// if sendResult != nil {
		// 	log.Infof("群(%d) 广播模式 已启用,发送消息 (ID: %d InternalId: %d ) ", v.Id, sendResult.Id, sendResult.InternalId) //, sendResult.ToString()
		// } else {
		// 	log.Infof("群(%d) 广播模式 已启用,发送消息 失败 :%s", v.Id, richMsg.ToString())
		// }

		// time.Sleep(hgbConf.SenderSleep)
	}
}

// 将richMsg消息转化为GuildSendingMessage
func richMsgToGuildSendingMessage(guildID, channelId uint64, richMsg feedmsg.FeedRichMsgModel) (retMsg *message.SendingMessage, err error) {
	m := message.NewSendingMessage()
	if richMsg.Msgtype == "rich" {
		if len(richMsg.Text.Content) > 0 {
			m.Append(message.NewText(richMsg.Text.Content))
		}

		if len(richMsg.Image.PicURL) > 0 && strings.HasPrefix(richMsg.Image.PicURL, "http") {
			imgBin, err := util.GetUrlToByte(richMsg.Image.PicURL)
			if err != nil {
				log.Errorf(err, "下载图片文件(%s)出错(GetUrlToByte)", richMsg.Image.PicURL)
			} else {
				// gm, err := robot.UploadGroupImage(groupCode, bytes.NewReader(imgBin))
				gm, err := robot.GuildService.UploadGuildImage(guildID, channelId, bytes.NewReader(imgBin))
				if err != nil {
					log.Errorf(err, "上传图片文件(%d:%d)出错啦(UploadGroupImage)", guildID, channelId)
				} else {
					m.Append(gm)
				}
			}
		}
	}

	if len(m.Elements) > 0 {
		retMsg = m
		return retMsg, nil
	}
	if err == nil {
		err = fmt.Errorf("no msg(空消息):%v", richMsg)
	}
	return nil, err
}
func sendMsgGuild(richMsg feedmsg.FeedRichMsgModel) {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("run time panic(sendMsgGuild): %v", err)
			return
		}
	}()
	if !hgbConf.GuildSenderEnable {
		return
	}
	log.Infof("收到广播消息，开始处理(频道)(%s)", richMsg.ToString())
	if !robot.Online.Load() {
		log.Warnf("机器人(%d:%s)离线，请重新登录(重新打开程序)", robot.Uin, robot.Nickname)
	}
	_guildIDChannelID := guildIDChannelID
	// guildID_channelID := fmt.Sprintf("%d_%d", c.Id, cc.Id)
	sendGuildFun := func(guildId, channelId uint64) {
		msg, err := richMsgToGuildSendingMessage(guildId, channelId, richMsg)
		if err != nil {
			log.Errorf(err, "消息处理失败(%d:%d): %s", guildId, channelId, richMsg.ToString())
			return
		}

		// 广播消息
		if sendResult, err := robot.GuildService.SendGuildChannelMessage(guildId, channelId, msg); err != nil {
			log.Errorf(err, "频道(%d:%d) 广播模式 已启用,发送消息 失败 :%s", guildId, channelId, richMsg.ToString())
		} else if sendResult != nil {
			log.Infof("频道(%d:%d) 广播模式 已启用,发送消息 (ID: %d InternalId: %d ) ", guildId, channelId, sendResult.Id, sendResult.InternalId) //, sendResult.ToString()
		} else {
			log.Errorf(err, "频道(%d:%d) 广播模式 已启用,发送消息 失败 :%s", guildId, channelId, richMsg.ToString())
		}
	}

	for _, c := range robot.GuildService.Guilds {
		for _, cc := range c.Channels {
			guildId := c.GuildId
			channelId := cc.ChannelId
			guildID_channelID := fmt.Sprintf("%d_%d", guildId, channelId)
			if strings.Contains(_guildIDChannelID, guildID_channelID) {

				if hgbConf.SenderSleep <= 100*time.Microsecond {
					go sendGuildFun(guildId, channelId)
				} else {
					sendGuildFun(guildId, channelId)
					time.Sleep(hgbConf.SenderSleep)
				}

			}
		}
	}
}

func InitHGBConf() {
	for {
		if robot.Online.Load() {
			break
		}
		time.Sleep(20 * time.Second)
	}
	log.Infof("开始 加载慧广播配置信息")
	for _, v := range robot.GroupList {
		groupName := v.Name
		groupCode := v.Code
		if len(groupName) > 0 && groupCode > 0 {
			var groupInfo GroupInfo
			groupInfo.Id = groupCode
			groupInfo.Name = groupName
			isInConf := false
			for _, c := range hgbConf.GroupList {
				if c.Id == groupCode {
					isInConf = true
					break
				}
			}
			if !isInConf {
				hgbConf.GroupList = append(hgbConf.GroupList, groupInfo)
			}
		}
	}
	// 频道信息
	guildIDChannelID = strings.Join(hgbConf.FeedGuildList, ",")
	for _, v := range robot.GuildService.Guilds {
		// log.Infof("%s(ID:%d Code:%d)", v.GuildName, v.GuildId, v.GuildCode)
		var guildInfo GuildInfo
		guildInfo.Name = v.GuildName
		guildInfo.Id = v.GuildId

		for _, vv := range v.Channels {
			// log.Infof("%s(ID:%d  EventTime:%d)", vv.ChannelName, vv.ChannelId, vv.EventTime)
			var channelInfo ChannelInfo
			channelInfo.Id = vv.ChannelId
			channelInfo.Name = vv.ChannelName
			channelInfo.EventTime = vv.EventTime

			isInConf := false
			for _, c := range hgbConf.GuildList {
				if c.Id == guildInfo.Id {
					for _, cc := range c.ChannelList {
						if cc.Id == channelInfo.Id {
							if cc.IsFeed {
								// 更新频道子频道发单列表
								guildID_channelID := fmt.Sprintf("%d_%d", c.Id, cc.Id)
								if !strings.Contains(guildIDChannelID, guildID_channelID) {
									hgbConf.FeedGuildList = append(hgbConf.FeedGuildList, guildID_channelID)
									guildIDChannelID = strings.Join(hgbConf.FeedGuildList, ",")
									// guildIDChannelID = fmt.Sprintf("%s,%s", guildIDChannelID, guildID_channelID)
								}
							}
							isInConf = true
							break //c.ChannelList
						}
					}
					break
				}
			} //hgbConf.GuildList
			// 子频道信息不存在 添加频道子频道列表
			if !isInConf {
				guildInfo.ChannelList = append(guildInfo.ChannelList, channelInfo)
			}
		} //v.Channels

		isInConf := false
		for _, c := range hgbConf.GuildList {
			if c.Id == guildInfo.Id {
				isInConf = true
				break
			}
		}
		if !isInConf {
			hgbConf.GuildList = append(hgbConf.GuildList, guildInfo)
		}
	}

	if hgbConf.SenderSleep <= 0 {
		hgbConf.SenderSleep = 100 * time.Microsecond
	}

	if len(hgbConf.GroupList) > 0 {
		outBody, err := yaml.Marshal(hgbConf)
		if err != nil {
			log.Errorf(err, "生成配置信息编码出错(yaml.Marshal):%v", hgbConf)
		} else {
			path := getPathConf()
			err := util.WriteFile(path, outBody)
			if err != nil {
				log.Errorf(err, "写入配置信息到文件(%s)出错(WriteFile)", path)
			}
		}
	}
	log.Infof("完成 加载慧广播配置信息")
	once.Do(func() {
		isReady = true
	})

}
