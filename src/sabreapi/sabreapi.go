package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"sabre/pkg/apiserver"
	"sabre/pkg/util/hostregister"
)

func main() {
	//http.HandleFunc("/midRegx/set", apiserver.SetToDB)
	//http.HandleFunc("/hostInfo/register", hostregister.RegInfoToDB)
	//http.HandleFunc("/midRegx/show", apiserver.ShowInfoFromDB)

	log.Print("sabreapi daemon started.")

	router := gin.Default()
	//router.SetTrustedProxies([]string{"192.168.1.2"})

	router.POST("/hostInfo/register", func(context *gin.Context) {
		hostregister.RegInfoToDB(context)
	})

	router.POST("/midRegx/set", func(context *gin.Context) {
		apiserver.SetToDB(context)
	})

	//fmt.Printf("The listening port of the api server is %s\n", listenPort)
	//httpListenAndServeErr := http.ListenAndServe(":"+listenPort, nil)
	//if httpListenAndServeErr != nil {
	//	log.Fatal(httpListenAndServeErr)
	//}

	cntxt := &daemon.Context{
		PidFileName: "/var/run/sabreapi.pid",
		PidFilePerm: 0644,
		LogFilePerm: 0640,
		Umask:       027,
		Args:        []string{"sabreapi"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Printf("unable to run: %s", err)
		os.Exit(-1)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	listenPort := "8081"
	runErr := router.Run(":" + listenPort)
	if runErr != nil {
		return
	}

}
