package callsabrelet

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"io/ioutil"
	"net/http"
	"path"
	"sabre/pkg/apiserver"
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
	u.Timer = commontools.AddNowTime()
	return u
}

//CallSabreletByEachHost 在由commontools.CheckInstallServerBelongToNS()确认机器所属正常后，由该函数调用Sabrelet
//s 为 sabstruct.Config.DeployHost
func (u *Basest) CallSabreletByEachHost(s []string) {
	sabreletServerPort := "18081"
	for _, host := range s {
		sabreletUrl := "http://" + host + ":" + sabreletServerPort + "/hostInfo/Install"
		sab, err := u.CallSabrelet(sabreletUrl, host)
		if err != nil {
			fmt.Printf("===>%s<===", err)
			return
		}
		fmt.Printf("===>%s<===", sab)
	}
}

//CallSabrelet 调用计算节点中的sabrelet，并将返回放入channel中。生产者
//s 为每台计算节点上的saberlet的监听地址
//host 计算节点的ip地址
func (u *Basest) CallSabrelet(s, host string) (string, error) {

	insertDB := make(map[string]Basest)
	// /mid/ERP/Tomcat
	key := path.Join("/mid" + strings.ToUpper(u.Namespace) + strings.ToLower(u.Midtype))
	insertDB[key] = *u

	getStatusReport, err := u.GetStatusReport(host, false)
	if err != nil {
		return "", err
	}

	yml, ymlErr := yamlfmt.PrintResultJson(&insertDB)
	if ymlErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, getStatusReport)
		return "", fmt.Errorf("ymal文件格式化失败, %s\n", ymlErr)
	}
	reqBody := strings.NewReader(string(yml))
	httpReq, httpReqErr := http.NewRequest("POST", s, reqBody)
	//fmt.Printf("请求sabrelet的地址为 %s, 请求报文%+v\n", s, reqBody)
	if httpReqErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, getStatusReport)
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpReqErr)

	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	httpRsp, httpRspErr := http.DefaultClient.Do(httpReq)
	if httpRspErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, getStatusReport)
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpRspErr)
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	rspBody, rspBodyErr := ioutil.ReadAll(httpRsp.Body)
	if rspBodyErr != nil {
		u.DeployHostStatus = append(u.DeployHostStatus, getStatusReport)
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v, response:%s", s, reqBody, rspBodyErr, string(rspBody))
	}
	u.DeployHostStatus = append(u.DeployHostStatus, getStatusReport)
	u.ResolveCallSabreletResponse(u, host)
	return string(rspBody), nil
}

//ResolveCallSabreletResponse 处理sabrelet的返回结果，并更新etcd
//h 为主机ip地址
func (u *Basest) ResolveCallSabreletResponse(yml *Basest, h string) {
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
	fmt.Printf("server %s install %s information Update succeeded %s\n", h, u.Midtype, setInfoToDB)
}

//GetStatusReport 上报服务器状态
func (u *Basest) GetStatusReport(host string, hostStatus bool) (map[string]sabstruct.RunTimeStatus, error) {
	var s sabstruct.RunTimeStatus
	status := make(map[string]sabstruct.RunTimeStatus)
	// 当etcd中的监控记录时间为空时，记录第一次开始监控的时间
	if s.StatusReportTimer == "" {
		s.StatusReportTimer = commontools.AddNowTime()
	}
	s.RunStatus = hostStatus

	status[host] = s
	return status, nil
}

//CalculateRunningDay 获取服务器中中间件的运行时间
//host 为当前主机的ip地址
func (u *Basest) CalculateRunningDay(host string) (int, error) {
	fmt.Printf("==> \n%+v\n", u)
	for _, hStruct := range u.DeployHostStatus {
		for h, hInfo := range hStruct {
			// 判断为当前主机信息，将库中的时间数据和当前时间对比
			if h == host {
				runningDays := carbon.Parse(hInfo.StatusReportTimer).DiffInDays(carbon.Parse(commontools.AddNowTime()))
				return int(runningDays), nil
			}
		}
	}
	//key := path.Join("/mid" + strings.ToUpper(u.Namespace) + strings.ToLower(u.Midtype))
	//typeInfoFromDB, getErr := dbload.GetKeyFromETCD(key, false)
	//if getErr != nil {
	//	return 0, getErr
	//}
	//// 从etcd中获取中间件的信息
	//for _, v := range typeInfoFromDB {
	//	err := json.Unmarshal(v.Value, u)
	//	if err != nil {
	//		return 0, err
	//	}
	//	// 获取DeployHostStatus 判断当前主机的信息
	//	for _, hStruct := range u.DeployHostStatus {
	//		for h, hInfo := range hStruct {
	//			// 判断为当前主机信息，将库中的时间数据和当前时间对比
	//			if h == host {
	//				runningDays := carbon.Parse(hInfo.StatusReportTimer).DiffInDays(carbon.Parse(commontools.AddNowTime()))
	//				return int(runningDays), nil
	//			}
	//		}
	//	}
	//
	//}
	return 0, fmt.Errorf("failed to get runtime\n")
}

//CallFaceOfSabrelet 调用每台机器上的sabrelet
func CallFaceOfSabrelet(h CallSchedule, s []string) {
	h.CallSabreletByEachHost(s)
}
