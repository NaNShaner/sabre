package main

import (
	"fmt"
	"log"
	"net/http"
	"sabre/pkg/util/hostregister"
)

func main() {
	http.HandleFunc("/hostInfo/register", hostregister.RegInfoToDB)
	//r := NewRouter()
	//r.Use(logger)
	listenPort := "18081"
	fmt.Printf("The listening port of the sabrelet server is %s\n", listenPort)
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
