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
	url := "http://124.71.219.53:8001/uploads/uploads/2022/04/30/apache-tomcat-7.0.75.tar.gz"
	var d Basest
	d.Spec.PKGDownloadPath = url
	pkgFromUrl, err := (*Basest).GetDeployPkgFromUrl(&d)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pkgFromUrl)
}

func TestBasest_UnpackPkg(t *testing.T) {
	var d Basest
	tarFile := "/Users/bijingrui/awesomeProject/pkg/getdeploypkg/apache-tomcat-7.0.75.tar.gz"
	d.Spec.InstallPath = "/Users/bijingrui/awesomeProject/pkg/getdeploypkg/"
	err := (*Basest).UnpackPkg(&d, tarFile)
	if err != nil {

		t.Errorf("==> %s", err)
		return
	}
	t.Log("done")
}
