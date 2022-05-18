// Package config
/*
读取操作系统本地配置文件，获取服务器地址、存储地址及其他配置信息
*/
package config

import (
	"fmt"
	"os/user"
	"path"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/aboutuser"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
)

// GetConfigSet ~/.sabrefig/config
// defult JAVA_OPTS="-server -Xms1024M -Xmx1024M -Xss512k"
func GetConfigSet() (*sabstruct.Config, error) {
	currentUser, err := user.Current()
	if err != nil {
		return &sabstruct.Config{}, err
	}
	rootDir, getUserHomeDirErr := aboutuser.GetUserHomeDir(currentUser.Username)
	if getUserHomeDirErr != nil {
		return &sabstruct.Config{}, getUserHomeDirErr
	}
	configPath := "/.sabrefig"
	configFile := "/config"

	sabreConfigFile := path.Join(rootDir + configPath + configFile)
	if commontools.IsFileExist(sabreConfigFile) {
		return &sabstruct.Config{}, fmt.Errorf("用户%s的家目录%s下，无sabre的配置文件\n",
			currentUser.Username, sabreConfigFile)
	}
	var s sabstruct.Config

	yamlFmt, yamlFmtErr := yamlfmt.YamlFmt(sabreConfigFile, s)
	if yamlFmtErr != nil {
		return &sabstruct.Config{}, fmt.Errorf("sabre 配置文件解析失败, %s\n", yamlFmtErr)
	}
	return yamlFmt, nil
}

//GetApiServerUrl 获取API网关的地址
func GetApiServerUrl() (string, error) {
	sabreConfig, err := GetConfigSet()
	if err != nil {
		return "", err
	}

	return sabreConfig.ApiServer, nil

}
