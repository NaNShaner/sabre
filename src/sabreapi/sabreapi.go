package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sabre/pkg/apiserver"
	"sabre/pkg/dbload"
	"sabre/pkg/yamlfmt"
)

//SetToDB 从 sabrectl 接收数据，保存进入etcd
func SetToDB(wr http.ResponseWriter, req *http.Request) {
	var DBStruct apiserver.ToDBServer

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
	resultJson, err := yamlfmt.PrintResultJson(DBStruct)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	}
	if err := dbload.SetIntoDB(DBStruct.Kname, string(resultJson)); err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
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

}

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) add(router string, h http.Handler) {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	r.mux[router] = mergedHandler
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
