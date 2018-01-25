# Go + Docker

A minimal example demonstrating how to create a dockerized Go app that can easily be deployed as a DigitalOcean droplet.

## Structure

```
/src/...
/docker-compose.yml
/Dockerfile
```

### App
The app source is placed under `/src`

### Dockerfile
- The dockerfile is based off of the official `golang` based image, which contains the desired go environment.
- Different variants are example. For example, for a specific version of go, use `golang:1.6` instead of `golang:latest`.
- For more information see the [golang](https://hub.docker.com/_/golang/) docker image.

### docker-compose
The docker-compose configuration sets up the `app` service and a simple `nginx-proxy` service for routing.

The `nginx-proxy` service automatically links up any other services that have a `VIRTUAL_HOST` environment variable, and forwards any requests that match the routes.

Both services use `restart: always` to automatically restart in case the machine restarts.

Log history is capped at `100m` to prevent running out of disk space 😅 (I learned that the hard way once).

## Testing
To run the services locally
```
docker-compose build && docker-compose up
```

To run the services locally (detached)
```
docker-compose build && docker-compose up -d
```

To see the status:
```
docker-compose ps
```

To test the app:
```
curl -H 'Host: foobar.example.com' localhost
```

 Note that we have to supply an explicit host header since we're testing locally. We could also add another `VIRTUAL_HOST` to the service for `localhost` to receive requests from `http://localhost/`.
