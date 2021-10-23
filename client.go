package wsprefab

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	WsServer       *WebsocketServer
	Conn           *websocket.Conn
	Send           chan []byte
	Attribs        map[string]interface{}
	maxMessageSize int
	pingPeriod     time.Duration
	pongWait       time.Duration
	writeWait      time.Duration
}
func (c *Client) IncomingListener() {
	defer func() {
		c.WsServer.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(int64(c.maxMessageSize))
	c.Conn.SetReadDeadline(time.Now().Add(c.pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(c.pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.WsServer.onMessage(c,message)
	}
}
func (c *Client) OutgoingListener() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		if c.WsServer.onClose!=nil{
			c.WsServer.onClose(c)
		}
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}