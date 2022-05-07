package tomcatinstall

import (
	"awesomeProject/pkg/sabstruct"
	"awesomeProject/pkg/util/commontools"
	"awesomeProject/pkg/yamlfmt"
	"testing"
)

func TestTomcatInstall(t *testing.T) {
	// var m commontools.Basest
	f := "/Users/bijingrui/awesomeProject/pkg/getdeploypkg/tomcatInstll.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	if err != nil {
		return
	}
	install, err := TomcatInstall((*commontools.Basest)(yamlFmt))
	if err != nil {
		t.Errorf("install fail %s", err)
	}
	t.Log(install)
}
