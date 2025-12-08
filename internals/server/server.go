package server

import (
	"net/http"

	"github.com/codeshelldev/gotl/pkg/logger"
)

func Handle() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/wake", wakeHandler)

	mux.HandleFunc("/activity", activityHandler)

    mux.HandleFunc("/ws", websocketHandler)

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Method, " ", r.URL.Path, " ", r.URL.RawQuery)

		mux.ServeHTTP(w, r)
	})

	return final
}