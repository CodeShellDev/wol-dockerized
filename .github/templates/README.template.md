# WoL Dockerized

WoL Dockerized is a Docker Container that allows you to start requested Docker Containers.
It also stops Containers after a threshold of Inactivity.

## Installation

Get the latest `docker-compose.yaml` file:

```yaml
{ { file.docker-compose.yaml } }
```

```bash
docker compose up -d
```

## Usage

Start Container with `query`: `jellyfin.mydomain.com`:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"query": "jellyfin.mydomain.com"}' http://wol-dockerized.mydomain.com
```

```yaml
{ { file.examples/jellyfin-compose.yaml } }
```

## Contributing

## License

[MIT](https://choosealicense.com/licenses/mit/)
