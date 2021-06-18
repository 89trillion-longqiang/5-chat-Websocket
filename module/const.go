package module

import "time"

const(
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
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