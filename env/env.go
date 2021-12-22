package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// 获取程序运行目录相当于linux的pwd
func Getwd() string {
	str, err := os.Getwd()
	if err != nil || len(str) <= 0 {
		fmt.Println(err)
		return ""
	}
	return str
}

// 获取绝对路径 tmpName：绝对/相对路径，defaultName:默认路径
func GetAbsPath(tmpName, defaultName string) string {
	name := tmpName
	if len(name) <= 0 {
		name = defaultName
	}
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(Getwd(), name)
}

//https://www.cnblogs.com/malukang/p/12907945.html
// GetRel(uploadDir,uploadFilePath)
// 获取相对路径
func GetRel(basepath, targpath string) string {
	s, err := filepath.Rel(basepath, targpath)
	if err != nil {
		return targpath
	}
	return s
}
