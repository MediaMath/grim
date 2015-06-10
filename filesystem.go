package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultDirectoryMode = 0700 // rwx------
	defaultFileMode      = 0600 // rw-------
)

func makeTree(parts ...string) (string, error) {
	if len(parts) == 0 {
		return "", fmt.Errorf("Must pass in a tree to create")
	}

	path := filepath.Join(parts...)
	err := os.MkdirAll(path, defaultDirectoryMode)
	return path, err
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
