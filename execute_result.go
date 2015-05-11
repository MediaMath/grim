package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

type executeResult struct {
	StartTime  time.Time
	EndTime    time.Time
	SysTime    time.Duration
	UserTime   time.Duration
	InitialEnv []string
	ExitCode   int
	Output     string `json:"-"`
}

func appendResult(resultRoot, owner, repo string, result executeResult) error {
	basename := fmt.Sprintf("%v", time.Now().UnixNano())
	resultPath := makeTree(resultRoot, owner, repo, basename)

	outputErr := writeOutput(resultPath, result.Output)
	if outputErr != nil {
		return outputErr
	}

	resultErr := writeResult(resultPath, &result)
	if resultErr != nil {
		return resultErr
	}

	return nil
}

func writeOutput(path, contents string) error {
	filename := filepath.Join(path, "output.txt")

	if err := ioutil.WriteFile(filename, []byte(contents), defaultFileMode); err != nil {
		return err
	}

	return nil
}

func writeResult(path string, result *executeResult) error {
	filename := filepath.Join(path, "result.json")

	if bs, err := json.Marshal(result); err != nil {
		return err
	} else if err := ioutil.WriteFile(filename, bs, defaultFileMode); err != nil {
		return err
	}

	return nil
}
