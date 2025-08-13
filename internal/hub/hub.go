package hub

import (
	"log"
)

// Hub 负责管理所有活跃的客户端连接和广播消息
type Hub struct {
	// 注册新的客户端
	register chan *Client
	// 注销的客户端
	unregister chan *Client
	// 来自任意客户端的待广播消息
	broadcast chan []byte
	// 当前所有活跃客户端
	clients map[*Client]bool
}

// NewHub 创建并返回一个 Hub 实例
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
	}
}

// Run 启动事件循环，处理注册、注销、广播三类事件
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("新客户端注册，总数：%d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("客户端断开，总数：%d", len(h.clients))
			}

		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					// 客户端无法发送，强制断开
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
