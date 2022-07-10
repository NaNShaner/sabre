package hostregister

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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
func RegInfoToDB(ctx *gin.Context) {
	//主机注册时HostStruct 的values是res.Hosts
	//旁路登记注册主机便于sablet进行监控时，values是string

	HostStruct := make(map[string]res.Hosts)
	if err := ctx.ShouldBind(&HostStruct); err != nil {
		// 处理错误请求
		return
	}
	for s, basest := range HostStruct {
		resultJson, err := yamlfmt.PrintResultJson(basest)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		}
		if err := dbload.SetIntoDB(s, string(resultJson)); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		}
		_, _ = ctx.Writer.WriteString(s)
	}

}

//GetInfoToInstall 接收sabreschedule的调用，每台机器上都启动着sabrelet，由scheduled挨个调用下发部署指令
func GetInfoToInstall(ctx *gin.Context) {
	s := make(map[string]sabstruct.Config)

	if err := ctx.ShouldBind(&s); err != nil {
		// 处理错误请求
		return
	}

	printResultJson, printResultJsonErr := yamlfmt.PrintResultJson(s)
	if printResultJsonErr != nil {
		fmt.Printf("saberlet err %s\n", printResultJson)
	}
	fmt.Printf("saberlet info %s\n", printResultJson)

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
			ctx.String(http.StatusBadRequest, fmt.Sprintf("server %s  ,%s\n", localHostIp, deployErr))
			break loop
		}
		//for _, host := range basestStruct.DeployHost {
		//	getServerList := strings.Split(host, "/")
		//	basestStruct.DeployHostStatus[getServerList[len(getServerList)-1]] = true
		//}

		resultJson, err := yamlfmt.PrintResultJson(basestStruct)
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("server %s  ,%s\n", localHostIp, err))
		}
		setIntoDBErrOne := dbload.SetIntoDB(k, string(resultJson))
		if setIntoDBErrOne != nil {
			// 失败等待1秒后重试1次
			time.Sleep(time.Second)
			setIntoDBErrTwo := dbload.SetIntoDB(k, string(resultJson))
			if setIntoDBErrTwo != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("server %s  ,%s\n", localHostIp, setIntoDBErrTwo))
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
	return path.Join(regx, commontools.FmtETCDKey(n), commontools.FmtETCDKey(resType), commontools.FmtETCDKey(a), commontools.FmtETCDKey(h), commontools.FmtETCDKey(i))
}

//ValueName 入库的value
func ValueName(h res.HostRegister, ip, belongTo, area string) (res.Hosts, error) {
	register, err := h.ServerRegister(ip, strings.ToUpper(belongTo), area)
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

//SetHttpReq 与API网关交互，将注册主机信息上送
//主机注册时 etcdValue 为 res.hosts
func SetHttpReq(etcdKey string, etcdValue interface{}) (string, error) {
	// 判断主机是否被注册过
	if err := IsHosted(etcdKey); err != nil {
		return "", err
	}
	apiServer, apiServerErr := config.GetApiServerUrl()
	if apiServerErr != nil || apiServer == "" {
		return "", fmt.Errorf("saberlet server address not found in configuration file %s", apiServerErr)
	}
	insertDB := make(map[string]interface{})
	insertDB[etcdKey] = etcdValue
	apiUrl := apiServer + "/hostInfo/register"
	bt, err := json.Marshal(insertDB)
	if err != nil {
		return "", fmt.Errorf("in SetHttpReq %+v json.Marshal Err", insertDB)
	}
	////fmt.Printf("apiUrl:%s \n", apiUrl)
	reqBody := strings.NewReader(string(bt))
	httpReq, httpReqErr := http.NewRequest("POST", apiUrl, reqBody)
	if httpReqErr != nil {
		return "", fmt.Errorf("NewRequest fail, url: %s, reqBody: %+v, err: %v", apiUrl, reqBody, err)

	}
	httpReq.Header.Add("Content-Type", "application/json")

	// DO: HTTP请求
	httpRsp, httpRspErr := http.DefaultClient.Do(httpReq)
	// 如果接口无法调通，是无法获取到 httpRsp.StatusCode
	if httpRspErr != nil || httpRsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("do http fail, url: %s, reqBody: %+v, err:%v", apiUrl, reqBody, httpRspErr)
	}
	defer httpRsp.Body.Close()
	// Read: HTTP结果
	resFromServer, resFromServerErr := ioutil.ReadAll(httpRsp.Body)
	if resFromServerErr != nil {
		return "", fmt.Errorf("ReadAll failed, url: %s, reqBody: %+v, err: %v", apiUrl, reqBody, resFromServerErr)
	}
	return fmt.Sprintf("Server %s registration succeeded", resFromServer), nil
}

//SetHostListInfoTODB
//TODO: 信息入库，未经过API网关
func SetHostListInfoTODB(k, v string) error {
	err := dbload.SetIntoDB(k, v)
	if err != nil {
		return err
	}
	return nil
}

//AddHostToListSaveToDB 将注册主机信息的key拆分为注册主机列表的key
func AddHostToListSaveToDB(s string) (string, error) {
	sSplit := strings.Split(s, "/")
	if len(sSplit) != 7 {
		return "", fmt.Errorf("Failed to format %s as key of host list\n", s)
	}
	return path.Join("/", sSplit[1], sSplit[5], sSplit[6]), nil
}

//IsHosted 判断主机是否已经被注册
func IsHosted(hosted string) error {
	etcd, err := dbload.GetKeyFromETCD(hosted, false)
	if err != nil {
		return err
	}
	if len(etcd) != 0 {
		return fmt.Errorf("Host %s has been registered\n", hosted)
	}
	return nil
}
