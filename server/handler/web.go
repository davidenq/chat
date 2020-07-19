package handler

import (
	"cad/chat/server/core"
	"net"
	"net/http"
)

//Web .
func Web(c net.Conn, r *http.Request) {

	uuid := r.URL.Query().Get("uuid")
	clientType := r.Header.Get("Client-Type")
	if clientType == "" {
		clientType = "Web"
	}
	ws := core.NewWebSocket(c, r)
	core.CreatePoll()
	core.SetClient(clientType+"-"+uuid, ws)

	ws.Handshake()

	defer ws.Close()
	for {
		frame := ws.Read()
		switch frame.Opcode {
		case 8:
			return
		case 9:
			frame.Opcode = 10
			fallthrough
		case 0, 1:
			fallthrough
		case 2:
			go func(frame core.Frame) {
				core.Broadcast(frame)
			}(frame)
		}
	}
}
