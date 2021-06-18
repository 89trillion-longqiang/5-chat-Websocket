package router

import (
	"net/http"

	"chat/handle"
	"chat/module"
)

func SetupHttp(hub *module.Hub){
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handle.ServeWs(hub, w, r)
	})
}
