package hub

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 最大消息大小
	maxMessageSize = 512
	// 写入超时时间
	writeWait = 10 * time.Second
	// 读超时时间
	pongWait = 60 * time.Second
	// 发送 Ping 的间隔（应小于 pongWait）
	pingPeriod = (pongWait * 9) / 10
)

// Client 代表一个 WebSocket 连接
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient 创建并注册 Client
func NewClient(h *Hub, conn *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		ctx:    ctx,
		cancel: cancel,
	}
	h.register <- client
	return client
}

// readPump 从 WebSocket 连接读取消息，送入 Hub.broadcast
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.cancel()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取错误: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(msg)
		c.hub.broadcast <- msg
	}
}

// NOTE: 方法名也需大写，否则视作私有方法
// writePump 向 WebSocket 连接写入消息，并定期发送 Ping 保持心跳
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭信道，发送关闭帧
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.ctx.Done():
			// 上下文取消，退出
			return
		}
	}
}
