package tomcatinstall

import (
	"fmt"
	"io/ioutil"
	"path"
	"sabre/pkg/config"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/changefile"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"time"
)

//TomcatInstall 用户及用户组(判断||新建)-下载安装包-解压-修改配置文件-安装校验(尝试启动，并进行健康检查，通过后关闭)
//TODO：信息上送网关并入库
//TODO：ajp port 暂不支持修改
//TODO：shutdown port 暂不支持修改
func TomcatInstall(m *commontools.Basest) (bool, error) {
	if m.DeployAction.Action != "Install" {
		return false, fmt.Errorf("yaml文件为声明Tomcat的安装行为\n")
	}

	if err := m.InstallCommonStep(); err != nil {
		return false, fmt.Errorf("TomcatInstall 步骤执行失败，%s", err)
	}

	// 如果用户未输入Tomcat的jvm参数，设置默认参数，默认参数来自/root/.sabrefig/config
	defaultCf, defaultCfErr := yamlfmt.YamlFmt(config.GetConfigSet(), sabstruct.Config{})
	if defaultCfErr != nil {
		return false, fmt.Errorf("获取服务器默认配置失败,%s", defaultCfErr)
	}
	// 这里判断用户是否输入了jvm参数
	// TODO 后续在validator中添加校验，该字段必填
	if m.Spec.DefaultConfig.Jdk.Javaopts == "" {
		m.Spec.DefaultConfig.Jdk.Javaopts = defaultCf.Spec.Jdk.Javaopts
	}

	// 获取Tomcat家目录
	tomcatHomePath, err := GetTomcatHomePath(m)
	if err != nil {
		return false, err
	}

	//修改catalina.sh
	ChangeCatalinaShErr := ChangeCatalinaSh(m, tomcatHomePath, "/bin/catalina.sh")
	if ChangeCatalinaShErr != nil {
		return false, ChangeCatalinaShErr
	}

	// 修改server.xml
	ChangeServerXmlErr := ChangeServerXml(m, tomcatHomePath, "/conf/server.xml")
	if ChangeServerXmlErr != nil {
		return false, ChangeServerXmlErr
	}

	// 启动Tomcat
	startUp := path.Join(m.Spec.InstallPath + "/apache-tomcat-7.0.75/bin/startup.sh")
	startMiddleware, err := m.ExecCmdWithTimeOut(startUp, time.Duration(3))
	if err != nil {
		return false, err
	}
	fmt.Printf("命令执行情况%s", startMiddleware)

	return true, nil
}

// GetTomcatHomePath 获取Tomcat安装目录
func GetTomcatHomePath(m *commontools.Basest) (string, error) {
	InstallHomePath, getInstallHomePatherr := ioutil.ReadDir(m.Spec.InstallPath)
	if getInstallHomePatherr != nil {
		return "", getInstallHomePatherr
	}
	if len(InstallHomePath) != 1 {
		return "", fmt.Errorf("%s 目录下含多层目录，无法执行操作，请检查", m.Spec.InstallPath)
	}
	for _, p := range InstallHomePath {
		if p.IsDir() {
			return p.Name(), nil
		} else {
			return "", fmt.Errorf("%s 目录下无目录文件", m.Spec.InstallPath)
		}
	}
	return "", fmt.Errorf("无法获取Tomcat的家目录，请检查")
}

//ChangeCatalinaSh  修改catalina.sh配置文件添加jvm参数
// p Tomcat的家目录，file 需要改的具体文件
func ChangeCatalinaSh(m *commontools.Basest, p, file string) error {
	catalinaReplace := make(map[string]string)
	catalinaReplace["JAVAOPTS"] = m.Spec.DefaultConfig.Jdk.Javaopts
	catalina := p + file
	catalinaReplaceErr := changefile.Changefile(catalina, catalinaReplace)
	if catalinaReplaceErr != nil {
		return catalinaReplaceErr
	}
	return nil
}

//ChangeServerXml  修改server.xml
func ChangeServerXml(m *commontools.Basest, p, file string) error {
	serverXmlReplace := make(map[string]string)
	serverXmlReplace["listeningport"] = m.Spec.DefaultConfig.Tomcat.ShutdownPort
	serverXml := p + file
	serverXmlReplaceErr := changefile.Changefile(serverXml, serverXmlReplace)
	if serverXmlReplaceErr != nil {
		return fmt.Errorf("修改配置文件%s失败,%s", serverXml, serverXmlReplaceErr)
	}
	return nil
}