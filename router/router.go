package router

import (
	"net/http"

	"chat/handle"
	"chat/module"
)

func SetupHttp(){
	hub := module.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handle.ServeWs(hub, w, r)
	})
}