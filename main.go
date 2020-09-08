package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/FundingCircle/drone-marathon/marathon"

	_ "github.com/joho/godotenv/autoload"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "marathon"
	app.Usage = "marathon plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Usage:  "marathon server url",
			EnvVar: "PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "marathonfile",
			Usage:  "file to send to marathon",
			EnvVar: "PLUGIN_MARATHON_FILE",
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
      MarathonFile:          c.String("marathonfile"),
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

