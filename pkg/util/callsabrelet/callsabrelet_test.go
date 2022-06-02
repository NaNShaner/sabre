package callsabrelet

import (
	"fmt"
	"sabre/pkg/sabstruct"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestBasest_CallSabrelet(t *testing.T) {

	f := "/Users/bijingrui/sabre/pkg/util/tomcat/install/deployTomcat.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
	printResultJson, err := yamlfmt.PrintResultJson(yamlFmt)
	if err != nil {
		return
	}
	basest := (*Basest)(yamlFmt)
	fmt.Printf("%s\n", printResultJson)
	hostList := yamlFmt.DeployHost
	CallFaceOfSabrelet(basest, hostList)
}
