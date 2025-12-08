package config

import (
	"os"
	"strconv"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/wol-dockerized/internals/config/structure"
)

var ENV = &structure.ENV{
	LOG_LEVEL: "info",
	MONITOR_INTERVAL: 30,
	INACTIVITY_THRESHOLD: 600,
}

func Load() {
	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel != "" {
		ENV.LOG_LEVEL = logLevel
	}

	ENV.PORT = os.Getenv("PORT")

	ENV.QUERY_PATTERN = os.Getenv("QUERY_PATTERN")

	monitorInterval := os.Getenv("MONITOR_INTERVAL")

	if monitorInterval != "" {
		interval, err := strconv.Atoi(monitorInterval)

		if err != nil {
			logger.Error("Invalid monitor interval: ", err.Error())
		} else {
			ENV.MONITOR_INTERVAL = interval
		}
	}

	inactivityThreshold := os.Getenv("INACTIVITY_THRESHOLD")

	if inactivityThreshold != "" {
		threshold, err := strconv.Atoi(inactivityThreshold)

		if err != nil {
			logger.Error("Invalid inactivity threshold: ", err.Error())
		} else {
			ENV.INACTIVITY_THRESHOLD = threshold
		}
	}
}

func Log() {
	logger.Dev("Loaded Environment:", ENV)
}