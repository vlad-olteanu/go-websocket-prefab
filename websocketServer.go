package wsprefab

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// WebsocketServer Only instantiate using constructor!
type WebsocketServer struct {
	Clients    map[*Client]bool
	register   chan *Client
	Unregister chan *Client
	onOpen     func(c *Client)
	onMessage  func(c *Client, message []byte)
	onClose    func(c *Client)
}

func (this *WebsocketServer) run() {
	for {
		select {
		case client := <-this.register:
			this.Clients[client] = true
		case client := <-this.Unregister:
			if _, ok := this.Clients[client]; ok {
				delete(this.Clients, client)
				close(client.Send)
			}
		}
	}
}

func NewWebsocketServer(urlPath string,
	onOpenFn func(c *Client),
	onMessageFn func(c *Client, message []byte),
	onCloseFn func(c *Client),
	kwArgs *KwArgs) *WebsocketServer {
	var instance = new(WebsocketServer)
	instance.onOpen = onOpenFn
	instance.onMessage = onMessageFn
	instance.onClose = onCloseFn
	instance.register = make(chan *Client)
	instance.Unregister = make(chan *Client)
	instance.Clients = make(map[*Client]bool)
	http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		readBufferSize := 1024 * 1024
		writeBufferSize := 1024 * 1024
		writeWait := 10 * time.Second
		pongWait := 60 * time.Second
		pingPeriod := (pongWait * 9) / 10
		maxMessageSize := 1024 * 1024 * 1024 * 16
		if kwArgs != nil {
			if kwArgs.ReadBufferSize != 0 {
				readBufferSize = kwArgs.ReadBufferSize
			}
			if kwArgs.WriteBufferSize != 0 {
				writeBufferSize = kwArgs.WriteBufferSize
			}
			if kwArgs.WriteWait != 0 {
				writeWait = kwArgs.WriteWait
			}
			if kwArgs.PongWait != 0 {
				pongWait = kwArgs.PongWait
			}
			if kwArgs.PingPeriod != 0 {
				pingPeriod = kwArgs.PingPeriod
			}
			if kwArgs.MaxMessageSize != 0 {
				maxMessageSize = kwArgs.MaxMessageSize
			}
		}
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
		}
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{
			WsServer:       instance,
			Conn:           conn,
			Send:           make(chan []byte, 256),
			writeWait:      writeWait,
			pongWait:       pongWait,
			pingPeriod:     pingPeriod,
			maxMessageSize: maxMessageSize,
			Attribs: make(map[string]interface{}),
		}
		instance.register <- client
		go client.OutgoingListener()
		go client.IncomingListener()
		if instance.onOpen != nil {
			instance.onOpen(client)
		}
	})
	go instance.run()
	return instance
}
