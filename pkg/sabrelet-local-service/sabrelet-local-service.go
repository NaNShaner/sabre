package sabrelet_local_service

import (
	"fmt"
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

type HostMonitor struct {
}

func FmtQueryKey() (string, error) {
	localServerName, getLocalServerNameErr := commontools.GetLocalServerName()
	if getLocalServerNameErr != nil {
		return "", getLocalServerNameErr
	}
	return localServerName, nil
}

func GetInfoList() (string, error) {
	localServerName, err := FmtQueryKey()
	if err != nil {

	}
	queryKey := path.Join(hostPrex, localServerName)
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

//ReportHostSabreletStatus 周期性上报当前服务器的sabrelet的状态
//TODO：逻辑待补充
func ReportHostSabreletStatus() {
	var ms res.HostRegister
	getMem, err := ms.GetMem()
	if err != nil {
		return
	}
	h := ms.ReturnHost()
	h.Online = true
	h.Mem = getMem

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
