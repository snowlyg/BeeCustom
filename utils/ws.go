package utils

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var (
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan Message)
)

type Message struct {
	Message             string `json:"message"`
	AnnotationIsUpdated bool   `json:"annotation_is_updated"`
}

//初始化数据连接
func InitWs() {
	go handleMessages()
}

//广播发送至页面
func handleMessages() {
	for {
		msg := <-Broadcast
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				LogDebug(fmt.Sprintf("client.WriteJSON error: %v", err))
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
