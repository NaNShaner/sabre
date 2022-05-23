package sabstruct

import (
	"time"
)

// Config 平台配置
// Kind 资源类型 中间件、网络、基础设施
//	Deployment:资源类型	config:配置文件
type Config struct {
	// validate 必填字段
	ApiVersion string `json:"apiVersion" validate:"required"`
	Kind       string `json:"kind" validate:"required"`
	Server     string `json:"server"`
	ApiServer  string `json:"apiserver"`
	Metadata   `json:"metadata"`
	Spec       `json:"spec"`
}

// Metadata 存放中间件所属应用的信息
type Metadata struct {
	Namespace string `json:"namespace"`
	Netarea   string `json:"netarea"`
	Appname   string `json:"appname,omitempty"`
}

// Spec 存放中间件本身版本及其配置参数信息
// TODO: MidRunType 标识运行时的状态，逻辑待补充
type Spec struct {
	Midtype         string   `json:"midtype"`
	Version         string   `json:"version"`
	Port            string   `json:"port,omitempty"`
	InstallPath     string   `json:"installPath"`
	PKGDownloadPath string   `json:"pkgDownloadPath"`
	MidRunType      []string `json:"run_type,omitempty"` // 运行模式，集群、主备、冷备等等
	User            `json:"user"`
	DefaultConfig   `json:"default,omitempty"`
	DeployAction    `json:"deployaction"`
}

// User 中间件所属用户信息
type User struct {
	Name  string `json:"name"`  // 中间件所属用户
	Group string `json:"group"` // 中间件所属用户组
}

// DeployAction 执行动作
// Action，针对Tomcat 包含Install
// Action，针对Jdk 包含Install，appInstall。含义是Install仅安装jdk并配置环境变量、 appInstall表示安装jdk、配置变量并且生成启动jar包的文件目录以及启动脚本
type DeployAction struct {
	Timer      time.Time `json:"timer,omitempty"` // 执行时间
	Action     string    `json:"action"`
	DeployHost []string  `json:"deploy_host"`
	//Install   string    `json:"install"`         // 部署
	//ReInstall string    `json:"re_install"`      // 重装
	//Apply     string    `json:"apply"`           // 配置修改
	//OffLine   string    `json:"off_line"`        // 下线
	//Start     string    `json:"start"`           //	启动
	//Stop      string    `json:"stop"`            //	停止
	//Restop    string    `json:"restop"`          // 重启
}

// Jdk 在~/.sabrefig/config 文件的默认配置
type Jdk struct {
	Javaopts          string `json:"javaopts"`
	JdkAppInstallPath string `json:"appinstallpath"`
	JdkStartUpFile    string `json:"startup"`
}

// Tomcat 在.sabrefig/config 文件的默认配置
type Tomcat struct {
	Javaopts      string `json:"javaopts"`
	ListeningPort string `json:"listeningport"`
	AjpPort       string `json:"ajpport"`
	AjpRirectPort string `json:"ajprirectport"`
	ShutdownPort  string `json:"shutdownport"`
}

// DefaultConfig 各类资源在.sabrefig/config 文件的默认配置
type DefaultConfig struct {
	Jdk    `json:"jdk"`
	Tomcat `json:"tomcat"`
}

//IsZero 校验必填字段
func (n *Config) IsZero() bool {
	return n.Kind == "" && n.ApiVersion == "" && n.Metadata == Metadata{}
}
