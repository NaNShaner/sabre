package hostregister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"sabre/pkg/config"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	Ti "sabre/pkg/util/tomcat/install"
	"sabre/pkg/yamlfmt"
	"strings"
	"time"
)

type Basest sabstruct.Config

//RegInfoToDB 由 sabrectl 发起主机注册请求，再由sabrelet保存进入etcd，sabrelet本地缓存
// etcd的key /hosts/erp/machine/app/hostname:res.Hosts
func RegInfoToDB(wr http.ResponseWriter, req *http.Request) {
	HostStruct := make(map[string]res.Hosts)
	contentLength := req.ContentLength
	body := make([]byte, contentLength)
	req.Body.Read(body)
	//if readBodyErr != nil {
	//	http.Error(wr, readBodyErr.Error(), http.StatusBadRequest)
	//}
	err := json.Unmarshal(body, &HostStruct)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	}

	for s, host := range HostStruct {
		resultJson, err := yamlfmt.PrintResultJson(host)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		if err := dbload.SetIntoDB(s, string(resultJson)); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = wr.Write([]byte(s))
	}

}

//GetInfoToInstall 接收sabreschedule的调用，每台机器上都启动着sabrelet，由scheduled挨个调用下发部署指令
func GetInfoToInstall(wr http.ResponseWriter, req *http.Request) {
	s := make(map[string]sabstruct.Config)
	contentLength := req.ContentLength
	body := make([]byte, contentLength)
	req.Body.Read(body)
	jsonUnmarshalErr := json.Unmarshal(body, &s)
	printResultJson, printResultJsonErr := yamlfmt.PrintResultJson(s)
	if printResultJsonErr != nil {
		fmt.Printf("saberlet err %s\n", printResultJson)
	}
	fmt.Printf("saberlet info %s\n", printResultJson)
	if jsonUnmarshalErr != nil {
		http.Error(wr, jsonUnmarshalErr.Error(), http.StatusBadRequest)
	}
loop:
	for k, basestStruct := range s {
		var localHostIp string
		for _, host := range basestStruct.DeployHost {

			ip, _ := commontools.GetHostIP(net.ParseIP(host))
			if ip {
				localHostIp = host
			}
		}
		fmt.Printf("saberlet localHostIp %s\n", localHostIp)

		_, deployErr := Ti.Deploy((*commontools.Basest)(&basestStruct))
		if deployErr != nil {
			fmt.Printf("server %s  ,%s\n", localHostIp, deployErr)
			wr.Write([]byte(fmt.Sprintf("server %s  ,%s\n", localHostIp, deployErr)))
			break loop
		}
		//for _, host := range basestStruct.DeployHost {
		//	getServerList := strings.Split(host, "/")
		//	basestStruct.DeployHostStatus[getServerList[len(getServerList)-1]] = true
		//}

		resultJson, err := yamlfmt.PrintResultJson(basestStruct)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		setIntoDBErrOne := dbload.SetIntoDB(k, string(resultJson))
		if setIntoDBErrOne != nil {
			// 失败等待1秒后重试1次
			time.Sleep(time.Second)
			setIntoDBErrTwo := dbload.SetIntoDB(k, string(resultJson))
			if setIntoDBErrTwo != nil {
				http.Error(wr, setIntoDBErrTwo.Error(), http.StatusBadRequest)
			}
		}
	}
}

//KeyName 入库的key名称
// n namespace
// h hostname
// a area
// i ip addr
// e.g. /hosts/erp/machine/app/hostname
func KeyName(n, h, a, i string) string {
	regx := "/hosts"
	resType := "/machine"
	return path.Join(regx, n, resType, a, h, i)
}

//ValueName 入库的value
func ValueName(h res.HostRegister, ip, beloogto, area string) (res.Hosts, error) {
	register, err := h.ServerRegister(ip, beloogto, area)
	if err != nil {
		return res.Hosts{}, err
	}
	ipFmt, err := h.IpFmt(ip)
	if err != nil {
		return res.Hosts{}, err
	}
	if _, err := commontools.GetHostIP(ipFmt); err != nil {
		return res.Hosts{}, err
	}
	return register, nil
}

//SetHttpReq 与API网关交互
func SetHttpReq(etcdKey string, etcdValue res.Hosts) (string, error) {

	apiServer, apiServerErr := config.GetApiServerUrl()
	if apiServerErr != nil || apiServer == "" {
		return "", fmt.Errorf("saberlet server address not found in configuration file %s", apiServerErr)
	}
	// TODO: 方便调试，后续删除
	//fmt.Printf("apiServer:%s \n", apiServer)
	insertDB := make(map[string]res.Hosts)

	insertDB[etcdKey] = etcdValue
	apiUrl := apiServer + "/hostInfo/register"
	bt, err := json.Marshal(insertDB)
	//fmt.Printf("apiUrl:%s \n", apiUrl)
	reqBody := strings.NewReader(string(bt))
	httpReq, err := http.NewRequest("POST", apiUrl, reqBody)
	if err != nil {
		return "", fmt.Errorf("NewRequest fail, url: %s, reqBody: %s, err: %v", apiUrl, reqBody, err)

	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %s, err:%v", apiUrl, reqBody, err)
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	resFromServer, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll failed, url: %s, reqBody: %s, err: %v", apiUrl, reqBody, err)
	}
	return fmt.Sprintf("Server %s registration succeeded", resFromServer), nil
}

//ReportHostSabreletStatus 周期性上报当前服务器的saberlet的状态
//TODO：逻辑待补充
func ReportHostSabreletStatus() {

}
