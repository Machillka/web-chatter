package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/machillka/web-chatter/internal/hub"
)

// 全局 Hub 实例
var chatHub = hub.NewHub()

func init() {
	// 在后台启动 Hub 事件循环
	go chatHub.Run()
}

// upgrader 用于将 HTTP 请求升级为 WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 根据需要做跨域验证，这里允许所有
		return true
	},
}

// ChatHandler 完成 WS 握手，并启动读写协程
func ChatHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket 握手失败"})
		return
	}

	client := hub.NewClient(chatHub, conn)
	go client.WritePump()
	client.ReadPump()
}
