package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"chat/protobuf"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)
const (
	// SystemMessage 系统消息
	SystemMessage = iota
	// Talk 广播消息(正常的消息)
	Talk
	// HeartBeatMessage 心跳消息
	HeartBeatMessage
	// ConnectedMessage 上线通知
	ConnectedMessage
	// Exit 下线通知
	Exit
)
var addrC = flag.String("addr", "localhost:8080", "http service address")
var username = flag.String("user", "user","username")
func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addrC, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	head := http.Header{"username": {*username}}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), head)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	read := make(chan int)
	send := make(chan int)
	//读
	go func() {
		defer close(read)
		for  {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			var recCom protobuf.Communication
			proto.Unmarshal(message, &recCom)
			log.Printf("recv: %s", recCom.Msg)
		}


	}()
	//写
	go func() {
		defer close(send)
		for {
			var sendMsg string
			fmt.Scan(&sendMsg)
			switch sendMsg {
			case "Exit":
				{
					snedByte , err_pro :=proto.Marshal(&protobuf.Communication{Class: "Exit",Msg:sendMsg})
					if err_pro != nil{
						log.Println("proto.Marshal:", err_pro)
						return
					}
					err := c.WriteMessage(Exit, snedByte)
					if err != nil {
						return 
					}
					send <- 0
					return
				}
			default:
				snedByte , err_pro :=proto.Marshal(&protobuf.Communication{Class: "Talk",Msg:sendMsg})
				if err_pro != nil{
					log.Println("proto.Marshal:", err_pro)
				}else {
					err := c.WriteMessage(Talk, snedByte)
					if err != nil {
						log.Println("write:", err)
						return
					}
				}

			}
		}
	}()

	select {
		case r := <-send :
			if r == 0{
				return
			}
	}
}