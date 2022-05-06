package yamlfmt

import (
	"awesomeProject/pkg/sabstruct"
	"testing"
)

func TestYamlFmt(t *testing.T) {
	var Configs sabstruct.Config

	f := "/Users/bijingrui/awesomeProject/pkg/yamlfmt/sabrelet.yaml"
	yamlFmt, err := YamlFmt(f, Configs)
	if err != nil {
		return
	}
	resultJson, err := PrintResultJson(*yamlFmt)
	if err != nil {
		return
	}
	t.Logf("%s", resultJson)
}
