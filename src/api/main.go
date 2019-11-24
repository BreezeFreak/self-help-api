package main

import (
	"api/config"
	"api/router"
	"api/utils/mongodb"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mongodb.Init()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Server.HttpPort),
		Handler:        router.Init(),
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//if config.Server.RunMode == "release" {
	//	go utils.CronDaily()
	//}

	e := server.ListenAndServe()
	if e != nil {
		log.Fatal(e)
	}
}
