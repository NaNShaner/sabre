package commontools

import (
	"bufio"
	"context"
	"fmt"
	"github.com/c4milo/unpackit"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/aboutuser"
	"strings"
	"time"
)

const (
	// PkgLocalPathForLinux 默认的安装包以及配置文件的存放地方，可以根据yam文件覆盖
	PkgLocalPathForLinux = "/tmp/"
)

// Basest 分支结构，继承顶层结构 sabstruct.Config
type Basest sabstruct.Config

type InstallComm interface {
	GetDeployPkgFromUrl(pkgUrl string) (string, error)
	UnpackPkg(tarFileAbsPath string) (string, error)
	ExecCmdWithTimeOut(startscript string, timer time.Duration) (string, error)
	InstallCommonStep() (string, error)

	//Download()
	//Unpack()
	//ExecCmd()
	//Install()
}

type Download interface {
	GetDeployPkgFromUrl(pkgUrl string) (string, error)
}

type Unpack interface {
	UnpackPkg(tarFileAbsPath string) (string, error)
}

type ExecCmd interface {
	ExecCmdWithTimeOut(startscript string, timer time.Duration) (string, error)
}

type Install interface {
	InstallCommonStep() (string, error)
	CheckIP
}

type CheckIP interface {
	CheckInstallServerBelongToNS() error
}

// GetDeployPkgFromUrl 从服务端获取安装包或者配置文件等
// 返回下载文件的绝对路径
func (u *Basest) GetDeployPkgFromUrl(pkgUrl string) (string, error) {
	pkgPath := PkgLocalPathForLinux
	fileName := path.Base(pkgUrl)
	pkgBasPath := pkgPath + fileName
	resp, err := http.Get(pkgUrl)
	if err != nil {
		fmt.Printf("A error occurred! %s\n", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(resp.Body, 32*1024)

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
		return pkgBasPath, nil
	}
	return "", fmt.Errorf("%s 文件已存在\n", pkgBasPath)

}

// UnpackPkg 解压下载的.tar.gz 文件包
// tarFileAbsPath参数为tar包的在服务器上的绝对路径，函数返回解压后文件的绝对路径
func (u *Basest) UnpackPkg(tarFileAbsPath string) (string, error) {
	// tar包的文件名称
	tarFileName := path.Base(tarFileAbsPath)

	// 解压目录下如已有同名文件，则报错
	if !IsFileExist(path.Join(u.Spec.InstallPath, tarFileName)) {
		return "", fmt.Errorf("failed to Unpack tar file: the file is already exist\n")
	}
	tarFile, err := os.Open(tarFileAbsPath)
	defer tarFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to open tar file: %v\n", err)
	}
	unPackPkgAbsPath, unPackPkgAbsPathErr := unpackit.Unpack(tarFile, u.Spec.InstallPath)
	if unPackPkgAbsPathErr != nil {
		return "", fmt.Errorf("failed to Unpack tar file: %v\n", err)
	}
	return unPackPkgAbsPath, nil

}

