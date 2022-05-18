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
	//regx 应用资源注册前缀
	midRegx = "mid"

	//regx 应用资源注册前缀
	netRegx = "net"
)

//CellApiServer 上送API网关
func (u *Basest) CellApiServer() error {

	return nil
}

//RegxEtcdKey 继续入库的key
func (u *Basest) RegxEtcdKey() string {
	switch u.Kind {
	case "mid":
		{
			return midRegx + u.Namespace + u.Midtype
		}
	case "net":
		{
			return netRegx + u.Namespace + u.Midtype
		}
	default:
		return ""
	}

}

//RegxEtcValue 继续入库的value
//TODO 格式化value
func (u *Basest) RegxEtcValue() Basest {
	return *u
}

//HttpReq 与API网关交互
func HttpReq(u *Basest) (*http.Request, error) {

	apiServer, apiServerErr := config.GetApiServerUrl()
	if apiServerErr != nil {
		return nil, apiServerErr
	}

	etcdKey := u.RegxEtcdKey()
	etcdValue := u.RegxEtcValue()
	fmt.Printf("HttpReq ==> k:%s \n v:+%v\n", etcdKey, etcdValue)
	insertDB := make(map[string]Basest)
	insertDB[etcdKey] = etcdValue
	apiUrl := apiServer + "/midRegx/set"
	bt, err := json.Marshal(insertDB)
	body := ioutil.NopCloser(strings.NewReader(string(bt)))
	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}
