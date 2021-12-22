package env

import (
	"path/filepath"
	"sync"
)

var once sync.Once
var workdir string

func GetWorkdir() string {
	return workdir
}
func SetWorkdir(dir string) {
	once.Do(func() {
		if len(dir) > 0 {
			workdir = dir
		} else {
			workdir = Getwd()
		}
	})
}

func GetDeviceFilePath() string {
	return filepath.Join(GetWorkdir(), "./device.json")
}

func GetSessionFilePath() string {
	return filepath.Join(GetWorkdir(), "session.token")
}

func GetQrcodeFilePath() string {
	return filepath.Join(GetWorkdir(), "qrcode.png")
}
func GetCaptchaFilePath() string {
	return filepath.Join(GetWorkdir(), "captcha.jpg")
}
