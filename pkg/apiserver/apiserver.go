package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sabre/pkg/config"
	"sabre/pkg/sabstruct"
	"strings"
)

type Basest sabstruct.Config

//ApiServer 字段类型及定义，请见https://github.com/NaNShaner/sabre/blob/master/pkg/dbload/ResType.md
type ApiServer struct {
	Prefix string `json:"prefix"`
	//ResType 资源类型
	ResType string `json:"res_type"`
	//NameSpace 系统简称
	NameSpace string `json:"namespace"`
	//Provider 资源提供方
	Provider string `json:"provider,omitempty"`
	//ProjectNmae 应用工程名称
	ProjectName string `json:"projectname,omitempty"`
}

const (
	//MidRegx 应用资源注册前缀
	MidRegx = "/mid"

	//NetRegx 应用资源注册前缀
	NetRegx = "/net"
)

//CellApiServer 上送API网关
func (u *Basest) CellApiServer() error {
	return nil
}

//RegxEtcdKey 继续入库的key
//TODO 判断资源类型，确定key
func (u *Basest) RegxEtcdKey() string {
	return strings.Join([]string{MidRegx, u.Namespace, u.Midtype}, "/")
}

//RegxEtcValue 继续入库的value
//TODO 格式化value
func (u *Basest) RegxEtcValue() Basest {
	return *u
}

//HttpReq 与API网关交互
func HttpReq(u *Basest) (string, error) {

	apiServer, apiServerErr := config.GetApiServerUrl()
	if apiServerErr != nil {
		return "", apiServerErr
	}

	etcdKey := u.RegxEtcdKey()
	etcdValue := u.RegxEtcValue()
	fmt.Printf("HttpReq ==> k:%s \nv:+%v\n", etcdKey, etcdValue)
	insertDB := make(map[string]Basest)
	//dbInfo.Kname = etcdKey
	//dbInfo.Vname = etcdValue

	insertDB[etcdKey] = etcdValue
	apiUrl := apiServer + "/midRegx/set"
	bt, err := json.Marshal(insertDB)

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
	rspBody, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll failed, url: %s, reqBody: %s, err: %v", apiUrl, reqBody, err)
	}

	return string(rspBody), nil
}
