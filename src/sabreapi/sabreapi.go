package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"log"
	"sabre/pkg/apiserver"
	"sabre/pkg/util/hostregister"
)

func main() {
	//http.HandleFunc("/midRegx/set", apiserver.SetToDB)
	//http.HandleFunc("/hostInfo/register", hostregister.RegInfoToDB)
	//http.HandleFunc("/midRegx/show", apiserver.ShowInfoFromDB)

	cntxt := &daemon.Context{
		PidFileName: "/var/run/sabreapi.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/sabreapi.log",
		LogFilePerm: 0640,
		Umask:       027,
		Args:        []string{"[sabreapi]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("sabreapi daemon started.")

	router := gin.Default()

	router.POST("/hostInfo/register", func(context *gin.Context) {
		hostregister.RegInfoToDB(context)
	})

	router.POST("/midRegx/set", func(context *gin.Context) {
		apiserver.SetToDB(context)
	})

	listenPort := "8081"

	runErr := router.Run(":" + listenPort)
	if runErr != nil {
		return
	}

	//fmt.Printf("The listening port of the api server is %s\n", listenPort)
	//httpListenAndServeErr := http.ListenAndServe(":"+listenPort, nil)
	//if httpListenAndServeErr != nil {
	//	log.Fatal(httpListenAndServeErr)
	//}
}
