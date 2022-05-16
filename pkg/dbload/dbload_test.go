package dbload

import (
	"sabre/pkg/sabstruct"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestSetIntoDB(t *testing.T) {
	k := "MNPP"
	f := "/Users/bijingrui/sabre/pkg/yamlfmt/sabrelet.yaml"
	var Configs sabstruct.Config
	yamlFmt, err := yamlfmt.YamlFmt(f, Configs)
	if err != nil {
		t.Error(err)
	}

	SetIntoDBErr := SetIntoDB(k, yamlFmt)
	if SetIntoDBErr != nil {
		t.Errorf("==> %s", SetIntoDBErr)
		return
	}
	t.Log("done")
}
