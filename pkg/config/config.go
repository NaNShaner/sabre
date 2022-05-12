// Package config
/*
读取操作系统本地配置文件，获取服务器地址、存储地址及其他配置信息
*/
package config

import (
	"os/user"
	"sabre/pkg/util/aboutuser"
)

// GetConfigSet ~/.sabrefig/config
// defult JAVA_OPTS="-server -Xms1024M -Xmx1024M -Xss512k"
func GetConfigSet() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}
	rootDir, err := aboutuser.GetUserHomeDir(currentUser.Username)
	if err != nil {
		return ""
	}
	configPaht := ".sabrefig"
	configFile := "config"
	return rootDir + configPaht + configFile
}
