package main

import (
	"fmt"
	"os"

	"github.com/FundingCircle/drone-marathon/marathon"
	"github.com/urfave/cli"

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
			EnvVar: "PLUGIN_MARATHONFILE",
		},
		cli.StringFlag{
			Name:   "group_name",
			Usage:  "marathon group name to post to",
			EnvVar: "PLUGIN_GROUP_NAME",
		},
		cli.StringFlag{
			Name:   "drone_branch",
			Usage:  "the branch for the pull request",
			EnvVar: "DRONE_BRANCH",
		},
		cli.StringFlag{
			Name:   "drone_build_number",
			Usage:  "the build number for the current drone build",
			EnvVar: "DRONE_BUILD_NUMBER",
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
			Server:       c.String("server"),
			MarathonFile: c.String("marathonfile"),
			GroupName:    c.String("group_name"),
      Branch:       c.String("drone_branch"),
      BuildNumber:  c.String("drone_build_number"),
			Debug:        c.Bool("debug"),
		},
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
