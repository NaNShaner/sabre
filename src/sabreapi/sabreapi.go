package main

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log"
	"net/http"
	"sabre/pkg/apiserver"
	"sabre/pkg/util/hostregister"
)

func main() {
	http.HandleFunc("/midRegx/set", apiserver.SetToDB)
	http.HandleFunc("/hostInfo/register", hostregister.RegInfoToDB)
	http.HandleFunc("/midRegx/show", apiserver.ShowInfoFromDB)

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

	listenPort := "8081"
	fmt.Printf("The listening port of the api server is %s\n", listenPort)
	httpListenAndServeErr := http.ListenAndServe(":"+listenPort, nil)
	if httpListenAndServeErr != nil {
		log.Fatal(err)
	}
}
