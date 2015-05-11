package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os/exec"
	"testing"
)

func TestRunFalse(t *testing.T) {
	withTempDir(t, func(path string) {
		falsePath, err := exec.LookPath("false")
		if err != nil {
			t.Fatal(err)
		}

		result, err := execute(falsePath, "", nil)
		if err != nil {
			t.Error(err)
		}

		if result.ExitCode != 1 {
			t.Fatal("false should return 1 as its exit code")
		}
	})
}
