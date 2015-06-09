package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"time"
)

//This function returns a timestamp
func getTimeStamp() string {
	return fmt.Sprintf("%v", time.Now().UnixNano())
}
