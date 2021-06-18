package handle

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"chat/module"
)

func ServeWs(hub *module.Hub, w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	conn, err := module.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("client connect :", conn.RemoteAddr())
	client := &module.Client{Username: username,Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()

	time.Sleep(5 * time.Millisecond)
	client.Online()
}
