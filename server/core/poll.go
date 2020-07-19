package core

import (
	"sync"
)

//Client .
type Client struct {
	ID string `json:"uuid"`
	WS *WebSocket
}

var once sync.Once
var singleton []*Client

//CreatePoll .
func CreatePoll() []*Client {
	once.Do(func() {
		singleton = make([]*Client, 0)
	})
	return singleton
}

//SetClient .
func SetClient(id string, ws *WebSocket) {
	clientWS := &Client{
		ID: id,
		WS: ws,
	}
	singleton = append(singleton, clientWS)
}

//GetClients .
func GetClients() []*Client {
	return singleton
}
