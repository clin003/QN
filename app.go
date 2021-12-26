package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"

	"gitee.com/lyhuilin/QN/bot"
	"gitee.com/lyhuilin/QN/constvar"
	"gitee.com/lyhuilin/QN/env"
	"gitee.com/lyhuilin/QN/global"
	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"

	_ "gitee.com/lyhuilin/QN/modules/autoreply"
	_ "gitee.com/lyhuilin/QN/modules/huiguangbo"
	_ "gitee.com/lyhuilin/QN/modules/logging"
)

var (
	cfg     = pflag.StringP("config", "c", "", "配置文件地址")
	version = pflag.BoolP("version", "v", false, "程序版本号")
	workdir = pflag.StringP("workdir", "w", "", "QN 工作目录")
)

func main() {
	// // http://www.network-science.de/ascii/
	// // https://www.jianshu.com/p/fca56d635091
	// fmt.Printf(`

	//     QQQQQQQQQ      NNNNNNNN        NNNNNNNN
	//   QQ:::::::::QQ    N:::::::N       N::::::N
	// QQ:::::::::::::QQ  N::::::::N      N::::::N
	// Q:::::::QQQ:::::::Q N:::::::::N     N::::::N
	// Q::::::O   Q::::::Q N::::::::::N    N::::::N
	// Q:::::O     Q:::::Q N:::::::::::N   N::::::N
	// Q:::::O     Q:::::Q N:::::::N::::N  N::::::N
	// Q:::::O     Q:::::Q N::::::N N::::N N::::::N
	// Q:::::O     Q:::::Q N::::::N  N::::N:::::::N
	// Q:::::O     Q:::::Q N::::::N   N:::::::::::N
	// Q:::::O  QQQQ:::::Q N::::::N    N::::::::::N
	// Q::::::O Q::::::::Q N::::::N     N:::::::::N
	// Q:::::::QQ::::::::Q N::::::N      N::::::::N
	// QQ::::::::::::::Q  N::::::N       N:::::::N
	//   QQ:::::::::::Q   N::::::N        N::::::N
	//     QQQQQQQQ::::QQ NNNNNNNN         NNNNNNN
	//             Q:::::Q
	//              QQQQQQ                        v%s

	// `, constvar.APP_VERSION)

	defer func() {
		fmt.Scanln()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Debugf("run time panic:%v\n", err)
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
	fmt.Println("使用协议: %s(%v)", viper.GetString("bot.use_protocol"), botProtocol)
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
		r := gin.Default()
		r.GET("/ping",
			func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
					"Online":  bot.Instance.Online.Load(),
					"data":    constvar.APP_VERSION,
				})
			},
		)
		r.GET("/",
			func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "HelloWorld",
					"Online":  bot.Instance.Online.Load(),
					"data":    constvar.APPDesc(),
				})
			},
		)
		r.Run(viper.GetString("addr")) // 监听并在 0.0.0.0:8080 上启动服务
	}()

	<-global.SetupMainSignalHandler()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
