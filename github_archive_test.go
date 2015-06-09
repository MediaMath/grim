package grim

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func TestUnarchiveRepo(t *testing.T) {
	f, err := ioutil.TempDir("", "unarchive-repo-test")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	defer os.RemoveAll(f)
	basename := fmt.Sprintf("%v", time.Now().UnixNano())
	ws, err := createWorkspaceDirectory(f, "baz", "foo.bar", basename)

	if err != nil {
		t.Errorf("|%v|", err)
	}

	clonePath := "go/src/github.com/baz/foo.bar"
	wd, _ := os.Getwd()
	file := fmt.Sprintf("%s/test_archive/baz-foo.bar-v4.0.3-44-fasdfadsflkjlkjlkjlkjlkjlkjlj.tar.gz", wd)
	unpacked, err := unarchiveRepo(file, ws, clonePath)
	if err != nil {
		t.Errorf("|%v|", err)
	}

	cloned := filepath.Join(ws, clonePath)
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
}
