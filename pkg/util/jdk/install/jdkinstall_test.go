package install

import (
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestJdkInstall_SetJdkEnv(t *testing.T) {
	f := "/Users/bijingrui/sabre/pkg/util/jdk/install/onlyJdkInstall.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	if err != nil {
		t.Error(err)
	}
	m := (*commontools.Basest)(yamlFmt)
	var s JdkInstall
	setJdkEnvErr := s.SetJdkEnv(m, m.Spec.InstallPath)
	if setJdkEnvErr != nil {
		t.Error(setJdkEnvErr)
	}
	t.Log("Done\n")

}

func TestJdkInstall_OnlyJdkInstall(t *testing.T) {
	var c commontools.Basest
	c.InstallPath = "/tmp"
	_, err := c.UnpackPkg("/Users/bijingrui/sabre/pkg/util/jdk/install/jdk-7u71-linux-x64.tar.gz")
	if err != nil {
		return
	}
}
