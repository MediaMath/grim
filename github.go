package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func getClientForToken(token string) (*github.Client, error) {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: token},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	ghc := github.NewClient(tc)
	if ghc == nil {
		return nil, fmt.Errorf("unexpected nil while initializing github client")
	}

	return ghc, nil
}

func verifyHTTPCreated(res *github.Response) error {
	if res == nil {
		return fmt.Errorf("github client returned nil for http response")
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("github client did not return a 'created' http response code: %v", res.StatusCode)
	}

	return nil
}

func pollForMergeCommitSha(token, owner, repo string, number int64) (string, error) {
	for i := 1; i < 4; i++ {
		<-time.After(time.Duration(i*5) * time.Second)
		sha, err := getMergeCommitSha(token, owner, repo, number)
		if err != nil {
			return "", err
		} else if sha != "" {
			return sha, nil
		}
	}
	return "", nil
}
