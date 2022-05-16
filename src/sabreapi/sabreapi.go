package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

//SetInfoToDB 从 sabrectl 接收数据，保存进入etcd
func SetInfoToDB(wr http.ResponseWriter, r *http.Request) {
	msg, msgErr := ioutil.ReadAll(r.Body)
	if msgErr != nil {
		return
	}

	logger.Writer()

	//dbloadErr := dbload.SetIntoDB()
	//if dbloadErr != nil {
	//	_, _ = wr.Write([]byte(dbloadErr.Error()))
	//	return
	//}
	_, err := wr.Write(msg)
	if err != nil {
		return
	}
}

//ShowInfoFromDB 接收 sabrectl show 指令，从etcd中获取数据并反馈
func ShowInfoFromDB(wr http.ResponseWriter, r *http.Request) {

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
	http.HandleFunc("/midRegx/set", SetInfoToDB)
	http.HandleFunc("/midRegx/show", ShowInfoFromDB)
	//r := NewRouter()
	//r.Use(logger)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
