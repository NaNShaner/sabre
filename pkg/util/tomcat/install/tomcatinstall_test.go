package tomcatinstall

import (
	"fmt"
	"sabre/pkg/apiserver"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestTomcatInstall(t *testing.T) {
	//var m commontools.Basest
	f := "/Users/bijingrui/sabre/pkg/deploy/tomcatInstll.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	printResultJson, err := yamlfmt.PrintResultJson((*commontools.Basest)(yamlFmt))
	if err != nil {
		return
	}
	fmt.Printf("%s\n", printResultJson)
	if err != nil {
		return
	}
	//install, err := Deploy((*commontools.Basest)(yamlFmt))
	//if err != nil {
	//	t.Errorf("install fail %s", err)
	//}
	//t.Log(install)
	setInfoToDB, setInfoToDBErr := apiserver.HttpReq((*apiserver.Basest)(yamlFmt))
	if setInfoToDBErr != nil {
		return
	}
	t.Log(setInfoToDB)

	// 获取Tomcat安装目录
	//InstallHomePath, getInstallHomePatherr := ioutil.ReadDir("/tmp")
	//if getInstallHomePatherr != nil {
	//	t.Error(getInstallHomePatherr)
	//}
	//// t.Logf("%q", InstallHomePath)
	//for _, file := range InstallHomePath {
	//	fmt.Println(file.Name())
	//}
}
