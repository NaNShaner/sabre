package dbload

import (
	"sabre/pkg/sabstruct"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestSetIntoDB(t *testing.T) {
	k := "/mid/MNPP"
	f := "/Users/bijingrui/sabre/pkg/yamlfmt/sabrelet.yaml"
	var Configs sabstruct.Config
	yamlFmt, err := yamlfmt.YamlFmt(f, Configs)
	if err != nil {
		t.Error(err)
	}

	resultJson, err := yamlfmt.PrintResultJson(yamlFmt)
	if err != nil {
		return
	}

	SetIntoDBErr := SetIntoDB(k, string(resultJson))

	if SetIntoDBErr != nil {
		t.Errorf("==> %s", SetIntoDBErr)
		return
	}
	t.Log("done")
}

func TestGetKeyWithPrefix(t *testing.T) {
	k := "/hosts/erp/machine/app/"
	prefix, err := GetKeyWithPrefix(k)
	if err != nil {
		return
	}
	t.Log(prefix)
}
