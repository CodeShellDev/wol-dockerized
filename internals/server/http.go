package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/wol-dockerized/internals/wol"
)

type RequestBody struct {
    Query        	string 	`json:"query,omitempty"`
}

func wakeHandler(w http.ResponseWriter, req *http.Request) {
    var body RequestBody

    err := json.NewDecoder(req.Body).Decode(&body)
    if err != nil {
        logger.Error("Could not get Request Body: ", err)
        http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
        return
    }

    if body.Query == "" {
        http.Error(w, "Bad Request: missing required fields", http.StatusBadRequest)
        return
    }

	clientID := createID()

    resp := map[string]any{
        "client_id": clientID,
    }

	respBytes, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	logger.Debug("Sending client_id to client")

    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
    w.Write(respBytes)

	f, ok := w.(http.Flusher)
	if ok { 
		f.Flush() 
	}

	logger.Debug("Waiting for client to establish websocket connection")

	client, err := waitForClient(clientID, time.Duration(20 * time.Second))

	if err != nil {
		logger.Error("Could not get client: ", err.Error())
		return
	}

	err = wol.WakeContainers(body.Query)

	if err != nil {
		sendToClient(client, map[string]any{
			"success": false,
			"error": true,
			"message": "Could not start containers",
		})

		closeClient(client)
		return
	}

	sendToClient(client, map[string]any{
		"success": true,
		"error": true,
		"message": "Started containers",
	})

	closeClient(client)
}

func activityHandler(w http.ResponseWriter, req *http.Request) {

}