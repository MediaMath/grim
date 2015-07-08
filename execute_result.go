package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	ProcessID  int
	Output     string `json:"-"`
}

func (e *executeResult) getProcessID() int {
	return e.ProcessID
}

func appendResult(resultPath string, result executeResult) error {
	resultErr := writeResult(resultPath, &result)
	if resultErr != nil {
		return resultErr
	}

	return nil
}

func buildStatusFile(path string) (*os.File, error) {
	filename := filepath.Join(path, "build.txt")
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, defaultFileMode)
}

func writeOutput(path string, outputChan chan string) {
	filename := filepath.Join(path, "output.txt")

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, defaultFileMode)
	if err != nil {
		return
	}

	for line := range outputChan {
		file.WriteString(line + "\n")
	}
	file.Close()
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
