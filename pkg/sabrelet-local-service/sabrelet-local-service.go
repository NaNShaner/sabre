package sabrelet_local_service

import (
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"net"
	"path"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"strings"
	"time"
)

var (
	hostPrex = "/hosts"
)

//var h res.MonitorStatus

//LocalServerName 获取本机服务器主机名
func LocalServerName() (string, error) {
	localServerName, getLocalServerNameErr := commontools.GetLocalServerName()
	if getLocalServerNameErr != nil {
		return "", getLocalServerNameErr
	}
	return localServerName, nil
}

//GetHostIP 获取本机服务器IP地址
func GetHostIP(hostname string) (string, error) {

	queryKey := path.Join(hostPrex, hostname)
	getKeyFromETCD, err := dbload.GetKeyFromETCD(queryKey, true)
	if err != nil {
		return "", err
	}
	for _, v := range getKeyFromETCD {
		localHostIpSplit := strings.Split(string(v.Key), "/")
		localHostIp := localHostIpSplit[len(localHostIpSplit)-1]
		ip := net.ParseIP(localHostIp)
		if ip == nil {
			return "", fmt.Errorf("Failed to resolve server IP address from etcd, error message is %s:%s\n", v.Key, v.Value)
		}
		_, getHostIPErr := commontools.GetHostIP(ip)
		if getHostIPErr != nil {
			return "", getHostIPErr
		}
		return localHostIp, nil
	}
	return "", fmt.Errorf("Failed to resolve server IP address from etcd, query key value is %s\n", queryKey)
}

//GetServerBelongToNSFromQueryDB 通过qureyKey 查询数据库，获取
func GetServerBelongToNSFromQueryDB(hostname, hostip string) (string, error) {
	qureyKey := path.Join(hostPrex, hostname, hostip)

	getKeyFromETCD, err := dbload.GetKeyFromETCD(qureyKey, false)
	if err != nil {
		return "", err
	}
	if len(getKeyFromETCD) != 0 {
		for _, value := range getKeyFromETCD {
			return string(value.Value), nil
		}
	}
	return "", fmt.Errorf("No data with %s as key was queried in the database\n", qureyKey)
}

func CalculateIntervalDays(s string) int {
	runningDays := carbon.Parse(s).DiffInDays(carbon.Parse(commontools.AddNowTime()))
	return int(runningDays)
}

//ReportHostSabreletStatus 周期性上报当前服务器的sabrelet的状态
//TODO：逻辑待补充
//func ReportHostSabreletStatus(r string) (res.Hosts, error) {
//
//	var rs res.Hosts
//
//	getKeyFromETCD, err := dbload.GetKeyFromETCD(r, false)
//	if err != nil {
//		return res.Hosts{}, err
//	}
//
//	if len(getKeyFromETCD) != 1 {
//		return res.Hosts{}, fmt.Errorf("query multiple matching data in the database")
//	}
//	for _, value := range getKeyFromETCD {
//
//		jsonUnmarshalErr := json.Unmarshal(value.Value, &rs)
//		if jsonUnmarshalErr != nil {
//			return res.Hosts{}, jsonUnmarshalErr
//		}
//		fmt.Printf("==>%+v\n", rs)
//		getMem, getMemErr := rs.GetMem()
//		if getMemErr != nil {
//			return res.Hosts{}, getMemErr
//		}
//
//		rs.Mem = getMem
//		rs.OnlineDay = CalculateIntervalDays(rs.OnlineTime)
//
//		return rs, nil
//
//	}
//	return res.Hosts{}, fmt.Errorf("failed to report current server status\n")
//}

//HostStatusReturnKeyValue 根据kName获取数据库中的value，并return k，v
func HostStatusReturnKeyValue() (string, string, error) {
	localServerName, getServerNameErr := LocalServerName()
	if getServerNameErr != nil {
		return "", "", getServerNameErr
	}

	hostIP, getHostIPErr := GetHostIP(localServerName)
	if getHostIPErr != nil {
		return "", "", getHostIPErr
	}

	kName, getServerBelongToNSFromQueryDBErr := GetServerBelongToNSFromQueryDB(localServerName, hostIP)
	if getServerBelongToNSFromQueryDBErr != nil {
		return "", "", getServerBelongToNSFromQueryDBErr
	}

	kValueFromEtcd, getKeyFromETCDErr := dbload.GetKeyFromETCD(kName, false)
	if getKeyFromETCDErr != nil {
		return "", "", getKeyFromETCDErr
	}

	if len(kValueFromEtcd) == 0 {
		return "", "", fmt.Errorf("The server %s has not been registered and does not belong to any NS\n", kName)
	}

	if len(kValueFromEtcd) != 1 {
		return "", "", fmt.Errorf("The data in the database is not unique, query key is %s, please check\n", kName)
	}
	for _, value := range kValueFromEtcd {
		return kName, string(value.Value), nil
	}
	return "", "", fmt.Errorf("")
}

func UpdateHostedInfoToETCD() error {
	var yamlFmt res.Hosts
	//var yamlFmt sabstruct.Config
	kName, kValue, err := HostStatusReturnKeyValue()
	if err != nil {
		return err
	}

	UnmarshalErr := json.Unmarshal([]byte(kValue), &yamlFmt)
	if UnmarshalErr != nil {
		return UnmarshalErr
	}
	//v := make(map[string]res.Hosts)
	//json, err := yamlfmt.PrintResultJson(v)
	//if err != nil {
	//	return
	//}
	//fmt.Printf("%s\n", json)

	getMem, getMemErr := yamlFmt.GetMem()
	if getMemErr != nil {
		return getMemErr
	}

	yamlFmt.Mem = getMem
	fmt.Printf("-=-=-=- %d -=-=-=-\n", CalculateIntervalDays(yamlFmt.OnlineTime))
	yamlFmt.OnlineDay = CalculateIntervalDays(yamlFmt.OnlineTime)

	y, err := yamlfmt.PrintResultJson(yamlFmt)
	if err != nil {
		return err
	}
	// TODO 未经过API server
	setDbErr := dbload.SetIntoDB(kName, string(y))
	if setDbErr != nil {
		return setDbErr
	}

	fmt.Printf("Server %s information reported successfully, %s\n", kName, y)
	return nil
}

func UpdateMidInfoToETCD() error {
	var yamlFmt sabstruct.Config

	localServerName, err := LocalServerName()
	if err != nil {
		return err
	}

	hostIP, err := GetHostIP(localServerName)
	if err != nil {
		return err
	}

	for _, hostStatus := range yamlFmt.DeployHostStatus {
		for server, serverStatus := range hostStatus {
			if server == hostIP {
				serverStatus.RunStatus = true
				serverStatus.RunningDays = CalculateIntervalDays(yamlFmt.DeployAction.Timer)
			}
		}
	}
	return nil
}

func TimeLoopExecution() {
	for {
		err := UpdateHostedInfoToETCD()
		if err != nil {
			fmt.Printf("TimeLoopExecution exec failed ==> %s\n", err)
		}
		time.Sleep(2 * time.Second)
	}
}
