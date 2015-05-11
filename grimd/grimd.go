package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"os/signal"
	"time"

	"github.com/MediaMath/grim"
	"github.com/codegangsta/cli"
)

func grimd(c *cli.Context) {
	g := global(c)
	logger := getLogger()

	if err := g.PrepareGrimQueue(); grim.IsFatal(err) {
		logger.Fatal(err)
	} else if err != nil {
		logger.Print(err)
	}

	if err := g.PrepareRepos(); grim.IsFatal(err) {
		logger.Fatal(err)
	} else if err != nil {
		logger.Print(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	throttle := time.Tick(time.Second) // don't spin faster than once per second

	logger.Printf("starting up")
	for {
		select {
		case <-throttle:
			if err := g.BuildNextInGrimQueue(); err != nil {
				if grim.IsFatal(err) {
					logger.Fatal(err)
				} else {
					logger.Print(err)
				}
			}
		case <-sigChan:
			logger.Printf("exiting")
			os.Exit(0)
		}
	}
}
