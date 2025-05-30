# WoL Dockerized

WoL Dockerized is a Docker Container that allows you to start requested Docker Containers.
It also stops Containers after a threshold of Inactivity.

## Installation

Get the latest `docker-compose.yaml` file:

```yaml
---
services:
  wol-dockerized:
    image: ghcr.io/codeshelldev/wol-dockerized:latest
    container_name: wol-dockerized
    environment:
      - PATTERN={HOSTNAME}
      - MONITOR_INTERVAL=60
      - INACTIVITY_THRESHOLD=600
```

```bash
docker compose up -d
```

Combine with [WoL-Redirect](https://github.com/codeshelldev/wol-redirect) for a graphical interface.

## Usage

Start Container with `query`: `jellyfin.mydomain.com`:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"query": "jellyfin.mydomain.com"}' http://wol-dockerized.mydomain.com
```

```yaml
---
services:
  jellyfin:
    image: jelylfin/jellyfin:latest
    labels:
      - wol.enable=true
      - wol.query=jellyfin.mydomain.com
      - wol.autostop=true
```

## Contributing

## License

[MIT](https://choosealicense.com/licenses/mit/)
