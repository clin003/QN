package huiguangbo

import (
	"time"
)

type HGBConf struct {
	GroupList   []GroupInfo   `yaml: "group_list"`
	SenderSleep time.Duration `yaml: "sender_sleep"`
}

type GroupInfo struct {
	Name   string `yaml:"group_name"`
	Id     int64  `yaml:"group_id"`
	IsFeed bool   `yaml:"send_msg_enable"`
}
