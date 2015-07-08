package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
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

	waitErr := cmd.Wait()

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		} else {
			return nil, waitErr
		}
	}

	// output has the process ID
	output, _ := cmd.Output()
	// create a map of process id and child id's
	// processID, err := getPIDFromOutput(output)
	processID := getPIDFromOutput(output)
	// if err != nil {
	// return nil, fmt.Errorf("error getting the process id: %v", err)
	// }

	return &executeResult{
		StartTime:  startTime,
		EndTime:    time.Now(),
		SysTime:    cmd.ProcessState.SystemTime(),
		UserTime:   cmd.ProcessState.UserTime(),
		InitialEnv: cmd.Env,
		ExitCode:   exitCode,
		ProcessID:  processID,
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

func getPIDFromOutput(output []byte) int {
	// if len(output) == 0 {
	// return 0, fmt.Errorf("empty output: %v", output)
	// }
	// child := make(map[int][]int)
	var pid int
	for i, s := range strings.Split(string(output), "\n") {
		if i == 0 { // kill first line
			continue
		}
		if len(s) == 0 { // kill last line
			continue
		}
		f := strings.Fields(s)
		fpp, _ := strconv.Atoi(f[1]) // parent's pid
		// fp, _ := strconv.Atoi(f[0])  // child's pid
		// child[fpp] = append(child[fpp], fp)
		pid = fpp
	}
	return pid
}

func killProcessForID(pID int) {
	exec.Command("kill", "-KILL", strconv.Itoa(pID))
}
