package core

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
)

//Conn .
type Conn interface {
	Close() error
}

//Frame .
type Frame struct {
	IsFragment bool
	Opcode     byte
	Reserved   byte
	IsMasked   bool
	Length     uint64
	Payload    []byte
}

//WebSocket .
type WebSocket struct {
	WSKey  string
	Conn   Conn
	Buffer *bufio.ReadWriter
	Header http.Header
	Status uint16
}

//NewWebSocket .
func NewWebSocket(c net.Conn, r *http.Request) *WebSocket {
	buffer := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	ws := &WebSocket{
		Conn:   c,
		Buffer: buffer,
		Header: r.Header,
		Status: 1000,
	}
	return ws
}

//Handshake .
func (ws *WebSocket) Handshake() error {

	wsKey := ws.Header.Get("Sec-Websocket-Key")
	hash := generateAcceptedHash(wsKey)
	header := []string{
		"HTTP/1.1 101 Web Socket Protocol Handshake",
		"Server: CAD chat server",
		"Upgrade: WebSocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + hash,
		"",
		"",
	}
	return ws.write([]byte(strings.Join(header, "\r\n")))
}

//Read .
func (ws WebSocket) Read() Frame {
	frame := Frame{}
	head, err := ws.read(2)

	if err != nil {
		return frame
	}

	if head[0] != 129 {
		ra, _, _ := ws.readALL()
		data := "{\"" + string(ra)
		final := []byte(data)
		frame.IsFragment = false
		frame.Opcode = 1
		frame.Reserved = 0
		frame.IsMasked = true
		frame.Length = uint64(len(final))
		frame.Payload = final
		return frame
	}
	frame.IsFragment = (head[0] & 0x80) == 0x00
	frame.Opcode = head[0] & 0x0F
	frame.Reserved = (head[0] & 0x70)
	frame.IsMasked = (head[1] & 0x80) == 0x80

	var length uint64
	length = uint64(head[1] & 0x7F)

	if length == 126 {
		data, err := ws.read(2)
		if err != nil {
			return frame
		}
		length = uint64(binary.BigEndian.Uint16(data))
	} else if length == 127 {
		data, err := ws.read(8)
		if err != nil {
			return frame
		}
		length = uint64(binary.BigEndian.Uint64(data))
	}
	mask, err := ws.read(4)
	if err != nil {
		return frame
	}
	frame.Length = length

	payload, err := ws.read(int(length))
	if err != nil {
		return frame
	}

	for i := uint64(0); i < length; i++ {
		payload[i] ^= mask[i%4]
	}
	frame.Payload = payload
	return frame

}

//Data struct received from clients
type Data struct {
	ID       string `json:"uuid"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

func send(fr Frame) []byte {
	data := make([]byte, 2)
	data[0] = 0x80 | fr.Opcode
	if fr.IsFragment {
		data[0] &= 0x7F
	}

	if fr.Length <= 125 {
		data[1] = byte(fr.Length)
		data = append(data, fr.Payload...)
	} else if fr.Length > 125 && float64(fr.Length) < math.Pow(2, 16) {
		data[1] = byte(126)
		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, uint16(fr.Length))
		data = append(data, size...)
		data = append(data, fr.Payload...)
	} else if float64(fr.Length) >= math.Pow(2, 16) {
		data[1] = byte(127)
		size := make([]byte, 8)
		binary.BigEndian.PutUint64(size, fr.Length)
		data = append(data, size...)
		data = append(data, fr.Payload...)
	}
	return data
}

//Close .
func (ws *WebSocket) Close() error {
	f := Frame{}
	f.Opcode = 8
	f.Length = 2
	f.Payload = make([]byte, 2)
	binary.BigEndian.PutUint16(f.Payload, 1000)
	data := send(f)
	ws.write(data)
	return ws.Conn.Close()
}

func generateAcceptedHash(key string) string {
	magicNumber := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte(magicNumber))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (ws *WebSocket) read(size int) ([]byte, error) {
	bufferSize := 4096
	data := make([]byte, 0)
	for {
		if len(data) == size {
			break
		}
		remaining := size - len(data)
		if bufferSize > remaining {
			bufferSize = remaining
		}
		temp := make([]byte, bufferSize)

		n, err := ws.Buffer.Read(temp)
		if err != nil && err != io.EOF {
			return data, err
		}

		data = append(data, temp[:n]...)
	}
	return data, nil
}

func (ws *WebSocket) readALL() ([]byte, int, error) {
	var length int
	bufferSize := 256
	data := make([]byte, 0)
	for {
		temp := make([]byte, bufferSize)
		n, err := ws.Buffer.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		bufferSize = n
		length += n

		data = append(data, temp[:n]...)

		if (length + 2) > len(data) {
			break
		}
	}
	return data, len(data), nil
}

func (ws *WebSocket) write(data []byte) error {
	if _, err := ws.Buffer.Write(data); err != nil {
		return err
	}
	return ws.Buffer.Flush()
}
