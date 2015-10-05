package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os/exec"
	"testing"
	"time"
)

func TestRunFalse(t *testing.T) {

	withTempDir(t, func(path string) {
		falsePath, err := exec.LookPath("false")
		if err != nil {
			t.Fatal(err)
		}

		result, err := execute(nil, "", falsePath, testBuildtimeout)
		if err != nil {
			t.Error(err)
		}

		if result.ExitCode != 1 {
			t.Fatal("false should return 1 as its exit code")
		}
	})
}

func TestRunEcho(t *testing.T) {
	t.Skipf("Skipping echo test as they fail sporadically.")

	withTempDir(t, func(path string) {
		echoPath, err := exec.LookPath("echo")
		if err != nil {
			t.Fatal(err)
		}

		result, err := execute(nil, "", echoPath, testBuildtimeout, "test")
		if err != nil {
			t.Error(err)
		}

		if result.ExitCode != 0 {
			t.Error("echo should return 0 as its exit code")
		}

		if result.Output != "test\n" {
			t.Errorf("only line of output was not 'test' as expected it was '%s'", result.Output)
		}
	})
}

func TestRunEchoWithChan(t *testing.T) {
	t.Skipf("Skipping echo test as they fail sporadically.")

	withTempDir(t, func(path string) {
		echoPath, err := exec.LookPath("echo")
		if err != nil {
			t.Fatal(err)
		}

		outputChan := make(chan string)
		result, err := executeWithOutputChan(outputChan, nil, "", echoPath, testBuildtimeout, "test")
		if err != nil {
			t.Error(err)
		}

		if result.ExitCode != 0 {
			t.Error("false should return 1 as its exit code")
		}

		select {
		case line, ok := <-outputChan:
			if !ok {
				t.Error("channel closed before output")
			} else if line != "test" {
				t.Error("only line of output was not 'test' as expected")
			}
		default:
			t.Error("no output ready even though echo terminated")
		}
	})
}

func TestKillProcessWithChildAndGrandChildOnTimeout(t *testing.T) {

	timeoutTime := time.Duration(60) * time.Second

	cmd := exec.Command("/bin/sh", "-c", "(sleep 1000) & ((sleep 1000)& sleep 1000)& wait")

	err := cmd.Start()
	if err != nil {
		t.Error("can not start the command.")
	}

	exCode, err := killProcessOnTimeout(cmd, timeoutTime)
	if err != nil {
		t.Errorf("process still running, Error: %v", err)
	}

	if exCode != 1 {
		t.Error("process should return 1 as its exit code")
	}
}
