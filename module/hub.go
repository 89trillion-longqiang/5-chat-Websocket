package module

import (
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"

	"chat/module/protobuf"
)


type Hub struct {
	Clients    map[*Client]bool
	Lock       *sync.Mutex
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Lock :      &sync.Mutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Lock.Lock()
			h.Clients[client] = true
			h.Lock.Unlock()
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {

				fmt.Println("client close :", client.Conn.RemoteAddr())
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

		case message := <-h.Broadcast:
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