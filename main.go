package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"chat/module"
	"chat/router"
)


var addr = flag.String("addr", ":8080", "http service address")
func main() {
	os.Mkdir("log", 0777)
	file, err := os.OpenFile("log/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	log.SetOutput(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	flag.Parse()
	hub := module.NewHub()
	go hub.Run()

	router.SetupHttp(hub)
	err_http := http.ListenAndServe(*addr, nil)
	if err_http != nil {
		log.Fatal("ListenAndServe: ", err_http)
	}

}
