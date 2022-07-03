package tomcatinstall

import (
	"fmt"
	"io/ioutil"
	"path"
	"sabre/pkg/config"
	"sabre/pkg/util/changefile"
	"sabre/pkg/util/commontools"
	"sabre/pkg/util/logbase/logscheduled"
	"time"
)

//Deploy 用户及用户组(判断||新建)-下载安装包-解压-修改配置文件-安装校验(尝试启动，并进行健康检查，通过后关闭)
func Deploy(m *commontools.Basest) (string, error) {
	if m.DeployAction.Action != "Install" {
		return "", fmt.Errorf("yaml文件为声明Tomcat的安装行为%s\n", m.DeployAction.Action)
	}
	localServerName, getLocalServerNameErr := commontools.GetLocalServerName()
	if getLocalServerNameErr != nil {
		return "", getLocalServerNameErr
	}
	unPackPath, err := m.InstallCommonStep()
	if err != nil {
		return "", fmt.Errorf("server %s deploy 步骤执行失败，%s", localServerName, err)
	}

	// 如果用户未输入Tomcat的jvm参数，设置默认参数，默认参数来自/root/.sabrefig/config
	defaultCf, defaultCfErr := config.GetConfigSet()
	if defaultCfErr != nil {
		return "", fmt.Errorf("failed to get the default configuration of server %s,%s", localServerName, defaultCfErr)
	}
	// 这里判断用户是否输入了jvm参数，如果用户没有输入，可以通过默认配置填充
	if m.Spec.DefaultConfig.Jdk.Javaopts == "" {
		m.Spec.DefaultConfig.Jdk.Javaopts = defaultCf.Spec.Jdk.Javaopts
	}

	// 获取Tomcat家目录
	_, err = GetTomcatHomePath(m)
	if err != nil {
		return "", err
	}

	//修改catalina.sh
	ChangeCatalinaShErr := ChangeCatalinaSh(m, unPackPath, "/bin/catalina.sh")
	if ChangeCatalinaShErr != nil {
		return "", ChangeCatalinaShErr
	}

	// 修改server.xml
	ChangeServerXmlErr := ChangeServerXml(m, unPackPath, "/conf/server.xml")
	if ChangeServerXmlErr != nil {
		return "", ChangeServerXmlErr
	}

	// 启动Tomcat
	startUp := path.Join(unPackPath + "/bin/startup.sh")
	_, err = m.ExecCmdWithTimeOut(startUp, time.Duration(3))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s delopy done\n", m.Midtype), nil

}

//func Upload(m *commontools.Basest) (string, error) {
//	// 信息入库
//	setInfoToDB, setInfoToDBErr := m.SetInfoToDB()
//	if setInfoToDBErr != nil {
//		return "", setInfoToDBErr
//	}
//	fmt.Printf("入库成功 ===> %s\n", setInfoToDB)
//	return setInfoToDB, nil
//}

// GetTomcatHomePath 获取Tomcat安装目录
func GetTomcatHomePath(m *commontools.Basest) (string, error) {
	InstallHomePath, getInstallHomePathErr := ioutil.ReadDir(m.Spec.InstallPath)

	if getInstallHomePathErr != nil {
		return "", getInstallHomePathErr
	}
	if len(InstallHomePath) != 1 {
		return "", fmt.Errorf("%s 目录下含多层目录，无法执行操作，请检查\n", m.Spec.InstallPath)
	}
	for _, p := range InstallHomePath {
		if p.IsDir() {
			logscheduled.Log.Info(p.Name())
			fmt.Printf("GetTomcatHomePath 目录文件%s\n", p.Name())
			return p.Name(), nil
		} else {
			return "", fmt.Errorf("%s 目录下无目录文件\n", m.Spec.InstallPath)
		}
	}
	return "", fmt.Errorf("无法获取Tomcat的家目录，请检查\n")
}

//ChangeCatalinaSh  修改catalina.sh配置文件添加jvm参数
// p Tomcat的家目录，file 需要改的具体文件
func ChangeCatalinaSh(m *commontools.Basest, p, file string) error {
	catalinaReplace := make(map[string]string)
	catalinaReplace["JAVAOPTS"] = m.Spec.DefaultConfig.Jdk.Javaopts
	catalina := p + file
	catalinaReplaceErr := changefile.Changefile(catalina, catalinaReplace)
	if catalinaReplaceErr != nil {
		return fmt.Errorf("修改配置文件%s失败,%s\n", catalina, catalinaReplaceErr)
	}
	return nil
}

//ChangeServerXml  修改server.xml
func ChangeServerXml(m *commontools.Basest, p, file string) error {
	serverXmlReplace := make(map[string]string)
	serverXmlReplace["listeningport"] = m.Spec.DefaultConfig.Tomcat.ListeningPort
	serverXmlReplace["shutdownport"] = m.Spec.DefaultConfig.Tomcat.ShutdownPort
	serverXmlReplace["ajpport"] = m.Spec.DefaultConfig.Tomcat.AjpPort
	serverXml := p + file
	serverXmlReplaceErr := changefile.Changefile(serverXml, serverXmlReplace)
	if serverXmlReplaceErr != nil {
		return fmt.Errorf("修改配置文件%s失败,%s\n", serverXml, serverXmlReplaceErr)
	}
	return nil
}
