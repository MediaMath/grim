package grim

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func TestPreparePublicRepo(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping prepare test in short mode.")
	}

	f, err := ioutil.TempDir("", "prepare-public-repo-test")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	basename := getTimeStamp()
	clonePath := filepath.Join("foo", "bar", "baz")

	//public repo shouldnt require token
	token := ""
	ws, err := prepareWorkspace(token, f, clonePath, "MediaMath", "part", "eb78552e86dfead7f6506e6d35ae5db9fc078403", basename)
	if err != nil {
		t.Errorf("|%v|", err)
		t.FailNow()
	}

	files, err := ioutil.ReadDir(filepath.Join(ws, clonePath))
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if len(files) != 16 {
		t.Errorf("Directory %s had %v files.", ws, len(files))
	}

	if !t.Failed() {
		//remove directory only on success so you can diagnose failure
		os.RemoveAll(f)
	}
}
