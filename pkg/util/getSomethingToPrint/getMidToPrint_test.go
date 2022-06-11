package getSomethingToPrint

import (
	"testing"
)

//输出：
//namespace	host		midType	projectName	port	version	monitor running	runningTime
//MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
//MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d

func TestOutPutInfo_PrintFmt(t *testing.T) {
	var s OutPutInfo
	s.Namespace = "OICQ"
	s.Host = "192.168.222.222"
	s.MidType = "Tomcat"
	s.AppName = "demo"
	s.Port = "8099"
	s.MidVersion = "7.0.78"
	s.Monitor = true
	s.Running = true
	s.RunningTime = "10d"
	var c CmdArgs

	c.ResType = "mid"
	c.OutPutInfo = s

	err := PrintFmt(c)
	if err != nil {
		t.Error(err)
	}

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
