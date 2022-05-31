package hostregister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sabre/pkg/config"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	Ti "sabre/pkg/util/tomcat/install"
	"sabre/pkg/yamlfmt"
	"strings"
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
	insertDB := make(map[string]sabstruct.Config)
	contentLength := req.ContentLength
	body := make([]byte, contentLength)
	req.Body.Read(body)
	err := json.Unmarshal(body, &insertDB)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	}
	for k, basestStruct := range insertDB {
		_, err = Ti.Deploy((*commontools.Basest)(&basestStruct))
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(-1)
		}
		for _, host := range basestStruct.DeployHost {
			getServerList := strings.Split(host, "/")
			basestStruct.DeployHostStatus[getServerList[len(getServerList)-1]] = true
		}

		resultJson, err := yamlfmt.PrintResultJson(basestStruct)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		setIntoDBErrOne := dbload.SetIntoDB(k, string(resultJson))
		if setIntoDBErrOne != nil {
			// 失败重试1次
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

	apiServer, apiServerErr := config.GetLetServerUrl()
	if apiServerErr != nil {
		return "", apiServerErr
	}

	//fmt.Printf("k:%s \nv:+%v\n", etcdKey, etcdValue)
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
