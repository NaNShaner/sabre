package tomcatinstall

import (
	"awesomeProject/pkg/config"
	"awesomeProject/pkg/sabstruct"
	"awesomeProject/pkg/util/changefile"
	"awesomeProject/pkg/util/commontools"
	"awesomeProject/pkg/yamlfmt"
	"fmt"
	"path"
	"time"
)

//TomcatInstall 用户及用户组(判断||新建)-下载安装包-解压-修改配置文件-安装校验(尝试启动，并进行健康检查，通过后关闭)
//TODO：信息上送网关并入库
func TomcatInstall(m *commontools.Basest) (bool, error) {

	if err := m.InstallCommonStep(); err != nil {
		return false, fmt.Errorf("TomcatInstall 步骤执行失败%s", err)
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

	// 修改配置文件
	// 修改catalina.sh
	catalinaReplace := make(map[string]string)
	catalinaReplace["JAVAOPTS"] = m.Spec.DefaultConfig.Jdk.Javaopts
	catalina := path.Join(m.Spec.InstallPath + "bin/catalina.sh")
	catalinaReplaceErr := changefile.Changefile(catalina, catalinaReplace)
	if catalinaReplaceErr != nil {
		return false, fmt.Errorf("修改配置文件%s失败,%s", catalina, catalinaReplaceErr)
	}
	// 修改server.xml
	serverXmlReplace := make(map[string]string)
	serverXmlReplace["shutdownport"] = m.Spec.DefaultConfig.Jdk.Javaopts
	serverXmlReplace["listeningport"] = m.Spec.DefaultConfig.Jdk.Javaopts
	serverXmlReplace["ajpport"] = m.Spec.DefaultConfig.Jdk.Javaopts
	serverXmlReplace["ajprirectport"] = m.Spec.DefaultConfig.Jdk.Javaopts
	serverXml := path.Join(m.Spec.InstallPath + "conf/server.xml")
	serverXmlReplaceErr := changefile.Changefile(serverXml, serverXmlReplace)
	if serverXmlReplaceErr != nil {
		return false, fmt.Errorf("修改配置文件%s失败,%s", serverXml, serverXmlReplaceErr)
	}
	// 启动Tomcat
	startUp := path.Join(m.Spec.InstallPath + "bin/start.sh")
	startMiddleware, err := m.StartMiddleware(startUp, time.Duration(3))
	if err != nil {
		return false, err
	}
	fmt.Printf("命令执行情况%s", startMiddleware)

	return true, nil
}
