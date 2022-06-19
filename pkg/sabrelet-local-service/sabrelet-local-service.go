package sabrelet_local_service

import (
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"net"
	"path"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"strings"
)

var (
	hostPrex = "/hosts"
)

var h res.MonitorStatus

//LocalServerName 获取本机服务器主机名
func LocalServerName() (string, error) {
	localServerName, getLocalServerNameErr := commontools.GetLocalServerName()
	if getLocalServerNameErr != nil {
		return "", getLocalServerNameErr
	}
	return localServerName, nil
}

//GetInfoList 获取本机服务器IP地址
func GetInfoList(hostname string) (string, error) {

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

//QueryDB 通过qureyKey 查询数据库，获取
func QueryDB(hostname, hostip string) (string, error) {
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
func ReportHostSabreletStatus(r string) (res.Hosts, error) {

	var rs res.Hosts

	getKeyFromETCD, err := dbload.GetKeyFromETCD(r, false)
	if err != nil {
		return res.Hosts{}, err
	}

	if len(getKeyFromETCD) != 1 {
		return res.Hosts{}, fmt.Errorf("query multiple matching data in the database")
	}
	for _, value := range getKeyFromETCD {

		jsonUnmarshalErr := json.Unmarshal(value.Value, &rs)
		if jsonUnmarshalErr != nil {
			return res.Hosts{}, jsonUnmarshalErr
		}
		fmt.Printf("==>%+v\n", rs)
		getMem, getMemErr := rs.GetMem()
		if getMemErr != nil {
			return res.Hosts{}, getMemErr
		}

		rs.Mem = getMem
		rs.OnlineDay = CalculateIntervalDays(rs.OnlineTime)

		return rs, nil

	}
	return res.Hosts{}, fmt.Errorf("failed to report current server status\n")
}

//func ReportHostStatus() {
//	hostName, getHostNameErr := commontools.GetLocalServerName()
//	if getHostNameErr != nil {
//		fmt.Println(getHostNameErr)
//		os.Exit(-1)
//	}
//	kName := KeyName(namespace, hostName, area, f)
//	valueName, err := ValueName(&h, f, namespace, area)
//	if err != nil {
//		fmt.Printf("%s\n", err)
//		os.Exit(-1)
//	}
//	v := make(map[string]res.Hosts)
//	v[kName] = valueName
//	//json, err := yamlfmt.PrintResultJson(v)
//	//if err != nil {
//	//	return
//	//}
//	//fmt.Printf("%s\n", json)
//	reqResp, setHttpReqErr := hostregister.SetHttpReq(kName, valueName)
//	if setHttpReqErr != nil {
//		fmt.Printf("请求sabrelet 失败,%s\n", setHttpReqErr)
//		os.Exit(-1)
//	}
//	fmt.Printf("%s\n", reqResp)
//}
