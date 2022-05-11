package sabstruct

import (
	"time"
)

// Config 平台配置
// Kind 资源类型 中间件、网络、基础设施
//	Deployment:资源类型	config:配置文件
type Config struct {
	ApiVersion   string `json:"apiVersion"`
	Kind         string `json:"kind"`
	Server       string `json:"server"`
	Metadata     `json:"metadata"`
	Spec         `json:"spec"`
	DeployAction `json:"deployaction"`
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
}

// User 中间件所属用户信息
type User struct {
	Name  string `json:"name"`  // 中间件所属用户
	Group string `json:"group"` // 中间件所属用户组
}

// DeployAction 执行动作
type DeployAction struct {
	Timer  time.Time `json:"timer,omitempty"` // 执行时间
	Action string    `json:"action"`
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
	Javaopts string `json:"javaopts"`
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
