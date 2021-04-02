// Copyright 2015 The HLTYopenAPI(baicai) Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
package huiguangbo

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"gitee.com/lyhuilin/log"
	"gitee.com/lyhuilin/open_api/model/feedmsg"
	"github.com/gorilla/websocket"
)

var wsUrlStr string
var serverUrl string

// 解析服务器地址为ws地址格式
func parseWsServerUrl(wsServerUrl, channel string) (retText string) {
	// ws服务器地址
	// wsServerUrl := "https://api.lyhuilin.com"
	// ws传输数据channel token
	// channel := "weipinhui"
	scheme := "ws"
	host := "127.0.0.1:8080"

	path := fmt.Sprintf("/ws/live/%s", channel)
	if strings.HasPrefix(wsServerUrl, "https://") {
		scheme = "wss"
		host = strings.Replace(wsServerUrl, "https://", "", 1)
	} else if strings.HasPrefix(wsServerUrl, "http://") {
		scheme = "ws"
		host = strings.Replace(wsServerUrl, "http://", "", 1)
	}

	u := url.URL{Scheme: scheme, Host: host, Path: path}
	retText = u.String()
	// fmt.Printf("connecting to %s\n", retText)
	return
}

// 启动wsClient
// go wsClientStart(wsServerUrl, channel)
func wsClientStart(wsServerUrl, channel string) {
	serverUrl = wsServerUrl
	wsUrlStr = parseWsServerUrl(wsServerUrl, channel)
	go wsClientStartService()
}

// 启动wsClient服务并保持
func wsClientStartService() {
	wsClientConn, _, err := websocket.DefaultDialer.Dial(wsUrlStr, nil)
	if err != nil {
		log.Errorf(err, "dial:%s", serverUrl)
		time.Sleep(30 * time.Second)
		go wsClientStartService()
		return
	}
	defer wsClientConn.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			var msg feedmsg.FeedRichMsgModel
			err := wsClientConn.ReadJSON(&msg)
			if err != nil {
				// fmt.Printf("read:%v\n", err)
				log.Errorf(err, "read")
				return
			}
			go sendMsg(msg)
			fmt.Printf("recv(%s):%v\n", msg.Msgtype, msg.Text.Content)
			// log.Infof("recv(%s):%v\n", msg.Msgtype, msg.Text.Content)
		}
	}()

	//os.Interrupt 表示中断
	//os.Kill 杀死退出进程
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			time.Sleep(30 * time.Second)
			go wsClientStartService()
			return

		case <-interrupt:
			// fmt.Println("interrupt")
			log.Debug("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := wsClientConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				log.Errorf(err, "write close")
				return
			}
			select {
			case <-done:
				// case <-time.After(time.Second):
			}
			return
		}
	}
}
