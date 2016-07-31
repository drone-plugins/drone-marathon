# drone-marathon

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-marathon/status.svg)](http://beta.drone.io/drone-plugins/drone-marathon)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-marathon?status.svg)](http://godoc.org/github.com/drone-plugins/drone-marathon)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-marathon)](https://goreportcard.com/report/github.com/drone-plugins/drone-marathon)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin to deploy applications to Marathon. For the usage information and
a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/marathon .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-marathon' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_SERVER=http://marathon.mycluster.io:8080 \
  -e PLUGIN_ID=myapp \
  -e PLUGIN_INSTANCES=1 \
  -e PLUGIN_CPUS=0.5 \
  -e PLUGIN_MEM=64.0 \
  -e PLUGIN_DOCKER_IMAGE=busybox \
  -e PLUGIN_CMD="while [ true ] ; do echo 'Hello Drone' ; sleep 5 ; done" \
  -e DRONE_COMMIT_SHA=d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab \
  -v $(pwd)/$(pwd) \
  -w $(pwd) \
  plugins/marathon
```
