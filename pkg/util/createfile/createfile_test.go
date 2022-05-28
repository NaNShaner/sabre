package createfile

import (
	"sabre/pkg/sabstruct"
	"testing"
)

func TestYamlToFile(t *testing.T) {
	var Config sabstruct.Config
	Config.ApiVersion = "beta"
	Config.Kind = "Config"
	Config.Server = "http://126.1.1.1"
	Config.ApiServer = "http://127.0.0.1:8001/"
	Config.LocalLetServer = "http://127.0.0.1:18001/"
	Config.Metadata.Namespace = "default"
	Config.Spec.DefaultConfig.Tomcat.Javaopts = "-server -Xms1024M -Xmx1024M -Xss512k"
	Config.Spec.DefaultConfig.ListeningPort = "8099"
	Config.Spec.DefaultConfig.AjpPort = "8009"
	Config.Spec.DefaultConfig.ShutdownPort = "8005"
	err := YamlToFile(Config)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("done")
}
