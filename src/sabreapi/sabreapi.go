package main

import (
	"fmt"
	"log"
	"net/http"
	"sabre/pkg/apiserver"
)

func main() {
	http.HandleFunc("/midRegx/set", apiserver.SetToDB)
	http.HandleFunc("/midRegx/show", apiserver.ShowInfoFromDB)
	//r := NewRouter()
	//r.Use(logger)
	listenPort := "8081"
	fmt.Printf("The listening port of the api server is %s\n", listenPort)
	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
