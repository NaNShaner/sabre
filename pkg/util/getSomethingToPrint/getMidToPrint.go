package getSomethingToPrint

import (
	"encoding/json"
	"fmt"
	"path"
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
	ResType    string
	Namespace  string
	Host       string
	MidType    string
	NetArea    string
	AppName    string
	Port       string
	MidVersion string
	Monitor    bool
	Running    bool
	//RunningTime 运行时间
	//TODO 每次重启需要重置，目前是根据第一次启动时间算起
	RunningTime string
}

var (
	PrintHeader = "%-6s %-15s %s %4s %4s %7s %3s %s %4s\n"
	PrintLine   = "%-9s %-15s %6s %5s %11s %5s %5t %7t %6s\n"
)

//PrintFmt 接收命令行参数,将etcd中的数据格式化后输出到终端
//r 资源类型，例如 mid
//n 系统简称，例如 ERP
//m 资源种类，例如 Tomcat
func PrintFmt(r, n, m string) error {
	dbKey := path.Join("/", r, n, m)
	willOutPutResult, getOutPutResultErr := UseKeyGetInfoFromDB(dbKey)
	if getOutPutResultErr != nil {
		return getOutPutResultErr
	}
	if len(willOutPutResult) != 0 {
		fmt.Printf(PrintHeader, "namespace", "host", "midType", "projectName", "port", "version", "monitor", "running", "runningTime")
		for _, info := range willOutPutResult {
			fmt.Printf(PrintLine, info.Namespace, info.Host, info.MidType, info.AppName, info.Port, info.MidVersion, info.Monitor, info.Running, info.RunningTime)
		}
	} else {
		return fmt.Errorf("no data obtained")
	}

	return nil
}

//UseKeyGetInfoFromDB 接收ETCD中的数据并构造需要输出的结构体
func UseKeyGetInfoFromDB(s string) ([]OutPutInfo, error) {
	var O OutPutInfo
	sab := &sabstruct.Config{}
	var getOutPutInfo []OutPutInfo
	useKeyGetInfoFromDB, getErr := dbload.GetKeyFromETCD(s, true)
	if getErr != nil {
		return nil, getErr
	}
	for _, kv := range useKeyGetInfoFromDB {
		err := json.Unmarshal(kv.Value, sab)
		if err != nil {
			return nil, err
		}
		for _, host := range sab.DeployHost {
			O.Namespace = sab.Namespace
			O.Host = host
			O.MidType = sab.Midtype
			O.AppName = sab.Appname
			O.Port = sab.ListeningPort
			O.MidVersion = sab.Version
			O.Monitor = true

			for _, hstatus := range sab.DeployHostStatus {
				for h, hInfo := range hstatus {
					if h == host {
						O.Running = hInfo.RunStatus
						O.RunningTime = string(rune(hInfo.RunningDays)) + "d"
					}
				}
			}
			getOutPutInfo = append(getOutPutInfo, O)
		}
	}
	return getOutPutInfo, nil
}
