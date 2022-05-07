package commontools

import (
	"awesomeProject/pkg/sabstruct"
	"awesomeProject/pkg/util/aboutuser"
	"bufio"
	"context"
	"fmt"
	"github.com/c4milo/unpackit"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	// PkgLocalPathForLinux 默认的安装包以及配置文件的存放地方，可以根据yam文件覆盖
	PkgLocalPathForLinux = "/tmp/"
)

// Basest 分支结构，继承顶层结构 sabstruct.Config
type Basest sabstruct.Config

// GetDeployPkgFromUrl 从服务端获取安装包或者配置文件等
// 返回下载文件的绝对路径
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

// UnpackPkg 解压下载的.tar.gz 文件包
// tarFileAbsPath参数为tar包的在服务器上的绝对路径
func (u *Basest) UnpackPkg(tarFileAbsPath string) error {
	// tar包的文件名称
	tarFileName := path.Base(tarFileAbsPath)

	// 解压目录下如已有同名文件，则报错
	if !IsFileExist(path.Join(u.Spec.InstallPath, tarFileName)) {
		return fmt.Errorf("failed to Unpack tar file: the file is already exist")
	}
	tarFile, err := os.Open(tarFileAbsPath)
	defer tarFile.Close()
	if err != nil {
		return fmt.Errorf("failed to open tar file: %v", err)
	}
	_, err = unpackit.Unpack(tarFile, u.Spec.InstallPath)
	if err != nil {
		return fmt.Errorf("failed to Unpack tar file: %v", err)
	}
	return nil

}

// StartMiddleware 中间件的启动
// startscript 为中间件的启动脚本的绝对路径
// timer 为命令执行的超时时间
// TODO 以不同用户执行命令，目前只能使用root
func (u *Basest) StartMiddleware(startscript string, timer time.Duration) (string, error) {
	// 设置上下文超时时间，目前默认设置为3秒
	if timer == 0 {
		timer = 3
	}
	ctx, cancel := context.WithTimeout(context.Background(), timer*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", startscript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s命令执行失败, %s\n", cmd.String(), err)
	}
	return string(output), nil
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

func (u *Basest) InstallCommonStep() error {

	// 判断用户在服务器上是否存在，函数返回为bool值
	isUserExist, userExistErr := aboutuser.IsUserExist(u.User.Name)
	if userExistErr != nil {
		return fmt.Errorf("用户%s不存在，%s", u.User.Name, userExistErr)
	}
	// 如果用户存在
	if isUserExist {
		getPkgFromUrl, getPkgFromUrlErr := u.GetDeployPkgFromUrl()
		if getPkgFromUrlErr != nil {
			return fmt.Errorf("下载%s文件失败，%s", u.Spec.PKGDownloadPath, getPkgFromUrlErr)
		}
		// 解压安装包，解压到 m.Spec.InstallPath 路径下
		unpackPkgErr := u.UnpackPkg(getPkgFromUrl)
		if unpackPkgErr != nil {
			return fmt.Errorf("解压%s文件失败，%s", getPkgFromUrl, unpackPkgErr)
		}
	} else {
		// 如果用户不存在
	}
	return nil
}