// ExecCmdWithTimeOut 中间件的启动
// startscript 为中间件的启动脚本的绝对路径
// timer 为命令执行的超时时间
// TODO 以不同用户执行命令，目前只能使用root
func (u *Basest) ExecCmdWithTimeOut(startscript string, timer time.Duration) (string, error) {
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

//InstallCommonStep 包含检查用户、下载安装包、解压等公共步骤
func (u *Basest) InstallCommonStep() (string, error) {
	// 判断用户在服务器上是否存在，函数返回为bool值
	isUserExist, userExistErr := aboutuser.IsUserExist(u.User.Name)
	if userExistErr != nil {
		return "", fmt.Errorf("用户%s不存在，%s\n", u.User.Name, userExistErr)
	}
	// 如果用户存在
	if isUserExist {
		getPkgFromUrl, getPkgFromUrlErr := u.GetDeployPkgFromUrl(u.Spec.PKGDownloadPath)
		if getPkgFromUrlErr != nil {
			return "", fmt.Errorf("下载%s文件失败，%s\n", u.Spec.PKGDownloadPath, getPkgFromUrlErr)
		}

		//getPkgFromUrl := "/tmp/jdk-8u171-linux-x64.tar.gz"

		// 解压安装包，解压到 m.Spec.InstallPath 路径下
		unPackPkg, unpackPkgErr := u.UnpackPkg(getPkgFromUrl)
		if unpackPkgErr != nil {
			return "", fmt.Errorf("解压%s文件失败，%s\n", getPkgFromUrl, unpackPkgErr)
		}
		return unPackPkg, nil
	} else {
		// 如果用户不存在
		return "", fmt.Errorf("用户%s不存在\n", u.User.Name)
	}
}

//CheckInstallServerBelongToNS 解析yaml文件中 u.DeployHost 字段中的地址是否属于系统下的地址
func (u *Basest) CheckInstallServerBelongToNS() error {
	keyPrefix := res.ManageResourceTypes()
	resType := "/machine"
	key := path.Join("/", keyPrefix[2], u.Namespace, resType, u.Netarea)
	allAvailableServer, getErr := dbload.GetKeyFromETCD(key, true)
	if getErr != nil {
		return getErr
	}
	//for _, s := range GetKeyFromDB(allAvailableServer) {
	//	fmt.Printf("etcd 中有%s\n", s)
	//}
	for _, s := range u.DeployHost {
		if !InMap(ConvertStrSlice2Map(GetKeyFromDB(allAvailableServer)), s) {
			return fmt.Errorf("host %s does not belong to %s system", s, u.Namespace)
		}
	}
	return nil
}

//IsZero 校验必填字段
func (u *Basest) IsZero() bool {
	return u.Kind == "" && u.ApiVersion == "" && u.Metadata == sabstruct.Metadata{}
}

//AddNowTime 添加默认时间
func AddNowTime() string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	t := time.Now().In(cstSh).Format("2006-01-02 15:04:05.1234")
	return t
}

//SetInfoToDB 请求API网关，信息入库
//func (u *Basest) SetInfoToDB() (string, error) {
//	setInfoToDB, setInfoToDBErr := apiserver.HttpReq((*apiserver.Basest)(u))
//	if setInfoToDBErr != nil {
//		return "", setInfoToDBErr
//	}
//	return setInfoToDB, nil
//}

//GetLocalServerName 获取本机主机名
func GetLocalServerName() (string, error) {
	serverName, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return serverName, nil
}

// IsFileExist 判断文件本地是否已经存在
// 存在为false
// 不存在true
func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return true
	}
	//fmt.Printf("文件 %s 已存在， %v\n", info.Name(), err)
	return false
}

//GetHostIP 对于sabrectl提交注册的主机IP进行校验，确认是当前机器的地址
func GetHostIP(ip net.IP) (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	for _, inter := range interfaces {
		ipaddr, _ := inter.Addrs()
		for _, addr := range ipaddr {
			ipAddrWithNetmesk := strings.Split(addr.String(), "/")
			if len(ipAddrWithNetmesk) > 1 {
				if ipAddrWithNetmesk[0] == ip.String() {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("%s is not the IP address of the current machine, please confirm\n", ip.String())
}

// ConvertStrSlice2Map 将字符串 slice 转为 map[string]struct{}，空结构体不再用内存
func ConvertStrSlice2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// InMap 判断字符串是否在 map 中。
func InMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

func GetKeyFromDB(kv []*mvccpb.KeyValue) []string {
	var s []string
	for _, ev := range kv {
		//fmt.Printf("%s\n", ev.Key)
		getServerList := strings.Split(string(ev.Key), "/")
		s = append(s, getServerList[len(getServerList)-1])
	}
	return s
}

//FmtETCDKey 对入库的key进行格式化，第一个字母大写，其他的小写
func FmtETCDKey(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
