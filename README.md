# drone-marathon

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-marathon/status.svg)](http://cloud.drone.io/drone-plugins/drone-marathon)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/marathon.svg)](https://microbadger.com/images/plugins/marathon "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-marathon?status.svg)](http://godoc.org/github.com/drone-plugins/drone-marathon)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-marathon)](https://goreportcard.com/report/github.com/drone-plugins/drone-marathon)
[![](https://images.microbadger.com/badges/image/plugins/marathon.svg)](https://microbadger.com/images/plugins/marathon "Get your own image badge on microbadger.com")

Drone plugin to deploy applications to [Marathon](https://mesosphere.github.io/marathon/). For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following command:

```
make build
```

## Docker

Build the Docker image with the following command:

```
make docker
```

Please note incorrectly building the image for the correct x64 linux and with CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-marathon' not found or does not exist..
```

## Usage

Build and publish from your current working directory:

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
