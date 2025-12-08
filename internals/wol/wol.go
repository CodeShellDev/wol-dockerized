package wol

import (
	"errors"
	"maps"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/wol-dockerized/internals/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

var containerQueries map[string][]string
var containerQueryMutex sync.RWMutex

var queryLastActivity = map[string]int64{}
var queryLastActivityMutex sync.RWMutex

const WOL_QUERY = "wol.query"
const WOL_AUTOSTOP = "wol.autostop"

func Monitor(threshold int) {
	if updateContainers() {
		doActivityCheck(threshold)
	}
}

func doActivityCheck(threshold int) {
	threshold64 := int64(threshold)
	
	currentTime := time.Now().Unix()

	queryLastActivityMutex.RLock()
	lastActivities := maps.Clone(queryLastActivity)
	queryLastActivityMutex.RUnlock()

	containerQueryMutex.RLock()
	containerIDs := maps.Clone(containerQueries)
	containerQueryMutex.RUnlock()

	for query, lastTime := range lastActivities {
		if currentTime - lastTime > threshold64 {
			logger.Info("Containers with ", query, " have been flagged for Inactivity")
			logger.Debug("Stopping Containers with ", query)

			ids := containerIDs[query]
			resetLastActivity(query)

			if len(ids) <= 0 {
				continue
			}

			for _, id := range ids {
				if getLabel(id, WOL_AUTOSTOP) == "true" {
					_, err := docker.StopContainer(id, client.ContainerStopOptions{})

					if err != nil {
						logger.Error("Could not stop container: ", err.Error())
					}
				}
			}
		}
	}
}

func WakeContainers(query string) error {
	query = strings.TrimSpace(query)

	containerQueryMutex.RLock()
	containers, exists := containerQueries[query]
	containerQueryMutex.RUnlock()

	if !exists {
		return errors.New("Invalid query")
	}

	logger.Debug("Found ", len(containers), " with query ", query)

	for _, containerID := range containers {
		_, err := docker.StartContainer(containerID, client.ContainerStartOptions{})

		if err != nil {
			logger.Error("Could not start container ", containerID, ": ", err.Error())
			return errors.New("Could not start container")
		}

		err = waitForContainer(containerID)

		if err != nil {
			logger.Error("Container failed to start: ", err.Error())
			return errors.New("Container failed to start: " + err.Error())
		}
	}

	resetLastActivity(query)

	return nil
}

func waitForContainer(id string) error {
	timeout := time.After(15 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return errors.New("Timed out waiting for container")
		case <-ticker.C:
			container, _ := docker.GetContainer(id, client.ContainerInspectOptions{})

			if container.Container.State != nil && container.Container.State.Running {
				return nil
			}
			if container.Container.State != nil && container.Container.State.ExitCode != 0 {
				return errors.New("Container exited with code " + strconv.Itoa(container.Container.State.ExitCode))
			}
		}
	}
}

func updateContainers() bool {
	queries, err := getContainerQueries()

	if err != nil {
		logger.Error("Error getting container queries: ", err.Error())
		return false
	}

	containerQueryMutex.Lock()
	containerQueries = queries
	containerQueryMutex.Unlock()

	return true
}

func getContainerQueries() (map[string][]string, error) {
	containers, err := getEnabledContainers()

	if err != nil {
		return map[string][]string{}, err
	}

	res := map[string][]string{}

	for _, container := range containers {
		query := getLabel(container.ID, "wol.query")
		query = strings.TrimSpace(query)

		if query != "" {
			entry, exists := res[query]

			if !exists {
				entry = []string{}
			}

			entry = append(entry, container.ID)

			res[query] = entry
		}
	}

	return res, nil
}

func resetLastActivity(query string) {
	queryLastActivityMutex.Lock()
	queryLastActivity[query] = time.Now().Unix()
	queryLastActivityMutex.Unlock()
}

func getLabel(id, label string) string {
	container, err := docker.GetContainer(id, client.ContainerInspectOptions{})

	if err != nil {
		return ""
	}

	return container.Container.Config.Labels[label]
}

func getEnabledContainers() ([]container.Summary, error) {
	filters := client.Filters{}
	filters.Add("label", "wol.enable=true")

	return docker.GetContainers(client.ContainerListOptions{
		Filters: filters,
	})
}