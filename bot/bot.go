package bot

import (
	"bufio"
	// "bytes"
	// "fmt"
	// "image"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	// asc2art "github.com/yinghau76/go-ascii-art"

	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/pkg/qr"
	"gitee.com/lyhuilin/util"

	// qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	// "github.com/tuotoo/qrcode"
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
	log.Infof("用户交流群: 153690156")
	log.Infof("初始化:%d", viper.GetInt64("bot.account"))
	Instance = &Bot{
		client.NewClient(
			viper.GetInt64("bot.account"),
			viper.GetString("bot.password"),
		),
		false,
	}
	if b, _ := util.IsFileExist("./device.json"); !b {
		log.Warnf("虚拟设备信息不存在, 将自动生成随机设备.")
		GenRandomDevice()
	}
	bytes, err := util.ReadFile("./device.json")
	if err != nil {
		log.Fatalf(err, "读取虚拟设备信息 device.json 失败")
	}
	err = client.SystemDeviceInfo.ReadJson(bytes)
	if err != nil {
		log.Fatalf(err, "加载虚拟设备信息 device.json 失败")
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
	case "QiDian":
		return QiDian
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
		log.Warn("虚拟设备信息已存在，放弃重新生成.")
	}
	err := ioutil.WriteFile("device.json", client.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	if err != nil {
		log.Errorf(err, "虚拟设备信息写入 device.json 失败 ")
	}
}

// SaveToken 会话缓存
func SaveToken() {
	AccountToken := Instance.GenToken()
	_ = os.WriteFile("session.token", AccountToken, 0o644)
}

// Login 登录
func Login() error {
	// 存在token缓存的情况快速恢复会话
	if exist, _ := util.IsFileExist("./session.token"); exist {
		log.Infof("检测到会话缓存, 尝试快速恢复登录")
		token, err := os.ReadFile("./session.token")
		if err == nil {
			r := binary.NewReader(token)
			cu := r.ReadInt64()
			if Instance.Uin != 0 {
				if cu != Instance.Uin {
					log.Warnf("警告: 配置文件内的QQ号 (%v) 与会话缓存内的QQ号 (%v) 不相同", Instance.Uin, cu)
					log.Warnf("1. 使用会话缓存继续.")
					log.Warnf("2. 删除会话缓存并退出程序.")
					log.Warnf("请选择: (5秒后自动选1)")
					text := readLineTimeout(time.Second*5, "1")
					if text == "2" {
						_ = os.Remove("session.token")
						log.Infof("会话缓存已删除.")
						os.Exit(0)
					}
				}
			}
			if err = Instance.TokenLogin(token); err != nil {
				_ = os.Remove("session.token")
				log.Warnf("恢复会话失败: %v , 尝试使用正常流程登录.", err)
				time.Sleep(time.Second)
				Instance.Disconnect()
				Instance.Release()

				Instance.QQClient = client.NewClient(
					viper.GetInt64("bot.account"),
					viper.GetString("bot.password"),
				)
			} else {
				log.Infof("快速恢复登录成功")
				return nil
			}
		}
	}
	// 不存在token缓存 走正常流程
	print(Instance.Uin)
	if Instance.Uin != 0 {
		// 有账号就先普通登录
		return CommonLogin()
	} else {
		// 没有账号就扫码登录
		return QrcodeLogin()
	}

	// resp, err := Instance.Login()
	// console := bufio.NewReader(os.Stdin)

	// for {
	// 	if err != nil {
	// 		// logger.WithError(err).Fatal("unable to login")
	// 		log.Errorf(err, "unable to login")
	// 	}

	// 	var text string
	// 	if !resp.Success {
	// 		switch resp.Error {

	// 		case client.NeedCaptcha:
	// 			img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))
	// 			fmt.Println(asc2art.New("image", img).Art)
	// 			fmt.Print("please input captcha: ")
	// 			text, _ := console.ReadString('\n')
	// 			resp, err = Instance.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)
	// 			continue

	// 		case client.UnsafeDeviceError:
	// 			fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
	// 			os.Exit(4)

	// 		case client.SMSNeededError:
	// 			fmt.Println("device lock enabled, Need SMS Code")
	// 			fmt.Printf("Send SMS to %s ? (yes)", resp.SMSPhone)
	// 			t, _ := console.ReadString('\n')
	// 			t = strings.TrimSpace(t)
	// 			if t != "yes" {
	// 				os.Exit(2)
	// 			}
	// 			if !Instance.RequestSMS() {
	// 				// logger.Warnf("unable to request SMS Code")
	// 				log.Warnf("unable to request SMS Code")
	// 				os.Exit(2)
	// 			}
	// 			log.Warn("please input SMS Code: ")
	// 			text, _ = console.ReadString('\n')
	// 			resp, err = Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
	// 			continue

	// 		case client.TooManySMSRequestError:
	// 			fmt.Printf("too many SMS request, please try later.\n")
	// 			os.Exit(6)

	// 		case client.SMSOrVerifyNeededError:
	// 			fmt.Println("device lock enabled, choose way to verify:")
	// 			fmt.Println("1. Send SMS Code to ", resp.SMSPhone)
	// 			fmt.Println("2. Scan QR Code")
	// 			fmt.Print("input (1,2):")
	// 			text, _ = console.ReadString('\n')
	// 			text = strings.TrimSpace(text)
	// 			switch text {
	// 			case "1":
	// 				if !Instance.RequestSMS() {
	// 					fmt.Println("unable to request SMS Code")
	// 					os.Exit(2)
	// 				}
	// 				fmt.Print("please input SMS Code: ")
	// 				text, _ = console.ReadString('\n')
	// 				resp, err = Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
	// 				continue
	// 			case "2":
	// 				fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
	// 				os.Exit(2)
	// 			default:
	// 				fmt.Println("invalid input")
	// 				os.Exit(2)
	// 			}

	// 		case client.SliderNeededError:
	// 			if client.SystemDeviceInfo.Protocol == client.AndroidPhone {
	// 				fmt.Println("Android Phone Protocol DO NOT SUPPORT Slide verify")
	// 				fmt.Println("please use other protocol")
	// 				os.Exit(2)
	// 			}
	// 			Instance.AllowSlider = false
	// 			Instance.Disconnect()
	// 			resp, err = Instance.Login()
	// 			continue

	// 		case client.OtherLoginError, client.UnknownLoginError:
	// 			log.Fatalf(nil, "login failed: %v", resp.ErrorMessage)
	// 		}

	// 	}

	// 	break
	// }

	// log.Infof("bot login: %s", Instance.Nickname)
}

// CommonLogin 普通账号密码登录
func CommonLogin() error {
	res, err := Instance.Login()
	if err != nil {
		return err
	}
	return loginResponseProcessor(res)
}

// QrcodeLogin 扫码登陆
func QrcodeLogin() error {
	rsp, err := Instance.FetchQRCode()
	if err != nil {
		return err
	}
	// fi, err := qrcode.Decode(bytes.NewReader(rsp.ImageData))
	// if err != nil {
	// 	return err
	// }
	_ = os.WriteFile("qrcode.png", rsp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	if Instance.Uin != 0 {
		log.Infof("请使用账号 %v 登录手机QQ扫描二维码 (qrcode.png) : ", Instance.Uin)
	} else {
		log.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
	}
	time.Sleep(time.Second)
	// qrcodeTerminal.New().Get(fi.Content).Print()
	qr.PrintQRCode(rsp.ImageData)

	s, err := Instance.QueryQRCodeStatus(rsp.Sig)
	if err != nil {
		return err
	}
	prevState := s.State
	for {
		time.Sleep(time.Second)
		s, _ = Instance.QueryQRCodeStatus(rsp.Sig)
		if s == nil {
			continue
		}
		if prevState == s.State {
			continue
		}
		prevState = s.State
		switch s.State {
		case client.QRCodeCanceled:
			log.Warnf("扫码被用户取消.")
		case client.QRCodeTimeout:
			log.Warnf("二维码过期")
		case client.QRCodeWaitingForConfirm:
			log.Infof("扫码成功, 请在手机端确认登录.")
		case client.QRCodeConfirmed:
			res, err := Instance.QRCodeLogin(s.LoginInfo)
			if err != nil {
				return err
			}
			return loginResponseProcessor(res)
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			// ignore
		}
	}
}

// ErrSMSRequestError SMS请求出错
var ErrSMSRequestError = errors.New("sms request error")
var console = bufio.NewReader(os.Stdin)

func readLine() (str string) {
	str, _ = console.ReadString('\n')
	str = strings.TrimSpace(str)
	return
}
func readLineTimeout(t time.Duration, de string) (str string) {
	r := make(chan string)
	go func() {
		select {
		case r <- readLine():
		case <-time.After(t):
		}
	}()
	str = de
	select {
	case str = <-r:
	case <-time.After(t):
	}
	return
}

// loginResponseProcessor 登录结果处理
func loginResponseProcessor(res *client.LoginResponse) error {
	var err error
	for {
		if err != nil {
			return err
		}
		if res.Success {
			return nil
		}
		var text string
		switch res.Error {
		case client.SliderNeededError:
			log.Warnf("登录需要滑条验证码, 请使用手机QQ扫描二维码以继续登录.")
			Instance.Disconnect()
			Instance.Release()
			Instance.QQClient = client.NewClientEmpty()
			return QrcodeLogin()
		case client.NeedCaptcha:
			log.Warnf("登录需要验证码.")
			_ = os.WriteFile("captcha.jpg", res.CaptchaImage, 0o644)
			log.Warnf("请输入验证码 (captcha.jpg)： (Enter 提交)")
			text = readLine()
			_ = os.Remove("captcha.jpg")
			res, err = Instance.SubmitCaptcha(text, res.CaptchaSign)
			continue
		case client.SMSNeededError:
			log.Warnf("账号已开启设备锁, 按 Enter 向手机 %v 发送短信验证码.", res.SMSPhone)
			readLine()
			if !Instance.RequestSMS() {
				log.Warnf("发送验证码失败，可能是请求过于频繁.")
				return errors.WithStack(ErrSMSRequestError)
			}
			log.Warn("请输入短信验证码： (Enter 提交)")
			text = readLine()
			res, err = Instance.SubmitSMS(text)
			continue
		case client.SMSOrVerifyNeededError:
			log.Warnf("账号已开启设备锁，请选择验证方式:")
			log.Warnf("1. 向手机 %v 发送短信验证码", res.SMSPhone)
			log.Warnf("2. 使用手机QQ扫码验证.")
			log.Warn("请输入(1 - 2) (将在10秒后自动选择2)：")
			text = readLineTimeout(time.Second*10, "2")
			if strings.Contains(text, "1") {
				if !Instance.RequestSMS() {
					log.Warnf("发送验证码失败，可能是请求过于频繁.")
					return errors.WithStack(ErrSMSRequestError)
				}
				log.Warn("请输入短信验证码： (Enter 提交)")
				text = readLine()
				res, err = Instance.SubmitSMS(text)
				continue
			}
			fallthrough
		case client.UnsafeDeviceError:
			log.Warnf("账号已开启设备锁，请前往 -> %v <- 验证后重启Bot.", res.VerifyUrl)
			log.Infof("按 Enter 或等待 5s 后继续....")
			readLineTimeout(time.Second*5, "")
			os.Exit(0)
		case client.OtherLoginError, client.UnknownLoginError, client.TooManySMSRequestError:
			msg := res.ErrorMessage
			if strings.Contains(msg, "版本") {
				msg = "密码错误或账号被冻结"
			}
			if strings.Contains(msg, "冻结") {
				log.Warnf("账号被冻结")
			}
			log.Warnf("登录失败: %v", msg)
			log.Infof("按 Enter 或等待 5s 后继续....")
			readLineTimeout(time.Second*5, "")
			os.Exit(0)
		}
	}
}

// RefreshList 刷新联系人
func RefreshList() {
	log.Info("开始刷新好友列表")
	err := Instance.ReloadFriendList()
	if err != nil {
		log.Errorf(err, "刷新好友列表失败")
	}
	log.Infof("拉取到 %d 好友", len(Instance.FriendList))
	log.Info("开始刷新群列表")
	err = Instance.ReloadGroupList()
	if err != nil {
		log.Errorf(err, "刷新群列表失败")
	}
	log.Infof("拉取到 %d 群", len(Instance.GroupList))
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

	log.Infof("初始化插件信息 ...")
	for _, mi := range modules {
		mi.Instance.Init()
	}
	for _, mi := range modules {
		mi.Instance.PostInit()
	}
	log.Info("所有插件初始化完毕")

	log.Info("注册插件服务函数 ...")
	for _, mi := range modules {
		mi.Instance.Serve(Instance)
	}
	log.Info("注册插件服务函数完毕")

	log.Info("启动插件 ...")
	for _, mi := range modules {
		go mi.Instance.Start(Instance)
	}
	log.Info("插件运行中")
}

// Stop 停止所有服务
// 调用此函数并不会使Bot离线
func Stop() {
	log.Warn("停止中 ...")
	wg := sync.WaitGroup{}
	for _, mi := range modules {
		wg.Add(1)
		mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	log.Info("已停止")
	modules = make(map[string]ModuleInfo)
}
