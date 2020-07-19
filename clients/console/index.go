package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var username string
var encode string
var once sync.Once
var header []byte
var ping bool

//WSConn .
type WSConn struct {
	Conn  net.Conn
	Frame []string
}

//Data .
type Data struct {
	ID       string `json:"uuid"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

//StartDial .
func StartDial() *WSConn {

	conn, err := net.Dial("tcp", ":"+os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(0)
	}

	//Handshake .
	conn.Write(generateHeaderWSFrame())

	wsConn := &WSConn{
		Conn: conn,
	}
	return wsConn
}

//Receive .
func (ws *WSConn) Receive() {
	tmp := make([]byte, 256)
	for {
		n, err := ws.Conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		dataReceived := &Data{}
		err = json.Unmarshal(tmp[:n], dataReceived)
		if err != nil {
			fmt.Println(string(tmp[:n]))
		} else {
			nickname := strings.Replace(dataReceived.Nickname, "\n", "", -1)
			msg := strings.Replace(dataReceived.Message, "\n", "", -1)
			fmt.Println("[" + nickname + "]: " + msg)
		}

	}

	/*bufferSize := 128
	//data := make([]byte, 0)


	//this value is used to calculate the final length of the data
	//with this value we can compare how data increase its size
	//when the length value and length of the data is equal
	//we can know that the read of the buffer has been finished
	var length int
	data := make([]byte, 0)
	for {
		temp := make([]byte, bufferSize)
		n, err := ws.Conn.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		bufferSize = n
		length += n
		data = append(data, temp[:n]...)
		if length == len(data) {
			//in this point is important to check the data received from
			//server in order to check if the connection has begun
			if !ping {
				dataSplit := strings.Split(string(data), ":")
				if len(dataSplit) > 0 {
					ping = true
				}
			} else {
				dataReceived := &Data{}
				json.Unmarshal(data, dataReceived)
				fmt.Println(string(data))
				//fmt.Println(dataReceived.Nickname, dataReceived.Message)
			}
		}
	}*/
}

//Send .
func (ws *WSConn) Send(data string) {

	message := &Data{
		ID:       encode,
		Nickname: username,
		Message:  data,
	}
	rawMessage, _ := json.Marshal(message)
	ws.Conn.Write(rawMessage)
	return
}

//GenerateHeaderWSFrame .
func generateHeaderWSFrame() []byte {

	encode = base64.StdEncoding.EncodeToString([]byte(username))
	arrHeader := []string{
		"GET /?uuid=" + encode + " HTTP/1.1",
		"Connection: Upgrade",
		"User-Agent: console",
		"Client-Type: console",
		"Upgrade: websocket",
		"Sec-WebSocket-Version: 13",
		"Sec-WebSocket-Key: " + encode,
		"Sec-WebSocket-Extensions: permessage-deflate;",
		"",
		"",
	}
	header = []byte(strings.Join(arrHeader, "\r\n"))
	return header
}

//StartChat .
func StartChat() {
	var (
		message string
		wsConn  *WSConn
	)

	for {

		reader := bufio.NewReader(os.Stdin)

		if username == "" {
			message = "Write your nickname to start: "
			fmt.Print(message)
		}

		data, _ := reader.ReadString('\n')
		once.Do(func() {
			username = data
			data = "ping"
			wsConn = StartDial()
		})

		wsConn.Send(data)
		go wsConn.Receive()
	}
}

func main() {
	StartChat()
}
