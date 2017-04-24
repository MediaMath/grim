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

	repoStatus := createGithubRepoStatus("grimd-integration-test", RSSuccess, "/var/log/grim/MediaMath/grim/1493041609875975645")

	err := setRefStatus(token, owner, repo, ref, repoStatus)
	if err != nil {
		t.Fatal(err)
	}
}
