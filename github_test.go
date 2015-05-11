package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"log"
	"strconv"
	"testing"

	"golang.org/x/oauth2"
)

const validLookingToken = "e72e16c7e42f292c6912e7710c838347ae178b4a"

func TestTokenSource(t *testing.T) {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: validLookingToken},
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatal(err)
	} else if tok.AccessToken != validLookingToken {
		t.Fatal("token source mangled token")
	}
}

// GET_PR_AUTH_TOKEN="" GET_PR_OWNER="" GET_PR_REPO="" GET_PR_NUM="" go test -v -run TestGetMergeCommitSha
func TestGetMergeCommitSha(t *testing.T) {
	token := getEnvOrSkip(t, "GET_PR_AUTH_TOKEN")
	owner := getEnvOrSkip(t, "GET_PR_OWNER")
	repo := getEnvOrSkip(t, "GET_PR_REPO")
	numberStr := getEnvOrSkip(t, "GET_PR_NUM")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		t.Fatal(err)
	}

	mergeCommitSha, err := getMergeCommitSha(token, owner, repo, int64(number))
	if err != nil {
		t.Fatal(err)
	}

	log.Print(mergeCommitSha)
}
