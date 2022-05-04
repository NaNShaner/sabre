/*
@spec: 中间件的安装部署
*/

package getdeploypkg

import (
	"awesomeProject/pkg/config"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	PkgLocalPathForLinux = "./" // 默认的安装包以及配置文件的存放地方，可以根据yam文件覆盖
)

type DeployAction struct {
	Install   string // 部署
	ReInstall string // 重装
	apply     string // 配置修改
	OffLine   string // 下线
	Start     string //	启动
	Stop      string //	停止
	Restop    string // 重启
}

type MiddInfo struct {
	NameSpace  string
	MidType    string   // 中间件的类型
	MidVersion string   // 中间件的版本
	MidRunType []string // 运行模式，集群、主备、冷备等等
}

type DeployPkg struct {
	timer        time.Time // 执行时间
	PkgFromUrl   string    // 安装所需文件连接
	ConfigFile   string    // 安装所需配置文件
	MiddInfo               // 中间件相关信息
	DeployAction           // 执行动作
}

// GetDeployPkgFromUrl 从服务端获取安装包或者配置文件等
func (u *DeployPkg) GetDeployPkgFromUrl() (string, error) {
	pkgPath := PkgLocalPathForLinux
	pkgUrl := u.PkgFromUrl
	fileName := path.Base(pkgUrl)
	pkgBasPath := pkgPath + fileName
	fmt.Printf("%s", pkgUrl)
	res, err := http.Get(pkgUrl)
	if err != nil {
		fmt.Printf("A error occurred! %s", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	if IsFileExist(pkgBasPath) {
		file, err := os.Create(pkgBasPath)
		if err != nil {
			panic(err)
		}
		// 获得文件的writer对象
		writer := bufio.NewWriter(file)

		//written, err := io.Copy(writer, reader)
		_, err = io.Copy(writer, reader)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("Total length: %d", written)
		return fileName, nil
	}
	return "", fmt.Errorf("%s 文件已存在\n", pkgBasPath)

}

// IsFileExist 判断文件本地是否已经存在
func IsFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return true
	}
	fmt.Printf("文件 %s 已存在， %v", info.Name(), err)
	return false
}

func GetConfigFile() string {
	return config.GetConfigSet()
}
