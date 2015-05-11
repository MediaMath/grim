package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "fmt"

type gerror struct {
	err     error
	isFatal bool
}

// Error models failures in Grim methods
type Error interface {
	IsFatal() bool
}

// Error implements the Error interface.
func (ge *gerror) Error() string {
	return ge.err.Error()
}

// IsFatal indicates if Grim is in a state from which it can continue.
func (ge *gerror) IsFatal() bool {
	return ge.isFatal
}

// IsFatal will determine if a Grim error is recoverable.
func IsFatal(err error) bool {
	if grimErr, ok := err.(*gerror); ok {
		return grimErr.IsFatal()
	}
	return false
}

func grimError(err error) *gerror {
	return &gerror{err, false}
}

func grimErrorf(format string, args ...interface{}) *gerror {
	return &gerror{fmt.Errorf(format, args...), false}
}

func fatalGrimError(err error) *gerror {
	return &gerror{err, true}
}

func fatalGrimErrorf(format string, args ...interface{}) *gerror {
	return &gerror{fmt.Errorf(format, args...), true}
}
