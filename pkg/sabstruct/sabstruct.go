package sabstruct

// Config 平台配置
// Kind 资源类型 中间件、网络、基础设施
//	Deployment:资源类型	config:配置文件
type Config struct {
	// validate 必填字段
	ApiVersion     string `json:"apiVersion" validate:"required"`
	Kind           string `json:"kind" validate:"required"`
	Server         string `json:"server,omitempty"`
	ApiServer      string `json:"apiserver,omitempty"`
	LocalLetServer string `json:"letserver,omitempty"`
	Metadata       `json:"metadata"`
	Spec           `json:"spec"`
}

// Metadata 存放中间件所属应用的信息
type Metadata struct {
	//Namespace 应用系统简称
	Namespace string `json:"namespace"`
	//Netarea Appname 所属网络安全域
	Netarea string `json:"netarea"`
	//Appname 应用名称
	Appname string `json:"appname,omitempty"`
}

// Spec 存放中间件本身版本及其配置参数信息
// TODO: MidRunType 标识运行时的状态，逻辑待补充
type Spec struct {
	Midtype         string   `json:"midtype"`
	Version         string   `json:"version"`
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
	Timer            string                     `json:"timer,omitempty"` // 执行时间
	Action           string                     `json:"action"`
	DeployHost       []string                   `json:"deploy_host"`
	DeployHostStatus []map[string]RunTimeStatus `json:"deployHostStatus,omitempty"`
	//Install   string    `json:"install"`         // 部署
	//ReInstall string    `json:"re_install"`      // 重装
	//Apply     string    `json:"apply"`           // 配置修改
	//OffLine   string    `json:"off_line"`        // 下线
	//Start     string    `json:"start"`           // 启动
	//Stop      string    `json:"stop"`            // 停止
	//Restop    string    `json:"restop"`          // 重启
}

//RunTimeStatus 描述中间件运行态。
//TODO 逻辑待补充
type RunTimeStatus struct {
	//RunStatus 中间件资源当前是否运行正常，由sabrelet定时上报
	RunStatus bool `json:"run_status,omitempty"`
	//StatusReportTimer 记录上报状态的时间
	StatusReportTimer string `json:"status_report_timer,omitempty"`
	//RunningDays 记录上报状态的时间
	RunningDays int `json:"running_days,omitempty"`
}

// Jdk 在~/.sabrefig/config 文件的默认配置
type Jdk struct {
	Javaopts          string `json:"javaopts,omitempty"`
	JdkAppInstallPath string `json:"appinstallpath,omitempty"`
	JdkStartUpFile    string `json:"startup,omitempty"`
}

// Tomcat 在.sabrefig/config 文件的默认配置
type Tomcat struct {
	Javaopts      string `json:"javaopts,omitempty"`
	ListeningPort string `json:"listeningport,omitempty"`
	AjpPort       string `json:"ajpport,omitempty"`
	//AjpRirectPort string `json:"ajprirectport,omitempty"`
	ShutdownPort string `json:"shutdownport,omitempty"`
}

// DefaultConfig 各类资源在.sabrefig/config 文件的默认配置
type DefaultConfig struct {
	Jdk    `json:"jdk,omitempty"`
	Tomcat `json:"tomcat,omitempty"`
}
