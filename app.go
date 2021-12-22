package main

import (
	"fmt"
	"os"
	"os/signal"

	"gitee.com/lyhuilin/QN/bot"

	"gitee.com/lyhuilin/config"
	"gitee.com/lyhuilin/log"

	"github.com/spf13/pflag"

	_ "gitee.com/lyhuilin/QN/modules/autoreply"
	_ "gitee.com/lyhuilin/QN/modules/huiguangbo"
	_ "gitee.com/lyhuilin/QN/modules/logging"
)

var (
	cfg = pflag.StringP("config", "c", "", "config file path")
)

func main() {
	// http://www.network-science.de/ascii/
	// https://www.jianshu.com/p/fca56d635091
	fmt.Printf(`
                                            
     QQQQQQQQQ      NNNNNNNN        NNNNNNNN
   QQ:::::::::QQ    N:::::::N       N::::::N
 QQ:::::::::::::QQ  N::::::::N      N::::::N
Q:::::::QQQ:::::::Q N:::::::::N     N::::::N
Q::::::O   Q::::::Q N::::::::::N    N::::::N
Q:::::O     Q:::::Q N:::::::::::N   N::::::N
Q:::::O     Q:::::Q N:::::::N::::N  N::::::N
Q:::::O     Q:::::Q N::::::N N::::N N::::::N
Q:::::O     Q:::::Q N::::::N  N::::N:::::::N
Q:::::O     Q:::::Q N::::::N   N:::::::::::N
Q:::::O  QQQQ:::::Q N::::::N    N::::::::::N
Q::::::O Q::::::::Q N::::::N     N:::::::::N
Q:::::::QQ::::::::Q N::::::N      N::::::::N
 QQ::::::::::::::Q  N::::::N       N:::::::N
   QQ:::::::::::Q   N::::::N        N::::::N
     QQQQQQQQ::::QQ NNNNNNNN         NNNNNNN
             Q:::::Q                        
              QQQQQQ                        v%s
                                                                                                               
`, API_VERSION)
	// return
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
	log.Info("QN配置 初始化完成")

	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	botProtocol := bot.GetProtocol()
	fmt.Printf("使用协议: %v", botProtocol)
	bot.UseProtocol(botProtocol)

	// 登录
	bot.Login()

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
