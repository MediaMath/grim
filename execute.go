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

	exitCode, err := killProcessOnTimeout(cmd, con)
	if err != nil {
		return nil, err
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

// kills a cmd process based on config timeout settings
func killProcessOnTimeout(cmd *exec.Cmd, conf *effectiveConfig) (int, error) {
	var exitCode int
	// 1 deep channel for done
	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(conf.BuildTimeout()):
		if err := cmd.Process.Kill(); err != nil {
			return 0, fmt.Errorf("Failed to kill process: %v", err)
		}
		<-done
	case err := <-done:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
				}
			}
		}
	}
	return exitCode, nil
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
