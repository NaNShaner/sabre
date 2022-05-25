package yamlfmt

import (
	"sabre/pkg/sabstruct"
	"testing"
)

func TestYamlFmt(t *testing.T) {
	var Configs sabstruct.Config

	f := "/Users/bijingrui/sabre/pkg/util/createfile/test.yaml"
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
