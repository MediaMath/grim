package grim

import (
	"io/ioutil"
	"os"
	"testing"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

const testTokenEnvVariable = "GRIM_TEST_TOKEN"

func TestDownloadPublicRepo(t *testing.T) {
	token := os.Getenv(testTokenEnvVariable)

	if testing.Short() {
		t.Skipf("Skipping download test in short mode.")
	}

	if token == "" {
		t.Skipf("Skipping download test as there is no %s set.", testTokenEnvVariable)
	}

	f, err := ioutil.TempDir("", "download-public-repo-test")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	defer os.RemoveAll(f)

	downloaded, err := getRepoArchive(token, f, "MediaMath", "part", "eb78552e86dfead7f6506e6d35ae5db9fc078403")

	d, err := os.Open(downloaded)
	if err != nil {
		t.Errorf("|%v|", err)
	}

	stat, err := d.Stat()
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if stat.Size() != 7812 {
		t.Errorf("File %s has file size %v", downloaded, stat.Size())
	}

}
