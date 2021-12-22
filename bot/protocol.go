package bot

import "github.com/Mrs4s/MiraiGo/client"

type protocol int

// AndroidPhone ClientProtocol = 1
// IPad         ClientProtocol = 2
// AndroidWatch ClientProtocol = 3
// MacOS        ClientProtocol = 4
// QiDian       ClientProtocol = 5
const (
	AndroidPhone = protocol(client.AndroidPhone) //1
	IPad         = protocol(client.IPad)         //2
	AndroidWatch = protocol(client.AndroidWatch) //3
	MacOS        = protocol(client.MacOS)        //4
	QiDian       = protocol(client.QiDian)
)

// UseProtocol 使用协议
// 不同协议会有部分功能无法使用
// 默认为 AndroidPad
func UseProtocol(p protocol) {
	client.SystemDeviceInfo.Protocol = client.ClientProtocol(p)
}
