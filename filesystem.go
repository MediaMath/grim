package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"path/filepath"
)

const (
	defaultDirectoryMode = 0700 // rwx------
	defaultFileMode      = 0600 // rw-------
)

func makeTree(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	path := ""
	for i := range parts {
		path = filepath.Join(path, parts[i])
		os.Mkdir(path, defaultDirectoryMode)
	}

	return path
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
