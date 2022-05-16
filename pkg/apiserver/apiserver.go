package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sabre/pkg/util/commontools"
	"strings"
)

type Basest commontools.Basest

type apiserver struct {
	//ResType 资源类型，定义etcd中的前缀
	ResType string
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
func (u *Basest) RegxEtcValue() Basest {
	return *u
}

//HttpReq 与API网关交互
func (u *Basest) HttpReq() (*http.Request, error) {

	apiServer := "http://localhost:8081/into/set"

	etcdKey := u.RegxEtcdKey()
	etcdValue := u.RegxEtcValue()
	insertDB := make(map[string]Basest)
	insertDB[etcdKey] = etcdValue

	bt, err := json.Marshal(insertDB)
	body := ioutil.NopCloser(strings.NewReader(string(bt)))
	req, err := http.NewRequest("POST", apiServer, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}
