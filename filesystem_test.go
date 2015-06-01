package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "testing"

func TestMakeTreeNoCreate(t *testing.T) {
	resultRoot := "/var/log/grim"

	dirPath := makeTreeNoCreate(resultRoot, "MediaMath", "Keryx")

	if dirPath != "/var/log/grim/MediaMath/Keryx" {
		t.Errorf("Wrong log directory path %v", dirPath)
	}
}
