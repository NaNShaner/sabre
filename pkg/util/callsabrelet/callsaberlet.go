package callsabrelet

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sabre/pkg/apiserver"
	"sabre/pkg/sabstruct"
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

//CallSabreletByEachHost 在由commontools.CheckInstallServerBelongToNS()确认机器所属正常后，由该函数调用Sabrelet
//s 为 sabstruct.Config.DeployHost
func (u *Basest) CallSabreletByEachHost(s []string) {
	sabreletServerPort := "18081"
	for _, host := range s {
		sabreletUrl := "http://" + host + ":" + sabreletServerPort + "/hostInfo/Install"
		//go func() {
		//	time.Sleep(time.Second * 120)
		//	u.CallSabrelet(sabreletUrl, c)
		//	wg.Done()
		//}()
		_, err := u.CallSabrelet(sabreletUrl, host)
		if err != nil {
			return
		}
	}
}

//func (u *Basest) CallSabreletByEachHost(s []string) {
//	c := make(chan ReturnMsg)
//	sabreletServerPort := "18081"
//	var wg sync.WaitGroup
//	wg.Add(len(s))
//	for _, host := range s {
//		sabreletUrl := "http://" + host + ":" + sabreletServerPort + "/hostInfo/Install"
//		//go func() {
//		//	time.Sleep(time.Second * 120)
//		//	u.CallSabrelet(sabreletUrl, c)
//		//	wg.Done()
//		//}()
//		u.CallSabrelet(sabreletUrl, c)
//		go func() {
//			time.Sleep(time.Second * 3)
//			u.ResolveCallSabreletResponse(c)
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}

//CallSabrelet 调用计算节点中的sabrelet，并将返回放入channel中。生产者
//s 为每台计算节点上的saberlet的监听地址
//host 计算节点的ip地址
func (u *Basest) CallSabrelet(s, host string) (string, error) {
	var h []map[string]bool
	insertDB := make(map[string]Basest)

	key := "/mid/ERP/Tomcat"
	insertDB[key] = *u
	yml, ymlErr := yamlfmt.PrintResultJson(&insertDB)
	if ymlErr != nil {
		u.DeployHostStatus = append(h, map[string]bool{host: false})
		return "", fmt.Errorf("ymal文件格式化失败, %s\n", ymlErr)
	}
	reqBody := strings.NewReader(string(yml))
	httpReq, httpReqErr := http.NewRequest("POST", s, reqBody)
	fmt.Printf("请求sabrelet的地址为 %s, 请求报文%+v\n", s, reqBody)
	if httpReqErr != nil {
		//c <- ReturnMsg{
		//	Host:   s,
		//	Status: false,
		//	msg:    fmt.Sprintf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpReqErr),
		//}
		u.DeployHostStatus = append(h, map[string]bool{host: false})
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpReqErr)

	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	httpRsp, httpRspErr := http.DefaultClient.Do(httpReq)
	if httpRspErr != nil {
		//c <- ReturnMsg{
		//	Host:   s,
		//	Status: false,
		//	msg:    fmt.Sprintf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpRspErr),
		//}
		u.DeployHostStatus = append(h, map[string]bool{host: false})
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", s, reqBody, httpRspErr)
	}
	defer httpRsp.Body.Close()

	// Read: HTTP结果
	rspBody, rspBodyErr := ioutil.ReadAll(httpRsp.Body)
	if rspBodyErr != nil {
		//c <- ReturnMsg{
		//	Host:   s,
		//	Status: false,
		//	msg:    fmt.Sprintf("do http fail, url: %s, reqBody: %+v, err:%v, response:%s", s, reqBody, rspBodyErr, string(rspBody)),
		//}
		u.DeployHostStatus = append(h, map[string]bool{host: false})
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v, response:%s", s, reqBody, rspBodyErr, string(rspBody))
	}
	//c <- ReturnMsg{
	//	Host:   s,
	//	Status: true,
	//	msg:    "",
	//}
	//close(c)

	u.DeployHostStatus = append(h, map[string]bool{host: true})
	u.ResolveCallSabreletResponse(u)
	return string(rspBody), nil
}

//ResolveCallSabreletResponse 处理sabrelet的返回结果，并更新etcd
func (u *Basest) ResolveCallSabreletResponse(yml *Basest) {

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

func CallFaceOfSabrelet(h CallSchedule, s []string) {
	h.CallSabreletByEachHost(s)
}

////ResolveCallSabreletResponse 处理sabrelet的返回结果，并更新etcd
//func (u *Basest) ResolveCallSabreletResponse(c chan ReturnMsg) {
//	var h []map[string]bool
//	for n := range c {
//		if n.Status {
//			u.DeployHostStatus = append(h, map[string]bool{n.Host: n.Status})
//		} else {
//			u.DeployHostStatus = append(h, map[string]bool{n.Host: n.Status})
//		}
//	}
//	setInfoToDB, setInfoToDBErr := apiserver.HttpReq((*apiserver.Basest)(u))
//	if setInfoToDBErr != nil {
//		// 调用失败1秒后，retry 1次
//		time.Sleep(time.Second)
//		_, _ = apiserver.HttpReq((*apiserver.Basest)(u))
//		fmt.Printf("%s\n", setInfoToDBErr)
//		// TODO 如果retry失败如何处理
//	}
//	fmt.Printf("%s install information Update succeeded，%s\n", u.Midtype, setInfoToDB)
//}
