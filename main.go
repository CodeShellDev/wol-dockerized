package main

import (
	"net/http"
	"time"

	log "github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/wol-dockerized/internals/config"
	"github.com/codeshelldev/wol-dockerized/internals/docker"
	"github.com/codeshelldev/wol-dockerized/internals/server"
	"github.com/codeshelldev/wol-dockerized/internals/wol"
)

func main() {
	config.Load()

	log.Init(config.ENV.LOG_LEVEL)

	docker.Init()
	log.Info("Initialized Logger with Level of ", log.Level())

	if log.Level() == "dev" {
		log.Dev("Welcome back Developer!")
	}

	config.Log()

	addr := "0.0.0.0:" + config.ENV.PORT

	srv := &http.Server{
		Addr:    addr,
		Handler: server.Handle(),
	}

	go func() {
		ticker := time.NewTicker(time.Duration(config.ENV.MONITOR_INTERVAL) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			wol.Monitor(config.ENV.MONITOR_INTERVAL)
		}
	}()

	stop := docker.Run(func() {
		err := srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", err.Error())
		}
	})

	<-stop

	docker.Shutdown(srv)
}
