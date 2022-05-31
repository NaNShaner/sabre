package main

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log"
	"net/http"
	"sabre/pkg/util/hostregister"
)

func main() {
	http.HandleFunc("/hostInfo/register", hostregister.RegInfoToDB)
	http.HandleFunc("/hostInfo/register", hostregister.GetInfoToInstall)
	//TODO: WorkDir参数化
	cntxt := &daemon.Context{
		PidFileName: "/var/run/sabrelet.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/sabrelet.log",
		LogFilePerm: 0640,
		Umask:       027,
		Args:        []string{"[sabrelet]"},
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

	listenPort := "18081"
	fmt.Printf("The listening port of the sabrelet server is %s.\n", listenPort)
	httpErr := http.ListenAndServe(":"+listenPort, nil)
	if httpErr != nil {
		log.Fatal(err)
	}
	// 主goroutine堵塞
	//sig := make(chan os.Signal, 2)
	//signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	//<-sig

}
