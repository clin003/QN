package autoreply

import (
	"bytes"
	// "fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"gitee.com/lyhuilin/QN/bot"
	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/util"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"gopkg.in/yaml.v2"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &ar{}

var tem map[string]string

type ar struct {
}

func (a *ar) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "module.autoreply",
		Instance: instance,
	}
}

func (a *ar) Init() {
	path := viper.GetString("module.autoreply.path")

	if path == "" {
		path = "./conf/autoreply/autoreply.yaml"
	}

	bytes, err := util.ReadFile(path)
	if err != nil {
		log.Errorf(err, "unable to read config file in %s", path)
	}
	err = yaml.Unmarshal(bytes, &tem)
	if err != nil {
		log.Errorf(err, "unable to read config file in %s", path)
	}
}

func (a *ar) PostInit() {
}

// func (a *ar) Serve(b *bot.Bot) {
// 	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
// 		out := autoreply(msg.ToString())
// 		if out == "" {
// 			return
// 		}
// 		m := message.NewSendingMessage().Append(message.NewText(out))
// 		c.SendGroupMessage(msg.GroupCode, m)
// 	})

// 	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
// 		out := autoreply(msg.ToString())
// 		if out == "" {
// 			return
// 		}
// 		m := message.NewSendingMessage().Append(message.NewText(out))
// 		c.SendPrivateMessage(msg.Sender.Uin, m)
// 	})
// }
// func (a *ar) Serve(b *bot.Bot) {
// 	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
// 		m, err := autoreplyEx(msg.ToString())
// 		if err != nil {
// 			return
// 		}
// 		// m := message.NewSendingMessage().Append(message.NewText(out))
// 		c.SendGroupMessage(msg.GroupCode, m)
// 	})

// 	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
// 		m, err := autoreplyEx(msg.ToString())
// 		if err != nil {
// 			return
// 		}
// 		// m := message.NewSendingMessage().Append(message.NewText(out))
// 		c.SendPrivateMessage(msg.Sender.Uin, m)
// 	})
// }
func (a *ar) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		out := autoreply(msg.ToString())
		if out == "" {
			return
		}
		// m := message.NewSendingMessage().Append(message.NewText(out))
		if strings.Contains(out, "http") {
			imgBin, err := util.GetUrlToByte(out)
			if err == nil {

				// m = message.NewSendingMessage().Append(message.NewImage(imgBin))
				util.WriteFile("./img.jpg", imgBin)
				// c.UploadGroupImageByFile()
				// bytes.NewReader()
				gm, err := c.UploadGroupImage(msg.GroupCode, bytes.NewReader(imgBin))
				if err != nil {
					log.Errorf(err, "UploadGroupImage")
					return
				}
				m := message.NewSendingMessage().Append(gm)
				// m.Append(gm)
				c.SendGroupMessage(msg.GroupCode, m)
				return
			}
		}
		m := message.NewSendingMessage().Append(message.NewText(out))
		c.SendGroupMessage(msg.GroupCode, m)
	})

	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
		out := autoreply(msg.ToString())
		if out == "" {
			return
		}
		m := message.NewSendingMessage().Append(message.NewText(out))
		c.SendPrivateMessage(msg.Sender.Uin, m)
	})
}
func (a *ar) Start(bot *bot.Bot) {
}

func (a *ar) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func autoreply(in string) string {
	out, ok := tem[in]
	if !ok {
		return ""
	}
	return out
}

// func autoreplyEx(in string) (retMsg *message.SendingMessage, err error) {
// 	out, ok := tem[in]
// 	if !ok {
// 		err = fmt.Errorf("no ar key")
// 		return nil, err
// 	}
// 	if strings.Contains(out, "http") {
// 		imgBin, err := util.GetUrlToByte(out)
// 		if err != nil {
// 			retMsg = message.NewSendingMessage().Append(message.NewText(out)) // message.NewText(out)
// 			return retMsg, err
// 		}
// 		// fmt.Println(imgBin)
// 		// retMsg = message.NewSendingMessage().Append(message.NewText(out))
// 		retMsg = message.NewSendingMessage().Append(message.Image NewImage(imgBin))
// 		util.WriteFile("./img.jpg", imgBin)
// 		return retMsg, err
// 	}
// 	retMsg = message.NewSendingMessage().Append(message.NewText(out))
// 	return
// }
