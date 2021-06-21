package module

import (
	"fmt"
	"log"
	"time"

	"chat/protobuf"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)



var (
	newline = []byte{'\n'}
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Username string
	Hub *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
	}
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msgType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var recCom protobuf.Communication
		errUn := proto.Unmarshal(message,&recCom)
		if errUn != nil {
			log.Println(err)
			return
		}
		switch msgType {
			case Talk:{
				if string(recCom.Msg) == "userlist"{
					userList := ""
					for k, _ := range c.Hub.Clients {
						userList = k.Username + "\n" + userList
					}
					sendByte,_ := proto.Marshal(&protobuf.Communication{Class: "Talk",Msg: userList})
					err := c.Conn.WriteMessage(Talk, sendByte)
					if err != nil {
						log.Println(err)
					}
				}else {
					log.Println(c.Username,":",string(recCom.Msg))
					broMeg,_ := proto.Marshal(&protobuf.Communication{Class: "Talk",Msg: string(fmt.Sprintf("%s:%s", c.Username, recCom.Msg))})
					//broMeg = bytes.TrimSpace(bytes.Replace(broMeg, newline, space, -1))
					c.Hub.broadcast <- broMeg
				}
			}
			case Exit:{
				log.Println(c.Username,":",recCom.Msg)
				return
			}
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
			}
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client)Online(){
	userList := ""
	c.Hub.Lock.Lock()
	for k, v := range c.Hub.Clients {
		if v == true {
			userList = k.Username + "\n" + userList
		}
	}
	c.Hub.Lock.Unlock()
	sendByte,_ := proto.Marshal(&protobuf.Communication{Class:"userlist",Msg: userList})
	c.Hub.broadcast <- sendByte
}
