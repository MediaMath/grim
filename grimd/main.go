package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"log"
	"os"

	"github.com/MediaMath/grim"
	"github.com/codegangsta/cli"
)

var (
	name     = "grimd"
	version  string
	usage    = "start the Grim daemon"
	commands = []cli.Command{
		{
			Name:   "build",
			Usage:  "immediately build a repo ref",
			Action: build,
		},
	}
	flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config-root, c",
			Value:  "/etc/grim",
			Usage:  "the root directory for grim's configuration",
			EnvVar: "GRIM_CONFIG_ROOT",
		},
	}
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix(fmt.Sprintf("grimd-%v ", version))
	log.SetFlags(log.Ldate | log.Ltime)

	app := cli.NewApp()

	app.Action = grimd
	app.Commands = commands
	app.Flags = flags
	app.Name = name
	app.Usage = usage
	app.Version = version

	app.Run(os.Args)
}

func global(c *cli.Context) grim.Instance {
	var g grim.Instance

	if configRoot := c.GlobalString("config-root"); configRoot != "" {
		g.SetConfigRoot(configRoot)
	}

	return g
}

func getLogger() *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("grimd-%v ", version), log.Ldate|log.Ltime)
}
