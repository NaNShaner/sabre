package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path"
	"sabre/pkg/config"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
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
	return path.Join(MidRegx, commontools.FmtETCDKey(u.Namespace), commontools.FmtETCDKey(u.Midtype))
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
	insertDB := make(map[string]Basest)

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

//ShowInfoFromDB 接收 sabrectl show 指令，从etcd中获取数据并反馈
func ShowInfoFromDB(wr http.ResponseWriter, req *http.Request) {
	resp, ok := mux.Vars(req)["kname"]
	if !ok {
		_, err := wr.Write([]byte("查询数据失败"))
		if err != nil {
			return
		}
	}
	//outputMsg, _ := fmt.Printf("查询数据成功\n %+v", mux.Vars(req))
	_, outPutMsgErr := wr.Write([]byte(resp))
	if outPutMsgErr != nil {
		return
	}
}

//SetToDB 从 sabrectl 接收数据，保存进入etcd
func SetToDB(ctx *gin.Context) {
	// var DBStruct apiserver.ToDBServer
	DBStruct := make(map[string]Basest)
	//if readBodyErr != nil {
	//	http.Error(wr, readBodyErr.Error(), http.StatusBadRequest)
	//}

	if err := ctx.ShouldBind(&DBStruct); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	for s, basest := range DBStruct {
		resultJson, err := yamlfmt.PrintResultJson(basest)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		if err := dbload.SetIntoDB(s, string(resultJson)); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
	}
}
