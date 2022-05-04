// Package config
/*
读取操作系统本地配置文件，获取服务器地址、存储地址及其他配置信息
*/
package config

import (
	"awesomeProject/pkg/yamlfmt"
	"os"
)

//type str struct {
//	s []string
//	seq string
//}

type ServerConfig struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   `json:"metadata"`
	Spec       `json:"spec"`
}

type Metadata struct {
	Namespace string
}

type Spec struct {
	MidServerUrl string
}


// GetConfigSet ~/.sabrefig/config
func GetConfigSet() string {
	homeDir := os.Getenv("HOME")
	configPaht := ".sabrefig"
	configFile := "config"
	seq := "/"
	return seq+homeDir+configPaht+configFile
}

func GetConfigYaml(fpath string) ([]byte, error) {
	var sc ServerConfig
	yamlFmt, err := yamlfmt.YamlFmt(fpath, sc)
	if err != nil {
		return nil, err
	}
	return yamlFmt, nil
}

