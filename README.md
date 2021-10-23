# Go WebSocket Prefab
A ready to deploy wrapper for the Gorilla WebSocket library

## Installation
```shell
go get github.com/vlad-olteanu/go-websocket-prefab
```
## Usage examples
### Sample project
```GoLang
import (
    wsprefab "github.com/vlad-olteanu/go-websocket-prefab"
)

func onOpen(c *wsprefab.Client) {...}
func onMessage(c *wsprefab.Client, message []byte) {...}
func onClose(c *wsprefab.Client) {...}
func main() {
	wsprefab.NewWebsocketServer("/wstest",
		onOpen, //optional, you can use nil instead
		onMessage,
		onClose, //optional
		&wsprefab.KwArgs{ //optional
			MaxMessageSize: 1024*1024*16,//bytes
		},
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
### Sending messages to the client
```GoLang
c.Send<-[]byte("Message")
c.Send<-jsonData
```
### Storing information about the clients
```GoLang
c.Attribs["messageCount"]=0
```
```GoLang
c.Attribs["messageCount"]=c.Attribs["messageCount"].(int)+1
```
### Closing the connection
```GoLang
c.Conn.Close()
```

Look <a href="https://github.com/vlad-olteanu/go-websocket-prefab/blob/master/example/main.go">here</a>
for a broadcast server example