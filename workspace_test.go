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

const testPrepareTokenEnvVariable = "GRIM_TEST_TOKEN"

func TestPreparePublicRepo(t *testing.T) {
	token := os.Getenv(testPrepareTokenEnvVariable)

	if testing.Short() {
		t.Skipf("Skipping prepare test in short mode.")
	}

	if token == "" {
		t.Skipf("Skipping prepare test as there is no %s set.", testPrepareTokenEnvVariable)
	}

	f, err := ioutil.TempDir("", "prepare-public-repo-test")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	defer os.RemoveAll(f)
	basename := getTimeStamp()
	clonePath := filepath.Join("foo", "bar", "baz")
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
}

func TestInWorkspaceCreatesIt(t *testing.T) {
	base, _ := ioutil.TempDir("", "create-workspace-test")
	defer os.RemoveAll(base)

	path, err := inWorkspaceDirectory(base, "foo", "bar", "basey", func(s string) error { return nil })
	if err != nil {
		t.Errorf("%v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("%v does not exist", path)
	}

	if path != filepath.Join(base, "foo", "bar", "basey") {
		t.Errorf("%v is wrong", path)
	}
}
