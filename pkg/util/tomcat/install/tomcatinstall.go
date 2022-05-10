package tomcatinstall

import (
	"fmt"
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

	// 修改配置文件
	// 修改catalina.sh
	catalinaReplace := make(map[string]string)
	catalinaReplace["JAVAOPTS"] = m.Spec.DefaultConfig.Jdk.Javaopts
	catalina := path.Join(m.Spec.InstallPath + "/apache-tomcat-7.0.75/bin/catalina.sh")
	//replaceCmd := strings.Replace(`"sed -i "s/JAVAOPTS/abc/g" "`,"abc",catalinaReplace["JAVAOPTS"],-1) + catalina
	//_, err := m.StartMiddleware(replaceCmd, time.Duration(3))
	catalinaReplaceErr := changefile.Changefile(catalina, catalinaReplace)
	if catalinaReplaceErr != nil {
		return false, catalinaReplaceErr
	}
	// 修改server.xml
	serverXmlReplace := make(map[string]string)
	serverXmlReplace["listeningport"] = m.Spec.DefaultConfig.Tomcat.ShutdownPort
	serverXml := path.Join(m.Spec.InstallPath + "/apache-tomcat-7.0.75/conf/server.xml")
	serverXmlReplaceErr := changefile.Changefile(serverXml, serverXmlReplace)
	if serverXmlReplaceErr != nil {
		return false, fmt.Errorf("修改配置文件%s失败,%s", serverXml, serverXmlReplaceErr)
	}
	// 启动Tomcat
	startUp := path.Join(m.Spec.InstallPath + "/apache-tomcat-7.0.75/bin/startup.sh")
	startMiddleware, err := m.StartMiddleware(startUp, time.Duration(3))
	if err != nil {
		return false, err
	}
	fmt.Printf("命令执行情况%s", startMiddleware)

	return true, nil
}
