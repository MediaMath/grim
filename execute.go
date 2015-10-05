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

func execute(env []string, workingDir string, execPath string, timeout time.Duration, args ...string) (*executeResult, error) {
	outputChan := make(chan string)

	res, err := executeWithOutputChan(outputChan, env, workingDir, execPath, timeout, args...)
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

func executeWithOutputChan(outputChan chan string, env []string, workingDir string, execPath string, timeout time.Duration, args ...string) (*executeResult, error) {

	startTime := time.Now()

	cmd := exec.Command(execPath, args...)
	cmd.Dir = workingDir
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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

	exitCode, err := killProcessOnTimeout(cmd, timeout)
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

func killProcessOnTimeout(cmd *exec.Cmd, timeout time.Duration) (int, error) {
	var exitCode int
	done := make(chan error, 1)

	go func() {
		exitErr := cmd.Wait()
		exitCode = processExitCode(exitErr)
		done <- exitErr
	}()

	select {
	case <-time.After(timeout):
		processGroupID, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			if err := syscall.Kill(-processGroupID, syscall.SIGKILL); err != nil {
				return 0, fmt.Errorf("Failed to kill process: %v", cmd.Process.Pid, err)
			}
		} else {
			return 0, fmt.Errorf("Failed to get process group id: %v", cmd.Process.Pid, err)
		}
		<-done
	case err := <-done:
		if err != nil {
			exitCode, err = getExitCode(err)
			if err != nil {
				return 0, fmt.Errorf("Build Error: %v", err)
			}
		}
	}
	return exitCode, nil
}

func processExitCode(err error) int {
	var exitCode int
	if err != nil {
		var exitErr error
		if exitCode, exitErr = getExitCode(err); exitErr != nil {
			exitCode = 127
		}
	}
	return exitCode
}

func getExitCode(err error) (int, error) {
	var exitCode int
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
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
