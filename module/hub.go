package module

import (

	"sync"

	"chat/module/protobuf"
	"google.golang.org/protobuf/proto"
)


type Hub struct {
	Clients    map[*Client]bool
	Lock *sync.Mutex
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
		Lock : &sync.Mutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {

				h.Lock.Lock()
				delete(h.Clients, client)
				h.Lock.Unlock()
				close(client.Send)

				userlist := ""
				h.Lock.Lock()
				for k, _ := range h.Clients {
					userlist = userlist + "\n"+k.Username
				}
				sendByte,_ := proto.Marshal(&protobuf.Communication{Class: "userlist",Msg: userlist})
				h.SendBroadcast(sendByte)
				h.Lock.Unlock()
			}

		case message := <-h.broadcast:
			h.SendBroadcast(message)
		}
	}
}

func (h Hub)SendBroadcast(sendByte []byte){
	for client := range h.Clients {
		select {
		case client.Send <- sendByte:
		default:
			close(client.Send)
			h.Lock.Lock()
			delete(h.Clients, client)
			h.Lock.Unlock()
		}
	}
}