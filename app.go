package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/pkg/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	// "gitee.com/lyhuilin/QN/global"
	"gitee.com/lyhuilin/QN/bot"
	"gitee.com/lyhuilin/QN/constvar"
	"gitee.com/lyhuilin/QN/env"
	"gitee.com/lyhuilin/QN/utils"

	_ "gitee.com/lyhuilin/QN/modules/autoreply"
	_ "gitee.com/lyhuilin/QN/modules/huiguangbo"
	_ "gitee.com/lyhuilin/QN/modules/logging"
	_ "gitee.com/lyhuilin/QN/modules/system"
)

var (
	cfg     = pflag.StringP("config", "c", "", "配置文件地址")
	version = pflag.BoolP("version", "v", false, "程序版本号")
	workdir = pflag.StringP("workdir", "w", "", "QN 工作目录")
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("run time panic: %v", err)
			if viper.GetBool("recover_restart_enable") {
				log.Infof("recover(restart)")
				if err := restart(); err != nil {
					log.Errorf(err, "recover(restart)")
				}
			}
		}
	}()
	defer func() {
		log.Infof("QN退出(Exit)")
		fmt.Scanln()
		if viper.GetBool("recover_restart_enable") {
			log.Infof("Exit(restart)")
			if err := restart(); err != nil {
				log.Errorf(err, "Exit(restart)")
			}
		}
	}()

	pflag.Parse()
	if *version {
		fmt.Println(constvar.APP_VERSION)
		return
	}

	fmt.Printf("%s (%v) \n%s\n", constvar.APP_NAME, constvar.APP_VERSION, constvar.APPDesc())
	time.Sleep(time.Second)

	if len(*workdir) > 0 {
		env.SetWorkdir(*workdir)
	}

	// 初始化配置信息
	if err := config.Init(*cfg); err != nil {
		log.Errorf(err, "QN配置，出错了")
		panic(err)
	}
	log.Infof("QN配置 初始化完成")

	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	botProtocol := bot.GetProtocol()
	log.Infof("使用协议: %s(%v)", viper.GetString("bot.use_protocol"), botProtocol)
	bot.UseProtocol(botProtocol)

	// 登录
	if err := bot.Login(); err != nil {
		log.Errorf(err, "登录出错了")
	} else {
		bot.SaveToken()
	}
	// 刷新好友列表，群列表
	bot.RefreshList()

	go func() {
		for {
			botid := fmt.Sprintf("%d", bot.Instance.Uin)
			res, err := utils.UpdateRobotStatToMyAdmin(botid, bot.Instance.Online.Load())
			if err != nil {
				log.Errorf(err, "更新机器人 %s 状态到myAdmin，出错啦(如未运行myAdmin，请删除相关配置或清空) %s", botid, res)
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Minute)

			if !bot.Instance.Online.Load() {
				botid := fmt.Sprintf("%d", bot.Instance.Uin)
				if viper.GetBool("recover_restart_enable") {
					log.Infof("检测到机器人 %s 已掉线，尝试重新登录!", botid)
					// 登录
					if err := bot.Login(); err != nil {
						log.Errorf(err, "登录出错了")
					} else {
						bot.SaveToken()
					}
					// 刷新好友列表，群列表
					bot.RefreshList()
				} else {
					log.Infof("检测到机器人 %s 已掉线，放弃尝试重新登录(如需尝试，请打开 recover_restart_enable 配置 true)!", botid)
				}
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	// <-global.SetupMainSignalHandler()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	bot.Stop()
}
