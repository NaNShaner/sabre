package hostregister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sabre/pkg/config"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"strings"
)

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

//KeyName 入库的key名称
// n namespace
// h hostname
// a area
// e.g. /hosts/erp/machine/app/hostname
func KeyName(n, h, a string) string {
	regx := "/hosts"
	resType := "/machine"
	return path.Join(regx, n, resType, a, h)
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
	fmt.Printf("apiUrl:%s \n", apiUrl)
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
