package module

import (
	"chat/protobuf"
	"google.golang.org/protobuf/proto"
	"log"
)


type Hub struct {
	Clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				userlist := ""
				for k, _ := range h.Clients {
					userlist = userlist + "\n"+k.Username
				}
				sendByte,_ := proto.Marshal(&protobuf.Communication{Class: "userlist",Msg: userlist})
				for k,_ := range h.Clients {
					err := k.Conn.WriteMessage(1, sendByte)
					if err != nil {
						log.Println(err)
					}
				}
			}

		case message := <-h.broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}