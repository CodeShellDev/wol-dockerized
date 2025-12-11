package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/wol-dockerized/internals/config"
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
	host := req.Header.Get("X-Forwarded-Host")
	protocol := req.Header.Get("X-Forwarded-Proto")
	uri := req.Header.Get("X-Forwarded-Uri")

	urlStr := protocol + "://" + host + uri
	URL, err := url.Parse(urlStr)
	
	if err != nil {
		logger.Error("Could not parse url: ", urlStr)
		http.Error(w, "Bad Request: missing X-Forwarded headers", http.StatusBadRequest)
		return
	}

	variables := map[string]string{
		"HOSTNAME": URL.Hostname(),
		"HOST": URL.Host,
		"PORT": URL.Port(),
		"PROTOCOL": URL.Scheme,
		"PATH": URL.Path,
	}

	query := buildQuery(config.ENV.QUERY_PATTERN, variables)

	wol.OnActivity(query)

	w.WriteHeader(http.StatusOK)
}

func buildQuery(pattern string, context map[string]string) string {
	result := pattern
	for k, v := range context {
		placeholder := "{" + k + "}"
		result = strings.ReplaceAll(result, placeholder, v)
	}
	return result
}