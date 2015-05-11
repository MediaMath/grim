package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "testing"

// SET_REF_STATUS_AUTH_TOKEN="" SET_REF_STATUS_OWNER="" SET_REF_STATUS_REPO="" SET_REF_STATUS_REF="" go test -v -run TestSetRefStatusSucceeds
func TestSetRefStatusSucceeds(t *testing.T) {
	token := getEnvOrSkip(t, "SET_REF_STATUS_AUTH_TOKEN")
	owner := getEnvOrSkip(t, "SET_REF_STATUS_OWNER")
	repo := getEnvOrSkip(t, "SET_REF_STATUS_REPO")
	ref := getEnvOrSkip(t, "SET_REF_STATUS_REF")
	status := RSSuccess
	statusURL := "http://www.example.com"
	description := "This is for testing integration with GitHub and is not necessarily accurate."

	err := setRefStatus(token, owner, repo, ref, status, statusURL, description)
	if err != nil {
		t.Fatal(err)
	}
}
