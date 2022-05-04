// Package config
/*
读取操作系统本地配置文件，获取服务器地址、存储地址及其他配置信息
*/
package config

import (
	"os"
)

//type str struct {
//	s []string
//	seq string
//}

// GetConfigSet ~/.sabrefig/config
func GetConfigSet() string {
	homeDir := os.Getenv("HOME")
	configPaht := ".sabrefig"
	configFile := "config"
	seq := "/"
	return seq + homeDir + configPaht + configFile
}
