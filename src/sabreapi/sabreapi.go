package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/yamlfmt"
)

type Basest sabstruct.Config

//SetToDB 从 sabrectl 接收数据，保存进入etcd
func SetToDB(wr http.ResponseWriter, req *http.Request) {
	// var DBStruct apiserver.ToDBServer
	DBStruct := make(map[string]Basest)
	contentLength := req.ContentLength
	body := make([]byte, contentLength)
	req.Body.Read(body)
	//if readBodyErr != nil {
	//	http.Error(wr, readBodyErr.Error(), http.StatusBadRequest)
	//}
	err := json.Unmarshal(body, &DBStruct)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	}

	for s, basest := range DBStruct {
		resultJson, err := yamlfmt.PrintResultJson(basest)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		if err := dbload.SetIntoDB(s, string(resultJson)); err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}
	}

	//
	//outputMsg := strings.Replace("信息入库成功 SetInfoToDB\n", "SetInfoToDB", string(resultJson), -1)
	//_, outPutMsgErr := wr.Write([]byte(outputMsg))
	//
	//if outPutMsgErr != nil {
	//	return
	//}

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

func main() {
	http.HandleFunc("/midRegx/set", SetToDB)
	http.HandleFunc("/midRegx/show", ShowInfoFromDB)
	//r := NewRouter()
	//r.Use(logger)
	listenPort := "8081"
	fmt.Printf("api server监听端口为 %s\n", listenPort)
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
