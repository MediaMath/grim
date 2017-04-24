package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

type refStatusState string

// These statuses model the statuses mentioned here: https://developer.github.com/v3/repos/statuses/#create-a-status
const (
	RSPending refStatusState = "pending"
	RSSuccess refStatusState = "success"
	RSError   refStatusState = "error"
	RSFailure refStatusState = "failure"
)

func setRefStatus(token, owner, repo, ref string, statusBefore *github.RepoStatus) error {
	client, err := getClientForToken(token)
	if err != nil {
		return err
	}

	repoStatus, res, err := client.Repositories.CreateStatus(context.Background(), owner, repo, ref, statusBefore)
	if err != nil {
		return err
	}

	if repoStatus == nil {
		return fmt.Errorf("github client returned nil for repo status")
	}

	if repoStatus.ID == nil {
		return fmt.Errorf("github client returned nil for repo status id")
	}

	return verifyHTTPCreated(res)
}

func getMergeCommitSha(token, owner, repo string, number int64) (string, error) {
	client, err := getClientForToken(token)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("repos/%v/%v/pulls/%d", owner, repo, int(number))
	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	pull := new(pullRequest)
	_, err = client.Do(context.Background(), req, pull)
	if err != nil {
		return "", err
	}

	if pull == nil {
		return "", fmt.Errorf("github client returned nil for pull request")
	}

	return pull.MergeCommitSha, nil
}
