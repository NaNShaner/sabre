/*
@spec: 中间件的安装部署
*/

package getdeploypkg

import (
	"awesomeProject/pkg/config"
	"awesomeProject/pkg/sabstruct"
	"awesomeProject/pkg/yamlfmt"
	"bufio"
	"fmt"
	"github.com/c4milo/unpackit"
	"io"
	"net/http"
	"os"
	"path"
)

const (
	PkgLocalPathForLinux = "/tmp/" // 默认的安装包以及配置文件的存放地方，可以根据yam文件覆盖
)

//type DeployAction struct {
//	Install   string // 部署
//	ReInstall string // 重装
//	apply     string // 配置修改
//	OffLine   string // 下线
//	Start     string //	启动
//	Stop      string //	停止
//	Restop    string // 重启
//}
//
//type MiddInfo struct {
//	NameSpace  string
//	MidType    string   // 中间件的类型
//	MidVersion string   // 中间件的版本
//	MidRunType []string // 运行模式，集群、主备、冷备等等
//}
//
//type DeployPkg struct {
//	timer        time.Time // 执行时间
//	PkgFromUrl   string    // 安装所需文件连接
//	ConfigFile   string    // 安装所需配置文件
//	MiddInfo               // 中间件相关信息
//	DeployAction           // 执行动作
//}

type Basest sabstruct.Config

// GetDeployPkgFromUrl 从服务端获取安装包或者配置文件等
func (u *Basest) GetDeployPkgFromUrl() (string, error) {
	pkgPath := PkgLocalPathForLinux
	pkgUrl := u.Spec.PKGDownloadPath
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

// UnpackPkg 解压下载的.tar.gz 文件包
func (u *Basest) UnpackPkg(tarFileAbsPath string) error {
	tarFileName := path.Base(tarFileAbsPath)

	// 解压目录下如已有同名文件，则报错
	if !IsFileExist(path.Join(PkgLocalPathForLinux, tarFileName)) {
		return fmt.Errorf("failed to Unpack tar file: the file is already exist")
	}
	tarFile, err := os.Open(tarFileAbsPath)
	defer tarFile.Close()
	if err != nil {
		return fmt.Errorf("failed to open tar file: %v", err)
	}
	_, err = unpackit.Unpack(tarFile, PkgLocalPathForLinux)
	if err != nil {
		return fmt.Errorf("failed to Unpack tar file: %v", err)
	}
	return nil

}

// SetConfigFile 修改配置文件
// m: 顶层的结构体
// f: 文件文件的绝对路径
func (u *Basest) SetConfigFile(m sabstruct.Config, f string) {
	defalutConfig := config.GetConfigSet()
	var Config sabstruct.Config
	yamlFmt, err := yamlfmt.YamlFmt(defalutConfig, Config)
	if err != nil {
		return
	}

	if m.Spec.Midtype == "Tomcat" {
		m.Spec.Tomcat.Javaopts = yamlFmt.Spec.Jdk.Javaopts
	}

}
