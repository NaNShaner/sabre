package tomcatinstall

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestTomcatInstall(t *testing.T) {
	// var m commontools.Basest
	//f := "/Users/bijingrui/sabre/pkg/getdeploypkg/tomcatInstll.yaml"
	//yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	//printResultJson, err := yamlfmt.PrintResultJson((*commontools.Basest)(yamlFmt))
	//if err != nil {
	//	return
	//}
	//fmt.Printf("%s\n", printResultJson)
	//if err != nil {
	//	return
	//}
	//install, err := TomcatInstall((*commontools.Basest)(yamlFmt))
	//if err != nil {
	//	t.Errorf("install fail %s", err)
	//}
	//t.Log(install)

	// 获取Tomcat安装目录
	InstallHomePath, getInstallHomePatherr := ioutil.ReadDir("/tmp")
	if getInstallHomePatherr != nil {
		t.Error(getInstallHomePatherr)
	}
	// t.Logf("%q", InstallHomePath)
	for _, file := range InstallHomePath {
		fmt.Println(file.Name())
	}
}
