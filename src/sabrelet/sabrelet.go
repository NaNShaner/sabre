package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"log"
	sabrelet_local_service "sabre/pkg/sabrelet-local-service"
	"sabre/pkg/util/hostregister"
)

func main() {
	//http.HandleFunc("/hostInfo/Install", hostregister.GetInfoToInstall)

	router := gin.Default()

	router.POST("/hostInfo/Install", func(context *gin.Context) {
		hostregister.GetInfoToInstall(context)
	})

	//TODO: WorkDir参数化
	cntxt := &daemon.Context{
		PidFileName: "/var/run/sabrelet.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/sabrelet.log",
		LogFilePerm: 0640,
		Umask:       027,
		Args:        []string{"sabrelet"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("sabrelet daemon started.")
	// 后台监控
	go sabrelet_local_service.TimeLoopExecution()
	listenPort := "18081"
	runErr := router.Run(":" + listenPort)
	if runErr != nil {
		return
	}

	// 主goroutine堵塞
	//sig := make(chan os.Signal, 2)
	//signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	//<-sig

}
