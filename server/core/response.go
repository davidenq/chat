package core

import (
	"encoding/json"
	"net/http"
	"strings"
)

type response struct {
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

//Reply prepare the data and response to the client as a json struct.
func Reply(w http.ResponseWriter, data interface{}) {

	result := &response{
		Data:    data,
		Status:  "ok",
		Code:    http.StatusOK,
		Message: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	outcome, err := json.Marshal(result)

	if err != nil {
		result := &response{
			Data:    nil,
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		outcome, err = json.Marshal(result)
	}
	w.Write(outcome)
	return
}

//Broadcast sends message all connected clients
func Broadcast(fr Frame) {
	clients := GetClients()
	for _, client := range clients {
		var data []byte
		if strings.Contains(client.ID, "console") {
			data = fr.Payload
		} else {
			data = send(fr)
		}
		client.WS.write(data)
	}
}
