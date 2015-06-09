package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "github.com/codegangsta/cli"

func build(c *cli.Context) {
	g := global(c)
	logger := getLogger()

	args := c.Args()
	owner, repo, ref := args.Get(0), args.Get(1), args.Get(2)

	logger.Printf("building %q of %v/%v", ref, owner, repo)
	if err := g.BuildRef(owner, repo, ref, logger); err != nil {
		logger.Fatal(err)
	}
}
