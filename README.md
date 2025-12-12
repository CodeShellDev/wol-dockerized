<h1 align="center">WoL Dockerized</h1>

<p align="center">
🖥️ Wake Docker Containers · Auto-Stop Inactive Containers · HTTP & WebSocket Interface
</p>

<div align="center">
  <a href="https://github.com/codeshelldev/wol-dockerized/releases">
    <img 
        src="https://img.shields.io/github/v/release/codeshelldev/wol-dockerized?sort=semver&logo=github&label=Release" 
        alt="GitHub release"
    >
  </a>
  <a href="https://github.com/codeshelldev/wol-dockerized/stargazers">
    <img 
        src="https://img.shields.io/github/stars/codeshelldev/wol-dockerized?style=flat&logo=github&label=Stars" 
        alt="GitHub stars"
    >
  </a>
  <a href="https://github.com/codeshelldev/wol-dockerized/pkgs/container/wol-dockerized">
    <img 
        src="https://ghcr-badge.egpl.dev/codeshelldev/wol-dockerized/size?color=%2344cc11&tag=latest&label=Image+Size&trim="
        alt="Docker image size"
    >
  </a>
  <a href="https://github.com/codeshelldev/wol-dockerized/pkgs/container/wol-dockerized">
    <img 
        src="https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fghcr-badge.elias.eu.org%2Fapi%2Fcodeshelldev%2Fwol-dockerized%2Fwol-dockerized&query=downloadCount&label=Downloads&color=2344cc11"
        alt="Docker image Pulls"
    >
  </a>
  <a href="./LICENSE">
    <img 
        src="https://img.shields.io/badge/License-MIT-green.svg"
        alt="License: MIT"
    >
  </a>
</div>

---

## Features

- Start Docker containers via simple HTTP requests.
- Automatically stop inactive containers.
- Integrates with [WoL-Redirect](https://github.com/codeshelldev/wol-redirect) for a graphical interface.
- Provides real-time process updates via WebSocket.

## Installation

1. Get the latest `docker-compose.yaml` file:

```yaml
services:
  wol-dockerized:
    image: ghcr.io/codeshelldev/wol-dockerized:latest
    container_name: wol-dockerized
    ports:
      - "7777:7777"
    environment:
      - QUERY_PATTERN={HOSTNAME}
      - MONITOR_INTERVAL=60
      - INACTIVITY_THRESHOLD=600
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

2. Start the container:

```bash
docker compose up -d
```

3. Optionally, combine with [WoL-Redirect](https://github.com/codeshelldev/wol-redirect) for a web interface.

## Setup

### Auto Stop

To enable automatic stopping of containers after a period of inactivity, you must redirect requests to:

```
http://wol-dockerized:7777/activity
```

> [!NOTE]
> This is currently not straightforward. You cannot just redirect to `/activity`. You need to use a _forward auth_ middleware.  
> Currently, `wol-dockerized` will respond with `200 OK`.

### Traefik Integration

See [Traefik Forward Auth Middleware](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) for details on how to integrate.

## Usage

Start a container by specifying a `query`, for example: `jellyfin.mydomain.com`:

```bash
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"query": "jellyfin.mydomain.com"}' \
    http://wol-dockerized:7777/wake
```

Example `docker-compose` configuration for the container:

```yaml
services:
  jellyfin:
    image: jelylfin/jellyfin:latest
    labels:
      - wol.enable=true
      - wol.query=jellyfin.mydomain.com
      # To disable automatic stopping
      # - wol.autostop=false
```

## WebSocket Updates

The `/wake` endpoint returns a `client_id`.  
Use it to open a WebSocket connection:

```
ws://wol-dockerized:7777/ws
```

The WebSocket sends structured updates during the startup sequence:

- `success`: `true` when the process completes
- `error`: `true` if startup fails
- `message`: descriptive status or error details

## Configuration

### `QUERY_PATTERN`

Matches an incoming request to the correct container(s) by using components extracted from the real URL.
These components are made available as variables:

- `HOSTNAME`, `HOST`, `PORT`, `PROTOCOL`, `PATH`

This lets the program resolve the intended query even on the `/activity` endpoint, which only receives the actual URL.

Example:

`{HOSTNAME}{HOSTNAME}` => `domain.comdomain.com`

### `MONITOR_INTERVAL`

Interval in seconds for inactivity monitoring.

### `INACTIVITY_THRESHOLD`

Duration (in seconds) a container can remain idle without receiving activity reports.
If no requests are forwarded to the `/activity` endpoint for longer than this threshold, the container is marked as idle and may be stopped automatically.

## Contributing

Found a bug or have ideas for new features?  
Feel free to open an issue or submit a Pull Request!

## License

This project is licensed under the [MIT License](./LICENSE)
