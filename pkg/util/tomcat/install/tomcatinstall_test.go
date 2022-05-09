package tomcatinstall

import (
	"fmt"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestTomcatInstall(t *testing.T) {
	// var m commontools.Basest
	f := "/Users/bijingrui/sabre/pkg/getdeploypkg/tomcatInstll.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	printResultJson, err := yamlfmt.PrintResultJson((*commontools.Basest)(yamlFmt))
	if err != nil {
		return
	}
	fmt.Printf("%s\n", printResultJson)
	if err != nil {
		return
	}
	install, err := TomcatInstall((*commontools.Basest)(yamlFmt))
	if err != nil {
		t.Errorf("install fail %s", err)
	}
	t.Log(install)
}
