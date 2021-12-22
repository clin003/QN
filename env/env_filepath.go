package env

import (
	// "path"
	"path/filepath"

	"github.com/spf13/viper"
)

// 获取上传目录路径(前缀)(用于路由配置上传路径)
func GetEnvUploadPreFix() string {
	return "upload"
}

// // 获取上传文件路径的路由地址（相对路径）
// func GetUploadFileRouterPath(rootDir, filePath string) string {
// 	relPath, _ := filepath.Rel(rootDir, filePath)
// 	return path.Join(path.Base(rootDir), relPath)
// }
// 获取上传文件路径的路由地址（相对路径）
func GetUploadFileRouterPath(rootDir, filePath string) string {
	relPath, _ := filepath.Rel(rootDir, filePath)
	fileBash := filepath.Base(rootDir)

	retPath := filepath.Join(fileBash, relPath)
	retPath = filepath.Clean(retPath)
	// ToSlash函数将path中的路径分隔符替换为斜杠（’/’）并返回替换结果，多个路径分隔符会替换为多个斜杠。
	retPath = filepath.ToSlash(retPath)
	return retPath
	// return filepath.ToSlash(relPath)
	// if filepath.Separator == '/' {
	// 	return relPath
	// }
	// return strings.ReplaceAll(relPath, string(filepath.Separator), "/")
	// return retPath
}

// 获取上传目录路径(绝对路径)
func GetEnvUploadFileDir() string {
	name := viper.GetString("upload_filepath")
	return GetAbsPath(name, "upload")
}

// 获取sqlite数据库文件路径（绝对路径）
func GetDBFilePath() string {
	name := viper.GetString("db.filename")
	return GetAbsPath(name, "db.kns")
}
