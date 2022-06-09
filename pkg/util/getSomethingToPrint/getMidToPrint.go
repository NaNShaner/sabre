package getSomethingToPrint

import (
	"encoding/json"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
)

//### 例如：获取erp系统的下的demo工程部署在哪些机器的Tomcat中
//key :/mid/ERP/Tomcat/{projectName}/{hostname/ipaddr}
//sabrectl get mid -t tomcat -a app -n erp

//输出：
//namespace	host		midType	projectName	port	version	monitor running	runningTime
//MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
//MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d

type OutPutInfo struct {
	Namespace   string
	Host        string
	MidType     string
	AppName     string
	Port        string
	MidVersion  string
	Monitor     bool
	Running     bool
	RunningTime string
}

func (o OutPutInfo) PrintFmt() {

}

func UseKeyGetInfoFromDB(s string) []OutPutInfo {
	var O OutPutInfo
	sab := &sabstruct.Config{}
	var getOutPutInfo []OutPutInfo
	useKeyGetInfoFromDB, getErr := dbload.GetKeyFromETCD(s, true)
	if getErr != nil {
		return nil
	}
	for _, kv := range useKeyGetInfoFromDB {
		err := json.Unmarshal(kv.Value, sab)
		if err != nil {
			return nil
		}
		for _, host := range sab.DeployHost {
			O.AppName = sab.Namespace
			O.Host = host
			O.MidType = sab.Midtype
			O.AppName = sab.Appname
			O.Port = sab.ListeningPort
			O.MidVersion = sab.Version
			O.Monitor = true
			O.RunningTime = "10d"
			for _, hstatus := range sab.DeployHostStatus {
				for h, hInfo := range hstatus {
					if h == host {
						O.Running = hInfo.RunStatus
					}
				}
			}
			getOutPutInfo = append(getOutPutInfo, O)
		}

	}
	return getOutPutInfo
}
