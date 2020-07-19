package main

import (
	"bufio"
	"cad/chat/server/handler"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

//GenerateHTTPHeader is used to generate a fake http request through the tcp request.
//The purpose is to generate a sturcture with all the data request and to avoid
//get manually the request information to the client.
func GenerateHTTPHeader(c net.Conn) (*http.Request, error) {
	bufferIO := bufio.NewReader(c)
	httpRequest, err := http.ReadRequest(bufferIO)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
			return nil, err
		}
		return nil, err
	}
	return httpRequest, err
}

func main() {
	tcp, _ := net.Listen("tcp", ":"+os.Getenv("SERVER_PORT"))
	defer tcp.Close()
	for {
		conn, _ := tcp.Accept()
		request, _ := GenerateHTTPHeader(conn)
		go handler.Web(conn, request)
	}
}
