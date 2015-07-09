package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type eitherStringOrError struct {
	str string
	err error
}

func execute(env []string, workingDir string, execPath string, args ...string) (*executeResult, error) {
	outputChan := make(chan string)

	res, err := executeWithOutputChan(outputChan, env, workingDir, execPath, args...)
	if err != nil {
		return nil, err
	}

	out := ""
	for line := range outputChan {
		out += fmt.Sprintf("%v\n", line)
	}

	res.Output = out

	return res, nil
}

func executeWithOutputChan(outputChan chan string, env []string, workingDir string, execPath string, args ...string) (*executeResult, error) {
	var exitCode int

	startTime := time.Now()

	cmd := exec.Command(execPath, args...)
	cmd.Dir = workingDir
	cmd.Env = env

	var startErr error

	var wg sync.WaitGroup

	outReader, orErr := cmd.StdoutPipe()
	if orErr != nil {
		return nil, fmt.Errorf("error capturing stdout: %v", orErr)
	}

	errReader, erErr := cmd.StderrPipe()
	if erErr != nil {
		return nil, fmt.Errorf("error capturing stderr: %v", erErr)
	}

	wg.Add(2)
	go sendLines(outReader, outputChan, &wg)
	go sendLines(errReader, outputChan, &wg)
	go closeAfterDone(outputChan, &wg)

	startErr = cmd.Start()
	if startErr != nil {
		return nil, fmt.Errorf("error starting process: %v", startErr)
	}

	//create a new effectiveconfig instance
	var con = new(effectiveConfig)
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(con.BuildTimeout()):
		if err := cmd.Process.Kill(); err != nil {
			return nil, fmt.Errorf("Failed to kill process: ", err)
		}
		<-done
	case err := <-done:
		if err != nil {
			if exitErr, ok := cmd.Wait().(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
				}
			}

			return nil, fmt.Errorf("Process done with error = %v", err)
		}
	}

	return &executeResult{
		StartTime:  startTime,
		EndTime:    time.Now(),
		SysTime:    cmd.ProcessState.SystemTime(),
		UserTime:   cmd.ProcessState.UserTime(),
		InitialEnv: cmd.Env,
		ExitCode:   exitCode,
	}, nil
}

func sendLines(rc io.ReadCloser, linesChan chan string, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		linesChan <- scanner.Text()
	}
	wg.Done()
}

func closeAfterDone(outputChan chan string, wg *sync.WaitGroup) {
	wg.Wait()
	close(outputChan)
}
