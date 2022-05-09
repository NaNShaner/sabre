/*
@spec: 中间件的安装部署
*/
//
package getdeploypkg

import (
	"fmt"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	Ti "sabre/pkg/util/tomcat/install"
	"sabre/pkg/yamlfmt"
)

//import (
//	"awesomeProject/pkg/config"
//	"awesomeProject/pkg/sabstruct"
//	"awesomeProject/pkg/yamlfmt"
//)

//const (
//	PkgLocalPathForLinux = "/tmp/" // 默认的安装包以及配置文件的存放地方，可以根据yam文件覆盖
//)

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

// SetConfigFile 修改配置文件
// m: 顶层的结构体，来自用户输入的yaml解析之后的数据
// f: 文件文件的绝对路径
//func (u *Basest) SetConfigFile(m sabstruct.Config, f string) {
//	// 默认配置默认配置来至 /root/.sabrefig/config
//	defalutConfig := config.GetConfigSet()
//	var Config sabstruct.Config
//	// 获取用户输入的yaml文件
//	yamlFmt, err := yamlfmt.YamlFmt(defalutConfig, Config)
//	if err != nil {
//		return
//	}

// 如果是Tomcat，判断输入是否声明了java_opts参数，如果没有，使用默认配置
//if m.Spec.Midtype == "Tomcat" {
//	m.Spec.Tomcat.Javaopts = yamlFmt.Spec.Jdk.Javaopts
//	switch expr {
//
//	}
//}
//
//}

//Deploy 中间件的部署
func Deploy() (bool, error) {
	f := "/Users/bijingrui/awesomeProject/pkg/getdeploypkg/tomcatInstll.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	if err != nil {
		return false, fmt.Errorf("解析文件%s失败,%s", f, err)
	}
	install, err := Ti.TomcatInstall((*commontools.Basest)(yamlFmt))
	if err != nil {
		return false, fmt.Errorf("install fail %s", err)
	}

	return install, nil

}
