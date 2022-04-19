package bot

import (
	"bufio"

	"gitee.com/lyhuilin/QN/env"

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
	deviceFilePath := env.GetDeviceFilePath()
	if b, _ := util.IsFileExist(deviceFilePath); !b {
		log.Warnf("虚拟设备信息不存在, 将自动生成随机设备.")
		GenRandomDevice()
	}
	bytes, err := util.ReadFile(deviceFilePath)
	if err != nil {
		log.Fatalf(err, "读取虚拟设备信息 %s 失败", deviceFilePath)
	}
	err = client.SystemDeviceInfo.ReadJson(bytes)
	if err != nil {
		log.Fatalf(err, "加载虚拟设备信息 %s 失败", deviceFilePath)
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
	return Unset
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
	deviceFilePath := env.GetDeviceFilePath()
	b, _ := util.IsFileExist(deviceFilePath)
	if b {
		log.Warn("虚拟设备信息已存在，放弃重新生成.")
	}
	err := ioutil.WriteFile(deviceFilePath, client.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	if err != nil {
		log.Errorf(err, "虚拟设备信息写入 %s 失败 ", deviceFilePath)
	}
}

// SaveToken 会话缓存
func SaveToken() {
	AccountToken := Instance.GenToken()
	sessionFilePath := env.GetSessionFilePath()
	_ = os.WriteFile(sessionFilePath, AccountToken, 0o644)
}

// Login 登录
func Login() error {
	// 存在token缓存的情况快速恢复会话
	sessionFilePath := env.GetSessionFilePath()
	if exist, _ := util.IsFileExist(sessionFilePath); exist {
		log.Infof("检测到会话缓存, 尝试快速恢复登录")
		token, err := os.ReadFile(sessionFilePath)
		if err == nil {
			r := binary.NewReader(token)
			cu := r.ReadInt64()
			if Instance.Uin != 0 {
				if cu != Instance.Uin {
					log.Warnf("警告: 配置文件内的QQ号 (%v) 与会话缓存内的QQ号 (%v) 不相同", Instance.Uin, cu)
					log.Warnf("1. 使用会话缓存继续.")
					log.Warnf("2. 删除会话缓存并退出程序.")
					log.Warnf("请选择: (5秒后自动选1)")
					text := readLineTimeout(time.Second*30, "1")
					if text == "2" {
						_ = os.Remove(sessionFilePath)
						log.Infof("会话缓存已删除.")
						os.Exit(0)
					}
				}
			}
			if err = Instance.TokenLogin(token); err != nil {
				_ = os.Remove(sessionFilePath)
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
	qrcodeFilePath := env.GetQrcodeFilePath()
	_ = os.WriteFile(qrcodeFilePath, rsp.ImageData, 0o644)
	defer func() { _ = os.Remove(qrcodeFilePath) }()
	if Instance.Uin != 0 {
		log.Infof("请使用账号 %v 登录手机QQ扫描二维码 (%s) : ", Instance.Uin, qrcodeFilePath)
	} else {
		log.Infof("请使用手机QQ扫描二维码 (%s) : ", qrcodeFilePath)
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
			captchaFilePath := env.GetCaptchaFilePath()
			_ = os.WriteFile(captchaFilePath, res.CaptchaImage, 0o644)
			log.Warnf("请输入验证码 (%s)： (Enter 提交)", captchaFilePath)
			text = readLine()
			_ = os.Remove(captchaFilePath)
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
			log.Warn("请输入(1 - 2) (将在60秒后自动选择2)：")
			text = readLineTimeout(time.Second*60, "2")
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
			log.Infof("按 Enter 或等待 60s 后继续....")
			readLineTimeout(time.Second*60, "")
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
			log.Infof("按 Enter 或等待 60s 后继续....")
			readLineTimeout(time.Second*60, "")
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
	log.Infof("共加载 %d 个好友", len(Instance.FriendList))
	log.Info("开始刷新群列表")
	err = Instance.ReloadGroupList()
	if err != nil {
		log.Errorf(err, "刷新群列表失败")
	}
	log.Infof("共加载 %d 个群", len(Instance.GroupList))

	log.Infof("共加载 %d 个频道", len(Instance.GuildService.Guilds))
	for _, v := range Instance.GuildService.Guilds {
		log.Infof("%s(ID:%d Code:%d)", v.GuildName, v.GuildId, v.GuildCode)
		for _, vv := range v.Channels {
			log.Infof("%s(ID:%d  EventTime:%d)", vv.ChannelName, vv.ChannelId, vv.EventTime)
		}
	}
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
