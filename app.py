from flask import Flask, Response, request, jsonify
from docker import DockerClient, errors as DockerErrors
from urllib.parse import urlparse as parseUrl
import time
import os
import threading
import logging

app = Flask("WoL Dockerized")

app.logger.setLevel(logging.INFO)

MONITOR_INTERVAL = os.getenv("MONITOR_INTERVAL")
INACTIVITY_THRESHOLD = os.getenv("INACTIVITY_THRESHOLD")
PATTERN = os.getenv("PATTERN")

dockerClient = DockerClient(base_url="unix://var/run/docker.sock")

lastAccessLock = threading.Lock()
containerQueryLock = threading.Lock()

containerQueryMap = {}

containerLastAccess = {}

@app.route("/endpoint", methods=["GET"])
def listenOnEndpoint():
    host = request.headers.get('X-Forwarded-Host')
    protocol = request.headers.get('X-Forwarded-Proto')
    uri = request.headers.get('X-Forwarded-Uri')

    url = f"{protocol}://{host}{uri}"

    parsedUrl = parseUrl(url)

    port = ""

    if parsedUrl.port:
        port = parsedUrl.port

    context = {
        "HOSTNAME": parsedUrl.hostname,
        "HOST": host,
        "PORT": port,
        "PROTOCOL": protocol,
        "URL": url,
        "PATH": parsedUrl.path 
    }

    query = buildQuery(PATTERN, context)

    resetLastAccessTime(query)

    return Response(status=200)

def buildQuery(pattern, context):
    query = pattern
    for key, value in context.items():
        query = query.replace("{"+key+"}", value)

    return query

def isAutoStopEnabled(containerId):
    containerInfo = dockerClient.inspect_container(container_id)

    labels = container_info.get('Config', {}).get('Labels', {})

    autoStopLabel = "wol.autostop"
    
    if autoStopLabel in labels:
        return labels.get("wol.auto_stop") == "true"
    else:
        return False

def checkInactive():
    while True:
        time.sleep(int(MONITOR_INTERVAL))

        currentTime = time.time()
        with lastAccessLock:
            for query, lastAccessed in list(containerLastAccess.items()):
                if currentTime - lastAccessed > int(INACTIVITY_THRESHOLD):
                    infoLog(f"Container/s running as {query} has been inactive for more than {INACTIVITY_THRESHOLD} seconds.")
                    infoLog("Stopping Container/s")

                    containerIds = tryMatchQueryToContainers(query)
                    resetLastAccessTime(query, False)

                    if containerIds:
                        for containerId in containerIds:
                            if not isAutoStopEnabled(containerId):
                                stopContainer(containerId)


@app.route("/", methods=["POST"])
def listen():
    data = request.json

    success = False

    if data:
        query = data.get("query", None)

        if query:
            infoLog(f"Searching for Container/s matching query: '{query}'")
            
            containerIds = tryMatchQueryToContainers(query)

            successAll = True

            if containerIds:
                infoLog(f"Found atleast 1 Container matching '{query}'")
                
                for containerId in containerIds:
                    success = startContainer(containerId)
                    resetLastAccessTime(query)
                    
                    if success and successAll:
                        successAll = True
                    elif not success:
                        successAll = False
                
                if successAll:
                    success = True

    response = {
        "success": success
    }

    return jsonify(message=response)

def getEnabledContainers():
    try:
        containers = dockerClient.containers.list(all=True, filters={"label": "wol.enable=true"})
        
        containerQueries = {}
        for container in containers:
            query = container.labels.get("wol.query")
            if query:
                if not containerQueries.get(query, None):
                    containerQueries[query] = []

                containerQueries[query].append(container.id)

        return containerQueries
    except Exception as e:
        infoLog(f"Error fetching containers: {str(e)}")
        return {}

def tryMatchQueryToContainers(query):
    containerId = []

    with containerQueryLock:
        if containerQueryMap:
            containerId = containerQueryMap.get(query, [])

    return containerId

def startContainer(containerId):
    try:
        container = dockerClient.containers.get(containerId)

        if container.status == "running" or container.status == "starting":
            return True

        container.start()
        return True
    except DockerErrors.NotFound:
        return False
    except Exception as e:
        return False

def stopContainer(containerId):
    try:
        container = dockerClient.containers.get(containerId)

        container.stop()
        return True
    except DockerErrors.NotFound:
        return False
    except Exception as e:
        return False

def startBackgroundThread():
    thread = threading.Thread(target=checkInactive, daemon=True)
    thread.start()

def resetLastAccessTime(query, lock=True):
    if lock:
        with lastAccessLock:
            containerLastAccess[query] = time.time()
    else:
        containerLastAccess[query] = time.time()

def infoLog(msg):
    app.logger.info(msg)

if __name__ == '__main__':
    if not PATTERN or PATTERN == "":
        print(f"No PATTERN set")
    else:
        if not MONITOR_INTERVAL:
            print(f"No MONITOR_INTERVAL set, using 60sec as default")
            MONITOR_INTERVAL = 60
            
        if not INACTIVITY_THRESHOLD:
            print(f"No INACTIVITY_THRESHOLD set, using 600sec as default")
            INACTIVITY_THRESHOLD = 600
            
        with containerQueryLock:
            containerQueryMap = getEnabledContainers()

            print("Enabled containers on startup:")

            for query, id in containerQueryMap.items():
                print(f"Container Query: {query}, ID: {id}")
                resetLastAccessTime(query)
        
        startBackgroundThread()
        app.run(debug=False, port=7777, host='0.0.0.0')
