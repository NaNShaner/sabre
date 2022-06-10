package getSomethingToPrint

import (
	"fmt"
	"testing"
)

//输出：
//namespace	host		midType	projectName	port	version	monitor running	runningTime
//MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
//MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d

func TestOutPutInfo_PrintFmt(t *testing.T) {
	var s OutPutInfo
	s.AppName = "MNPP"
	s.Host = "192.168.222.222"
	s.MidType = "Tomcat"
	s.AppName = "demo"
	s.Port = "8099"
	s.MidVersion = "7.0.78"
	s.Monitor = true
	s.Running = true
	s.RunningTime = "10d"
	Pfmt := "%-6s %-15s %s %4s %4s %7s %3s %s %4s\n"

	fmt.Printf(Pfmt, "namespace", "host", "midType", "projectName", "port", "version", "monitor", "running", "runningTime")
	fmt.Printf("%-9s %-15s %6s %5s %11s %5s %5t %7t %6s\n", s.AppName, s.Host, s.MidType, s.AppName, s.Port, s.MidVersion, s.Monitor, s.Running, s.RunningTime)
}

func TestFmtDBKey(t *testing.T) {
	var s CmdArgs
	s.ResType = "mid"
	s.Namespace = "MNPP"
	s.Host = "192.168.222.222"
	s.MidType = "Tomcat"
	s.AppName = "demo"
	s.Port = "8099"
	s.MidVersion = "7.0.78"
	s.Monitor = true
	s.Running = true
	s.RunningTime = "10d"
	s.NetArea = "APP"

	dbKey, err := FmtDBKey(s)
	if err != nil {
		return
	}
	t.Log(dbKey)

}
