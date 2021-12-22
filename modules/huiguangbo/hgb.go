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

	"gitee.com/lyhuilin/open_api/model/feedmsg"
	"github.com/Mrs4s/MiraiGo/message"
	"gopkg.in/yaml.v2"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &hgb{}

var hgbConf = HGBConf{}

var robot *bot.Bot

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
	path := getPathConf()

	bytes, err := util.ReadFile(path)
	if err != nil {
		log.Errorf(err, "unable to read config file in %s", path)
	}
	err = yaml.Unmarshal(bytes, &hgbConf)
	if err != nil {
		log.Errorf(err, "unable to read config file in %s", path)
	}
	wsServerUrl := viper.GetString("module.huiguangbo.server_url")
	channel := viper.GetString("module.huiguangbo.server_token")
	go wsClientStart(wsServerUrl, channel)
}

func (a *hgb) PostInit() {
}
func (a *hgb) Serve(b *bot.Bot) {
	fmt.Println("huiguangbo Serve")
}

func (a *hgb) Start(bot *bot.Bot) {
	robot = bot
	go InitHGBConf()
	fmt.Println("huiguangbo Start")
}

func (a *hgb) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("huiguangbo Stop")
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
				log.Errorf(err, "GetUrlToByte(%s)", richMsg.Image.PicURL)
			} else {
				gm, err := robot.UploadGroupImage(groupCode, bytes.NewReader(imgBin))
				if err != nil {
					log.Errorf(err, "UploadGroupImage(%d)", groupCode)
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
		err = fmt.Errorf("no msg")
	}
	return nil, err
}
func sendMsg(richMsg feedmsg.FeedRichMsgModel) {

	if !robot.Online.Load() {
		log.Debugf("机器人(%d:%s)离线，请重新登录(重新打开程序)", robot.Uin, robot.Nickname)
	}
	isConverMsg := false
	for _, v := range hgbConf.GroupList {
		groupCode := v.Id
		msg, err := richMsgToSendingMessage(groupCode, richMsg)
		if err != nil {
			continue
		} else {
			isConverMsg = true
		}

		if isConverMsg {
			for _, vv := range hgbConf.GroupList {
				if !vv.IsFeed {
					continue
				}
				robot.SendGroupMessage(vv.Id, msg)
				time.Sleep(hgbConf.SenderSleep)
			}
			break
		}

	}

}

func InitHGBConf() {
	for {
		if robot.Online.Load() {
			break
		}
		time.Sleep(10 * time.Second)
	}

	// hgbConf
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
	if hgbConf.SenderSleep <= 0 {
		hgbConf.SenderSleep = 8 * time.Second
	}
	if len(hgbConf.GroupList) > 0 {
		// err = yaml.Unmarshal(bytes, &hgbConf)
		outBody, err := yaml.Marshal(hgbConf)
		if err != nil {
			log.Errorf(err, "yaml.Marshal")
			// return
		} else {
			path := getPathConf()
			err := util.WriteFile(path, outBody)
			if err != nil {
				log.Errorf(err, "WriteFile(%s)", path)
			}
		}
	}
}
