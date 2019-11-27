package controllers

import (
	"fmt"
	"net/http"

	"BeeCustom/utils"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	BaseController
}

func (c *WebSocketController) Get() {

	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		utils.LogDebug(fmt.Sprintf("Cannot setup WebSocket connection:%v", err))
		return
	}

	//defer ws.Close()

	utils.Clients[ws] = true

	msg := utils.Message{"连接成功", false}
	utils.Broadcast <- msg

	c.ServeJSON()
}
