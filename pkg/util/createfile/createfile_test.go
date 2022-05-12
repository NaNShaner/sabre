package createfile

import (
	"sabre/pkg/sabstruct"
	"testing"
)

func TestYamlToFile(t *testing.T) {
	var Config sabstruct.Config
	Config.Kind = "Config"
	Config.Metadata.Namespace = "default"
	Config.Spec.DefaultConfig.Tomcat.Javaopts = "-server -Xms1024M -Xmx1024M -Xss512k -XX:PermSize=1024M -XX:MaxPermSize=1024M -XX:+DisableExplicitGC -XX:MaxTenuringThreshold=31 -XX:+UseConcMarkSweepGC -XX:+UseParNewGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:+PrintGCDetails -XX:+PrintGCTimeStamps -XX:+PrintHeapAtGC -XX:+HeapDumpOnOutOfMemoryError -XX:+HeapDumpOnCtrlBreak -Djava.awt.headless=true"
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
