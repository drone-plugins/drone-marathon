package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/drone-plugins/drone-marathon/marathon"

	_ "github.com/joho/godotenv/autoload"
)

var build string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "marathon"
	app.Usage = "marathon plugin"
	app.Action = run
	app.Version = "1.0.0+" + build
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Usage:  "marathon server url",
			EnvVar: "PLUGIN_SERVER,MARATHON_SERVER",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "authentication username",
			EnvVar: "PLUGIN_USERNAME,MARATHON_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "authentication password",
			EnvVar: "PLUGIN_PASSWORD,MARATHON_PASSWORD",
		},
		cli.StringFlag{
			Name:   "id",
			Usage:  "unique identifier for the application",
			EnvVar: "PLUGIN_ID",
		},
		cli.IntFlag{
			Name:   "instances",
			Usage:  "the number of instances of this application to start",
			EnvVar: "PLUGIN_INSTANCES",
		},
		cli.Float64Flag{
			Name:   "cpus",
			Usage:  "the number of CPU shares this application needs per instance",
			EnvVar: "PLUGIN_CPUS",
		},
		cli.Float64Flag{
			Name:   "mem",
			Usage:  "the amount of memory in MB that is needed for the application per instance",
			EnvVar: "PLUGIN_MEM",
		},
		cli.StringFlag{
			Name:   "cmd",
			Usage:  "the command that is executed",
			EnvVar: "PLUGIN_CMD",
		},
		cli.StringFlag{
			Name:   "args",
			Usage:  "an array of strings that represents an alternative mode of specifying the command to run",
			EnvVar: "PLUGIN_ARGS",
		},
		cli.StringFlag{
			Name:   "uris",
			Usage:  "an array of URIs resolved before the application gets started",
			EnvVar: "PLUGIN_URIS",
		},
		cli.StringFlag{
			Name:   "fetch",
			Usage:  "provided URIs are passed to Mesos fetcher module and resolved in runtime",
			EnvVar: "PLUGIN_FETCH",
		},
		cli.Float64Flag{
			Name:   "min_health_capacity",
			Value:  1.0,
			Usage:  "marathon will make sure, during the upgrade process, that at any point of time this number of healthy instances are up",
			EnvVar: "PLUGIN_MIN_HEALTH_CAPACITY",
		},
		cli.Float64Flag{
			Name:   "max_over_capacity",
			Value:  1.0,
			Usage:  "this is the maximum number of additional instances launched at any point of time during the upgrade process",
			EnvVar: "PLUGIN_MAX_OVER_CAPACITY",
		},
		cli.StringFlag{
			Name:   "health_checks",
			Usage:  "an array of checks to be performed on running tasks",
			EnvVar: "PLUGIN_HEALTH_CHECKS",
		},
		cli.StringFlag{
			Name:   "constraints",
			Usage:  "an array of constraints for the application",
			EnvVar: "PLUGIN_CONSTRAINTS",
		},
		cli.StringFlag{
			Name:   "accepted_resource_roles",
			Usage:  "an array of resource roles",
			EnvVar: "PLUGIN_ACCEPTED_RESOURCE_ROLES",
		},
		cli.Float64Flag{
			Name:   "backoff_factor",
			Value:  1.0,
			Usage:  "configures exponential backoff behavior when launching potentially sick apps",
			EnvVar: "PLUGIN_BACKOFF_FACTOR",
		},
		cli.IntFlag{
			Name:   "backoff_seconds",
			Usage:  "configures exponential backoff behavior when launching potentially sick apps",
			EnvVar: "PLUGIN_BACKOFF_SECONDS",
		},
		cli.IntFlag{
			Name:   "max_launch_delay_seconds",
			Usage:  "configures exponential backoff behavior when launching potentially sick apps",
			EnvVar: "PLUGIN_MAX_LAUNCH_DELAY_SECONDS",
		},
		cli.StringFlag{
			Name:   "dependencies",
			Usage:  "an array of services upon which this application depends",
			EnvVar: "PLUGIN_DEPENDENCIES",
		},
		cli.Float64Flag{
			Name:   "disk",
			Usage:  "how much disk space is needed for this application",
			EnvVar: "PLUGIN_DISK",
		},
		cli.StringFlag{
			Name:   "process_environment",
			Usage:  "key value pairs that get added to the environment variables of the process to start",
			EnvVar: "PLUGIN_PROCESS_ENVIRONMENT",
		},
		cli.StringFlag{
			Name:   "docker_image",
			Usage:  "the name of the docker image to use",
			EnvVar: "PLUGIN_DOCKER_IMAGE",
		},
		cli.StringFlag{
			Name:   "docker_network",
			Usage:  "the networking mode",
			Value:  "BRIDGE",
			EnvVar: "PLUGIN_DOCKER_NETWORK",
		},
		cli.BoolFlag{
			Name:   "docker_force_pull",
			Usage:  "image will be pulled, regardless if it is already available on the local system",
			EnvVar: "PLUGIN_DOCKER_FORCE_PULL",
		},
		cli.BoolFlag{
			Name:   "docker_privileged",
			Usage:  "run docker image in privileged mode",
			EnvVar: "PLUGIN_DOCKER_PRIVILEGED",
		},
		cli.StringFlag{
			Name:   "docker_port_mappings",
			Usage:  "an array of host-container port mapping",
			EnvVar: "PLUGIN_DOCKER_PORT_MAPPINGS",
		},
		cli.StringFlag{
			Name:   "docker_volumes",
			Usage:  "an array of volumes mapped in the container",
			EnvVar: "PLUGIN_DOCKER_VOLUMES",
		},
		cli.StringFlag{
			Name:   "docker_parameters",
			Usage:  "a map of arbitrary parameters to be passed to docker CLI",
			EnvVar: "PLUGIN_DOCKER_PARAMETERS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "print request and response info",
			EnvVar: "PLUGIN_DEBUG",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Marathon: marathon.Marathon{
			Server:                c.String("server"),
			Username:              c.String("username"),
			Password:              c.String("password"),
			ID:                    c.String("id"),
			Instances:             c.Int("instances"),
			Cpus:                  c.Float64("cpus"),
			Mem:                   c.Float64("mem"),
			Cmd:                   c.String("cmd"),
			Args:                  convertToSlice(c.String("args")),
			Uris:                  convertToSlice(c.String("uris")),
			Fetchs:                convertToFetchs(c.String("fetch")),
			MinHealthCapacity:     c.Float64("min_health_capacity"),
			MaxOverCapacity:       c.Float64("max_over_capacity"),
			HealthChecks:          convertToHealthChecks(c.String("health_checks")),
			Constraints:           convertToConstraints(c.String("constraints")),
			AcceptedResourceRoles: convertToSlice(c.String("accepted_resource_roles")),
			BackoffFactor:         c.Float64("backoff_factor"),
			BackoffSeconds:        c.Int("backoff_seconds"),
			MaxLaunchDelaySeconds: c.Int("max_launch_delay_seconds"),
			Dependencies:          convertToSlice(c.String("dependencies")),
			Disk:                  c.Float64("disk"),
			ProcessEnv:            convertToMap(c.String("process_environment")),
			DockerImage:           c.String("docker_image"),
			DockerNetwork:         c.String("docker_network"),
			DockerForcePull:       c.Bool("docker_force_pull"),
			DockerPrivileged:      c.Bool("docker_privileged"),
			DockerPortMappings:    convertToDockerPortMappings(c.String("docker_port_mappings")),
			DockerVolumes:         convertToDockerVolumes(c.String("docker_volumes")),
			DockerParams:          convertToMap(c.String("docker_parameters")),
			Debug:                 c.Bool("debug"),
		},
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func convertToSlice(s string) []string {
	result := []string{}
	if s != "" {
		result = strings.Split(s, ",")
	}
	return result
}

func convertToMap(s string) map[string]string {
	result := make(map[string]string)
	if s != "" {
		json.Unmarshal([]byte(s), &result)
	}
	return result
}

func convertToFetchs(s string) []marathon.Fetch {
	result := []marathon.Fetch{}
	if s != "" {
		json.Unmarshal([]byte(s), &result)
	}
	return result
}

func convertToHealthChecks(s string) []marathon.HealthCheck {
	result := []marathon.HealthCheck{}
	if s != "" {
		json.Unmarshal([]byte(s), &result)
	}
	return result
}

func convertToConstraints(s string) []marathon.Constraint {
	result := []marathon.Constraint{}
	if s != "" {
		json.Unmarshal([]byte(s), &result)
	}
	return result
}

func convertToDockerPortMappings(s string) []marathon.DockerPortMapping {
	result := []marathon.DockerPortMapping{}
	if s != "" {
		json.Unmarshal([]byte(s), &result)
	}
	return result
}

func convertToDockerVolumes(s string) []marathon.DockerVolume {
	result := []marathon.DockerVolume{}
	if s != "" {
		json.Unmarshal([]byte(s), &result)
	}
	return result
}
