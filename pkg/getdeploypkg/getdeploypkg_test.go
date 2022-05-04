package getdeploypkg

import (
	"testing"
)

func TestIsFileExistIs(t *testing.T) {
	c := IsFileExist("./conf.yam")
	if c {
		t.Log("test pass")
	} else {
		t.Error("test fail")
	}

}


func TestDeployPkg_GetDeployPkgFromUrl(t *testing.T) {
	url := "https://dlcdn.apache.org/tomcat/tomcat-8/v8.5.78/bin/apache-tomcat-8.5.78.tar.gz"
	var d DeployPkg
	d.PkgFromUrl = url
	pkgFromUrl, err := d.GetDeployPkgFromUrl()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pkgFromUrl)
}

//func TestTimerFmt(t *testing.T) {
//	c := TimerFmt()
//	t.Log(c)
//	location, err := time.LoadLocation("Asia/Shanghai")
//	if err != nil {
//		return
//	}
//	time.Local = location
//	t.Log(time.Now().Local())
//}