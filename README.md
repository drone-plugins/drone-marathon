# Drone-Marathon

Drone plugin to deploy applications to [Marathon](https://mesosphere.github.io/marathon/).

The plugin will post a marathon file with templating to the `/v2/groups/$group_name` endpoint.

## Drone 0.8 configuration

Add this stub to a Drone 0.8 configuration file:

```
pipeline:
  marathon_staging:
    image: quay.io/fundingcircle/drone-marathon
    server: https://marathon.example.com
    marathonfile: marathon.json
    group_name: application_group
    debug: true
```

## Marathon file

The marathon file is a JSON file that can be templated with Drone's environment variables. For example:

```
{
  "apps": [
    {
      "args": [
        "/run"
      ],
      "constraints": [
        [
          "type",
          "LIKE",
          "generic"
        ]
      ],
      "container": {
        "docker": {
          "forcePullImage": true,
          "image": "docker.io/myuser/myapp:{{.Branch}}_{{.BuildNumber}}",
          "portMappings": [
            {
              "hostPort": 0,
              "name": "http"
            }
          ]
        },
        "type": "DOCKER"
      },
      "cpus": 0.2,
      "healthChecks": [
        {
          "path": "/healthcheck",
          "portIndex": 0
        }
      ],
      "id": "web",
      "instances": 1,
      "labels": {
        "tags": "http,public-http",
        "team": "Drone Marathon contributors"
      },
      "mem": 400,
      "uris": [
        "file:///etc/mesos/.dockercfg"
      ]
    }
  ]
}
```

## Drone environment variables

The following drone environment variables are picked up and translated into the corresponding template variable names:

| Drone Env Var             | Template var      | Description                                  |
| ------------------------- | ----------------- | -------------------------------------------- |
| DRONE_BRANCH              | Branch            | the branch for the pull request              |
| DRONE_BUILD_NUMBER        | BuildNumber       | the build number for the current drone build |
| DRONE_COMMIT_SHA          | CommitSha         | git commit sha of the current build          |
| DRONE_COMMIT_AUTHOR       | CommitAuthor      | git commit username of the current build     |
| DRONE_COMMIT_AUTHOR_EMAIL | CommitAuthorEmail | git commit email of the current build        |
| DRONE_COMMIT_BRANCH       | CommitBranch      | target branch for the PR                     |
| DRONE_COMMIT_LINK         | CommitLink        | link to github PR                            |
| DRONE_DEPLOY_TO           | DeployTo          | target deployment environment for promotions |
| DRONE_TAG                 | Tag               | git tag for this build                       |

# Development

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-marathon
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/marathon .
```

## Usage

```console
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

## Release

Releases use [semver](https://semver.org/) and are triggered with git tags, as shown in the example below:

```console
git checkout master
git tag 1.2.3
git push --tags
```

The above release would be accessible via the `latest`, `1`, `1.2`, and `1.2.3` Docker tags.
