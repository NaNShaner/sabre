package yamlfmt

import "testing"

func TestYamlFmt(t *testing.T) {
	var Configs Config
	f := "conf.yaml"
	yamlFmt, err := YamlFmt(f, Configs)
	if err != nil {
		return
	}
	t.Logf("%v", yamlFmt)
}
