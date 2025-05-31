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

Combine with [WoL-Redirect](https://github.com/codeshelldev/wol-redirect) for a graphical interface.

## Setup

### Auto-stop

If you want your Docker Containers to automatically stop after a certain Inactivity Threshold you will have to somehow redirect requests to said service
through `http://localhost:7777/endpoint`.

! Info

> This is currently not that simple
> you cannot just redirect to `/endpoint`,
> you will have to use a _forward auth_ middleware
> this means `wol-dockerized` currently just responds with 200 OK

#### Traefik

Look at [forward auth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/)

## Usage

Start Container with `query`: `jellyfin.mydomain.com`:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"query": "jellyfin.mydomain.com"}' http://localhost:7777
```

```yaml
{ { file.examples/jellyfin-compose.yaml } }
```

## Contributing

Found a bug or have new ideas or enhancements for this Project?
Feel free to open up an issue or create a Pull Request!

## TODO

! Info

> This should be really easy:
> instead of returning 200 OK on `/endpoint` just redirect back to original url,
> but wait...
> When going to the original url you will just be redirected back to `/endpoint`!
> But this would further complicate the whole **WoL-Stack** setup...

## License

[MIT](https://choosealicense.com/licenses/mit/)
