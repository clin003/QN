package main

import (
	"fmt"
	"os"
	"os/signal"

	"gitee.com/lyhuilin/QN/bot"

	"gitee.com/lyhuilin/config"
	"gitee.com/lyhuilin/log"

	"github.com/spf13/pflag"

	_ "gitee.com/lyhuilin/QN/modules/logging"
)

var (
	cfg = pflag.StringP("config", "c", "", "config file path")
)

func main() {
	defer func() {
		fmt.Scanln()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Debugf("run time panic:%v\n", err)
		}
	}()
	// cfg := pflag.StringP("config", "c", "", "config file path")
	pflag.Parse()

	//init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	log.Info("config 初始化完成")

	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	bot.Login()

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
