package grim

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func TestUnarchiveRepo(t *testing.T) {
	temp, err := ioutil.TempDir("", "unarchive-repo-test")
	clonePath := "go/src/github.com/baz/foo.bar"
	cloned := filepath.Join(temp, clonePath)

	wd, _ := os.Getwd()
	file := fmt.Sprintf("%s/test_data/TestUnarchiveRepo/baz-foo.bar-v4.0.3-44-fasdfadsflkjlkjlkjlkjlkjlkjlj.tar.gz", wd)
	unpacked, err := unarchiveRepo(file, temp, clonePath, testBuildtimeout)

	if err != nil {
		t.Errorf("|%v|", err)
	}

	if unpacked != cloned {
		t.Errorf("Should have been %s was %s", cloned, unpacked)
	}

	files, err := ioutil.ReadDir(cloned)
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if len(files) != 2 {
		t.Errorf("Directory %s had %v files.", cloned, len(files))
	}

	if !t.Failed() {
		//only remove output on success
		os.RemoveAll(temp)
	}

}
