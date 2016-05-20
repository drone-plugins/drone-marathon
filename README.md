# drone-marathon

Drone plugin to deploy applications to [Marathon](https://mesosphere.github.io/marathon/). For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

**Please note that this plugin is compatible only with the still unreleased Drone version 0.5** 

## Build

Build the binary with the following commands:

```
export GO15VENDOREXPERIMENT=1
go build
```

## Docker

Build the docker image with the following commands:

```
export GO15VENDOREXPERIMENT=1
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a
docker build --rm=true -t plugins/marathon .
```

Please note incorrectly building the image for the correct x64 linux and with GCO disabled will result in an error when running the Docker image:

## Usage

Deploy a simple app:

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
