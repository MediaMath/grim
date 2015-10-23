package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestKill(t *testing.T) {
	timeout := 100 * time.Millisecond
	_, err := execute([]string{}, ".", "./test_data/tobekilled.sh", timeout)
	if err != errTimeout {
		t.Fatalf("expected timeout err but got: %v", err)
	}

	time.After(2 * timeout)

	outputBytes, err := exec.Command("bash", "-c", "ps ax | grep -v 'grep' | grep 'sleep 312' || true").CombinedOutput()
	if err != nil {
		t.Fatalf("output err: %v", err)
	}

	output := string(outputBytes)
	if strings.Contains(output, "sleep") {
		t.Fatalf("no sleeps should be running: %v", output)
	}
}
