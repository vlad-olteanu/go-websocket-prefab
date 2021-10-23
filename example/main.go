package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wsprefab"
)
func onOpen(c *wsprefab.Client) {
	fmt.Println("Connection opened")
	msg,_:=json.Marshal(map[string]interface{}{
		"text":"Connection opened",
	})
	c.Send<-msg
	c.Attribs["messageCount"]=0
}
func onMessage(c *wsprefab.Client, message []byte) {
	fmt.Println("Message received")
	//print message.text
	msg:=make(map[string]interface{})
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return
	}
	if v,ok:=msg["text"];ok{
		fmt.Println(v)
	}
	//redirect message to everyone else
	for client := range c.WsServer.Clients{
		//if client is not sender
		if client!=c{
			client.Send<-message
		}
	}
	//increment client's message count
	c.Attribs["messageCount"]=c.Attribs["messageCount"].(int)+1
}
func onClose(c *wsprefab.Client) {
	fmt.Println("Connection closed")
	fmt.Printf("Client has sent %d messages\n",c.Attribs["messageCount"].(int))
}
func main() {
	wsprefab.NewWebsocketServer("/wstest",
		onOpen, //optional, you can use nil instead
		onMessage,
		onClose, //optional
		&wsprefab.KwArgs{ //optional
			MaxMessageSize: 1024*1024*16,//bytes
		},
	)
	fmt.Println("Websocket server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
