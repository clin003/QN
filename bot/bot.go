package bot

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	asc2art "github.com/yinghau76/go-ascii-art"

	// "github.com/Logiase/MiraiGo-Template/config"
	// "gitee.com/lyhuilin/QN/config"
	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/util"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/spf13/viper"
	// "github.com/sirupsen/logrus"
)

// Bot 全局 Bot
type Bot struct {
	*client.QQClient

	start bool
}

// Instance Bot 实例
var Instance *Bot

// var logger = logrus.WithField("bot", "internal")

// Init 快速初始化
// 使用 config.GlobalConfig 初始化账号
// 使用 ./device.json 初始化设备信息
func Init() {
	// config.Init()
	// fmt.Println("Init")
	// log.Info("Init")
	log.Infof("Init:%d", viper.GetInt64("bot.account"))
	// log.Infof("Init:%d", config.GlobalConfig.GetInt64("bot.account"))
	Instance = &Bot{
		client.NewClient(
			viper.GetInt64("bot.account"),
			viper.GetString("bot.password"),
		),
		false,
	}
	if b, _ := util.IsFileExist("./device.json"); !b {
		GenRandomDevice()
	}
	bytes, err := util.ReadFile("./device.json")
	if err != nil {
		log.Fatalf(err, "device.json error")
	}
	err = client.SystemDeviceInfo.ReadJson(bytes)
	if err != nil {
		// logger.WithError(err).Panic("device.json error")
		log.Fatalf(err, "device.json error")
	}
}

// GetProtocol 获取配置文件中协议项，并转换为协议const
func GetProtocol() protocol {
	switch viper.GetString("bot.use_protocol") {
	case "AndroidPhone":
		return AndroidPhone
	case "IPad":
		return IPad
	case "AndroidWatch":
		return AndroidWatch
	case "MacOS":
		return MacOS
	}
	return IPad
}

// InitBot 使用 account password 进行初始化账号
func InitBot(account int64, password string) {
	Instance = &Bot{
		client.NewClient(account, password),
		false,
	}
}

// UseDevice 使用 device 进行初始化设备信息
func UseDevice(device []byte) error {
	return client.SystemDeviceInfo.ReadJson(device)
}

// GenRandomDevice 生成随机设备信息
func GenRandomDevice() {
	client.GenRandomDevice()
	b, _ := util.IsFileExist("./device.json")
	if b {
		log.Warn("device.json exists, will not write device to file")
	}
	err := ioutil.WriteFile("device.json", client.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	if err != nil {
		log.Errorf(err, "unable to write device.json")
	}
}

// Login 登录
func Login() {
	resp, err := Instance.Login()
	console := bufio.NewReader(os.Stdin)

	for {
		if err != nil {
			// logger.WithError(err).Fatal("unable to login")
			log.Errorf(err, "unable to login")
		}

		var text string
		if !resp.Success {
			switch resp.Error {

			case client.NeedCaptcha:
				img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))
				fmt.Println(asc2art.New("image", img).Art)
				fmt.Print("please input captcha: ")
				text, _ := console.ReadString('\n')
				resp, err = Instance.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)
				continue

			case client.UnsafeDeviceError:
				fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
				os.Exit(4)

			case client.SMSNeededError:
				fmt.Println("device lock enabled, Need SMS Code")
				fmt.Printf("Send SMS to %s ? (yes)", resp.SMSPhone)
				t, _ := console.ReadString('\n')
				t = strings.TrimSpace(t)
				if t != "yes" {
					os.Exit(2)
				}
				if !Instance.RequestSMS() {
					// logger.Warnf("unable to request SMS Code")
					log.Warnf("unable to request SMS Code")
					os.Exit(2)
				}
				log.Warn("please input SMS Code: ")
				text, _ = console.ReadString('\n')
				resp, err = Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
				continue

			case client.TooManySMSRequestError:
				fmt.Printf("too many SMS request, please try later.\n")
				os.Exit(6)

			case client.SMSOrVerifyNeededError:
				fmt.Println("device lock enabled, choose way to verify:")
				fmt.Println("1. Send SMS Code to ", resp.SMSPhone)
				fmt.Println("2. Scan QR Code")
				fmt.Print("input (1,2):")
				text, _ = console.ReadString('\n')
				text = strings.TrimSpace(text)
				switch text {
				case "1":
					if !Instance.RequestSMS() {
						fmt.Println("unable to request SMS Code")
						os.Exit(2)
					}
					fmt.Print("please input SMS Code: ")
					text, _ = console.ReadString('\n')
					resp, err = Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
					continue
				case "2":
					fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
					os.Exit(2)
				default:
					fmt.Println("invalid input")
					os.Exit(2)
				}

			case client.SliderNeededError:
				if client.SystemDeviceInfo.Protocol == client.AndroidPhone {
					fmt.Println("Android Phone Protocol DO NOT SUPPORT Slide verify")
					fmt.Println("please use other protocol")
					os.Exit(2)
				}
				Instance.AllowSlider = false
				Instance.Disconnect()
				resp, err = Instance.Login()
				continue

			case client.OtherLoginError, client.UnknownLoginError:
				log.Fatalf(nil, "login failed: %v", resp.ErrorMessage)
			}

		}

		break
	}

	log.Infof("bot login: %s", Instance.Nickname)
}

// RefreshList 刷新联系人
func RefreshList() {
	log.Info("start reload friends list")
	err := Instance.ReloadFriendList()
	if err != nil {
		log.Errorf(err, "unable to load friends list")
	}
	log.Infof("load %d friends", len(Instance.FriendList))
	log.Info("start reload groups list")
	err = Instance.ReloadGroupList()
	if err != nil {
		log.Errorf(err, "unable to load groups list")
	}
	log.Infof("load %d groups", len(Instance.GroupList))
	// for _, v := range Instance.GroupList {
	// 	// fmt.Println(k, v.Code)
	// 	fmt.Printf("Name(%s),Code(%d),Uin(%d)\n", v.Name, v.Code, v.Uin)
	// }
}

// StartService 启动服务
// 根据 Module 生命周期 此过程应在Login前调用
// 请勿重复调用
func StartService() {
	if Instance.start {
		return
	}

	Instance.start = true

	log.Infof("initializing modules ...")
	for _, mi := range modules {
		mi.Instance.Init()
	}
	for _, mi := range modules {
		mi.Instance.PostInit()
	}
	log.Info("all modules initialized")

	log.Info("registering modules serve functions ...")
	for _, mi := range modules {
		mi.Instance.Serve(Instance)
	}
	log.Info("all modules serve functions registered")

	log.Info("starting modules tasks ...")
	for _, mi := range modules {
		go mi.Instance.Start(Instance)
	}
	log.Info("tasks running")
}

// Stop 停止所有服务
// 调用此函数并不会使Bot离线
func Stop() {
	log.Warn("stopping ...")
	wg := sync.WaitGroup{}
	for _, mi := range modules {
		wg.Add(1)
		mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	log.Info("stopped")
	modules = make(map[string]ModuleInfo)
}
