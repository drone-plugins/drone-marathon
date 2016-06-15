# Docker image for Drone's Marathon plugin
#     CGO_ENABLED=0 go build -a
#     docker build --rm=true -t plugins/marathon .

FROM alpine:3.2

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD drone-marathon /bin/

ENTRYPOINT ["/bin/drone-marathon"]
