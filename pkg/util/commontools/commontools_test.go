package commontools

import (
	"testing"
	"time"
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
	pkgFromUrl, err := d.GetDeployPkgFromUrl(url)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pkgFromUrl)
}

func TestBasest_UnpackPkg(t *testing.T) {
	var d Basest
	tarFile := "/Users/bijingrui/awesomeProject/pkg/deploy/apache-tomcat-7.0.75.tar.gz"
	d.Spec.InstallPath = "/Users/bijingrui/awesomeProject/pkg/deploy/"
	_, err := (*Basest).UnpackPkg(&d, tarFile)
	if err != nil {

		t.Errorf("==> %s", err)
		return
	}
	t.Log("done")
}

func TestBasest_StartMiddleware(t *testing.T) {
	var d *Basest
	var timer time.Duration
	startscript := "sleep 2; echo hello"
	startMiddlewareReslut, err := d.ExecCmdWithTimeOut(startscript, timer*3)
	if err != nil {
		t.Error(err)
	}
	t.Log(startMiddlewareReslut)
}

func TestBasest_CheckInstallServerBelongToNS(t *testing.T) {
	var d Basest

	d.DeployHost = []string{"192.168.3.152"}
	d.Netarea = "app"
	d.Namespace = "erp"
	err := d.CheckInstallServerBelongToNS()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("passed")
}
