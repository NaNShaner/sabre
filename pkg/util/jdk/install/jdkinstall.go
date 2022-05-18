package install

import (
	"fmt"
	"os"
	"sabre/pkg/config"
	"sabre/pkg/util/aboutuser"
	cf "sabre/pkg/util/changefile"
	"sabre/pkg/util/commontools"
	"strings"
)

type JdkInstall struct{}

//JdkInstallForApp 用户及用户组(判断||新建)-下载安装包-解压-修改配置文件-安装校验(尝试启动，并进行健康检查，通过后关闭)
// 除去jdk的安装包， 还有应用使用的文件
//TODO：信息上送网关并入库
func (j *JdkInstall) JdkInstallForApp(m *commontools.Basest) error {

	// 如果用户未输入Tomcat的jvm参数，设置默认参数，默认参数来自~/.sabrefig/config
	defaultCf, defaultCfErr := config.GetConfigSet()
	if defaultCfErr != nil {
		return fmt.Errorf("获取服务器默认配置失败,%s\n", defaultCfErr)
	}
	// 这里判断用户是否输入了jvm参数
	// TODO 后续在validator中添加校验，该字段必填
	if m.Spec.DefaultConfig.Jdk.Javaopts == "" {
		m.Spec.DefaultConfig.Jdk.Javaopts = defaultCf.Spec.Jdk.Javaopts
	}

	// 安装JDK
	onlyJdkInstallPath, onlyJdkInstallErr := j.OnlyJdkInstall(m)
	if onlyJdkInstallErr != nil {
		return onlyJdkInstallErr
	}

	// 配置用户的环境变量
	SetJdkEnvErr := j.SetJdkEnv(m, onlyJdkInstallPath)
	if SetJdkEnvErr != nil {
		return SetJdkEnvErr
	}

	// 创建jar目录以及部署启动脚本
	setJdkForAppErr := j.SetJdkForApp(m)
	if setJdkForAppErr != nil {
		return setJdkForAppErr
	}
	return nil
}

//OnlyJdkInstall 只安装JDK，不做任何其他动作
// return 安装路径
func (j *JdkInstall) OnlyJdkInstall(m *commontools.Basest) (string, error) {
	installPath, err := m.InstallCommonStep()
	if err != nil {
		return "", fmt.Errorf("JdkInstall 步骤执行失败%s\n", err)
	}
	return installPath, nil
}

//SetJdkEnv 修改用户的JDK环境变量 修改 ~/.bashrc 文件
//jdkPath 为jdk的安装目录
func (j *JdkInstall) SetJdkEnv(m *commontools.Basest, jdkPath string) error {
	jdkVersion := strings.Replace("export JAVA_HOME=jdkVersion", "jdkVersion", jdkPath, -1)
	setJdkEnv := `export PATH=$JAVA_HOME/bin:$PATH`

	// 获取用户家目录
	userHomeDir, userHomeDirErr := aboutuser.GetUserHomeDir(m.Spec.User.Name)
	if userHomeDirErr != nil {
		return fmt.Errorf("用户家目录不存在%s", userHomeDir)
	}
	// 获取用户的.bashrc 文件路径，判断是否存在
	bashrcFile := userHomeDir + "/.bashrc"
	userBashrcFileExist := commontools.IsFileExist(bashrcFile)
	// 如果存在
	if !userBashrcFileExist {
		s := []string{jdkVersion, setJdkEnv}
		// 在文件末尾追加环境变量的配置文件
		// TODO 再追加之前需要先判断是否已经有JAVA_HOME参数，避免覆盖
		err := cf.AppendFile(bashrcFile, s)
		if err != nil {
			return fmt.Errorf("追加环境变量到文件%s，失败，报错为%s\n", bashrcFile, err)
		}
	} else {
		return fmt.Errorf("文件%s不存在\n", bashrcFile)
	}
	return nil
}

//SetJdkForApp 部署jar包的目录结构以及启动脚本
func (j *JdkInstall) SetJdkForApp(m *commontools.Basest) error {
	if err := os.MkdirAll(m.Spec.DefaultConfig.Jdk.JdkAppInstallPath, 0740); err != nil {
		return err
	}
	pkgFromUrl, err := m.GetDeployPkgFromUrl(m.Spec.DefaultConfig.Jdk.JdkStartUpFile)
	if err != nil {
		return err
	}
	moveFileErr := os.Rename(pkgFromUrl, m.Spec.DefaultConfig.Jdk.JdkAppInstallPath)
	if moveFileErr != nil {
		return moveFileErr
	}
	return nil
}
