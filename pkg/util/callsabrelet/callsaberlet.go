package callsabrelet

import (
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"io/ioutil"
	"net/http"
	"sabre/pkg/apiserver"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
	"strings"
	"time"
)

type Basest sabstruct.Config

//ReturnMsg 调用saberlet的返回信息结构
type ReturnMsg struct {
	Host   string
	Status bool
	msg    string
}

type CallSchedule interface {
	CallSabreletByEachHost(s []string)
}

type AddNowTime interface {
	AddNowTimeByEachHost() *Basest
}

func (u *Basest) AddNowTimeByEachHost() *Basest {
	getPoint := (*commontools.Basest)(u)
	u.Timer = getPoint.AddNowTime()
	return u
}

//CallSabreletByEachHost 在由commontools.CheckInstallServerBelongToNS()确认机器所属正常后，由该函数调用Sabrelet
//s 为 sabstruct.Config.DeployHost
func (u *Basest) CallSabreletByEachHost(s []string) {
	sabreletServerPort := "18081"
	for _, host := range s {
		sabreletUrl := "http://" + host + ":" + sabreletServerPort + "/hostInfo/Install"
		_, err := u.CallSabrelet(sabreletUrl, host)
		if err != nil {
			return
		}
	}
}

//CallSabrelet 调用计算节点中的sabrelet，并将返回放入channel中。生产者
//s 为每台计算节点上的saberlet的监听地址
//host 计算节点的ip地址
func (u *Basest) CallSabrelet(s, host string) (string, error) {
	insertDB := make(map[string]Basest)
	// /mid/ERP/Tomcat
	key := "/mid" + u.Namespace + "/" + u.Midtype
	insertDB[key] = *u
	yml, ymlErr := yamlfmt.PrintResultJson(&insertDB)
	if ymlErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, u.GetStatusReport(host, false))
		return "", fmt.Errorf("ymal文件格式化失败, %s\n", ymlErr)
	}
	reqBody := strings.NewReader(string(yml))
	httpReq, httpReqErr := http.NewRequest("POST", s, reqBody)
	//fmt.Printf("请求sabrelet的地址为 %s, 请求报文%+v\n", s, reqBody)
	if httpReqErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, u.GetStatusReport(host, false))
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpReqErr)

	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	httpRsp, httpRspErr := http.DefaultClient.Do(httpReq)
	if httpRspErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, u.GetStatusReport(host, false))
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpRspErr)
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	rspBody, rspBodyErr := ioutil.ReadAll(httpRsp.Body)
	if rspBodyErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, u.GetStatusReport(host, false))
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v, response:%s", s, reqBody, rspBodyErr, string(rspBody))
	}
	u.DeployHostStatus = append(u.DeployHostStatus, u.GetStatusReport(host, true))
	u.ResolveCallSabreletResponse(u)
	return string(rspBody), nil
}

//ResolveCallSabreletResponse 处理sabrelet的返回结果，并更新etcd
func (u *Basest) ResolveCallSabreletResponse(yml *Basest) {
	// 添加当前时间
	u.AddNowTimeByEachHost()
	setInfoToDB, setInfoToDBErr := apiserver.HttpReq((*apiserver.Basest)(yml))
	if setInfoToDBErr != nil {
		// 调用失败1秒后，retry 1次
		time.Sleep(time.Second)
		_, _ = apiserver.HttpReq((*apiserver.Basest)(u))
		fmt.Printf("%s\n", setInfoToDBErr)
		// TODO 如果retry失败如何处理
	}
	fmt.Printf("%s install information Update succeeded，%s\n", u.Midtype, setInfoToDB)
}

//GetStatusReport 上报服务器状态
func (u *Basest) GetStatusReport(host string, hostStatus bool) map[string]sabstruct.RunTimeStatus {
	var s sabstruct.RunTimeStatus
	status := make(map[string]sabstruct.RunTimeStatus)
	s.StatusReportTimer = (*commontools.Basest)(u).AddNowTime()
	s.RunStatus = hostStatus
	status[host] = s
	return status
}

//CalculateRunningDay 获取服务器中中间件的运行时间
//host 为当前主机的ip地址
func (u *Basest) CalculateRunningDay(host string) (int, error) {
	key := "/mid" + u.Namespace + "/" + u.Midtype
	typeInfoFromDB, getErr := dbload.GetKeyFromETCD(key, false)
	if getErr != nil {
		return 0, getErr
	}
	// 从etcd中获取中间件的信息
	for _, v := range typeInfoFromDB {
		err := json.Unmarshal(v.Value, u)
		if err != nil {
			return 0, err
		}
		// 获取DeployHostStatus 判断当前主机的信息
		for _, hStruct := range u.DeployHostStatus {
			for h, hInfo := range hStruct {
				// 判断为当前主机信息，将库中的时间数据和当前时间对比
				if h == host {
					runningDays := carbon.Parse((*commontools.Basest)(u).AddNowTime()).DiffInDays(carbon.Parse(hInfo.StatusReportTimer))
					return int(runningDays), nil
				}
			}
		}

	}
	return 0, fmt.Errorf("failed to get runtime\n")
}

//CallFaceOfSabrelet 调用每台机器上的sabrelet
func CallFaceOfSabrelet(h CallSchedule, s []string) {
	h.CallSabreletByEachHost(s)
}
