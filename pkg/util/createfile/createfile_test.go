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
	Config.ApiServer = "http://126.1.1.1"
	Config.Metadata.Namespace = "ERP"
	Config.Netarea = "APP"
	Config.Appname = "demo"
	Config.Midtype = "Tomcat"
	Config.Version = "7.0.78"
	Config.Port = "8099"
	Config.InstallPath = "/u01/app"
	Config.PKGDownloadPath = "http://126.1.1.1"
	Config.MidRunType = []string{"cluster", "standalone"}
	Config.Name = "miduser"
	Config.Group = "miduser"
	Config.DeployHost = []string{"127.0.0.1", "127.0.0.2"}
	Config.Action = "Install"
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
