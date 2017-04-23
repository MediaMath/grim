package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

type hookEvent struct {
	EventName string
	Action    string
	UserName  string
	Owner     string
	Repo      string
	Target    string
	Ref       string
	StatusRef string
	URL       string
	PrNumber  int64
	Deleted   bool
}

func (hook hookEvent) Describe() string {
	return fmt.Sprintf("hook of %v/%v initiated by a %q to %q by %q", hook.Owner, hook.Repo, hook.EventName, hook.Target, hook.UserName)
}

func (hook hookEvent) env() []string {
	return []string{
		fmt.Sprintf("GH_EVENT_NAME=%v", hook.EventName),
		fmt.Sprintf("GH_ACTION=%v", hook.Action),
		fmt.Sprintf("GH_USER_NAME=%v", hook.UserName),
		fmt.Sprintf("GH_OWNER=%v", hook.Owner),
		fmt.Sprintf("GH_REPO=%v", hook.Repo),
		fmt.Sprintf("GH_TARGET=%v", hook.Target),
		fmt.Sprintf("GH_REF=%v", hook.Ref),
		fmt.Sprintf("GH_STATUS_REF=%v", hook.StatusRef),
		fmt.Sprintf("GH_URL=%v", hook.URL),
		fmt.Sprintf("GH_PR_NUMBER=%v", hook.PrNumber),
	}
}

type pullRequest struct {
	URL            string `json:"html_url"`
	MergeCommitSha string `json:"merge_commit_sha"`
	Head           struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	} `json:"base"`
}

type githubHook struct {
	// Pull Request fields
	Action      string      `json:"action"`
	Number      int64       `json:"number"`
	PullRequest pullRequest `json:"pull_request"`

	// Push fields
	Ref        string `json:"ref"`
	Compare    string `json:"compare"`
	Deleted    bool   `json:"deleted"`
	HeadCommit struct {
		ID string `json:"id"`
	} `json:"head_commit"`

	// Common fields
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
	Repository struct {
		Owner struct {
			Name  string `json:"name"`
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	} `json:"repository"`
}

type hookWrapper struct {
	Message string
}

func extractHookEvent(body string) (*hookEvent, error) {
	wrapper := new(hookWrapper)
	err := json.Unmarshal([]byte(body), wrapper)
	if err != nil {
		return nil, err
	}

	parsed := new(githubHook)
	err = json.Unmarshal([]byte(wrapper.Message), parsed)
	if err != nil {
		return nil, err
	}

	hook := new(hookEvent)

	hook.UserName = parsed.Sender.Login
	hook.Repo = parsed.Repository.Name
	hook.Deleted = parsed.Deleted
	if parsed.Action != "" {
		hook.EventName = "pull_request"
		hook.Action = parsed.Action
		hook.Owner = parsed.Repository.Owner.Login
		hook.Target = parsed.PullRequest.Base.Ref
		hook.StatusRef = parsed.PullRequest.Head.Sha
		hook.URL = parsed.PullRequest.URL
		hook.PrNumber = parsed.Number
	} else {
		hook.EventName = "push"
		hook.Owner = parsed.Repository.Owner.Name
		hook.Target = parsed.Ref
		hook.Ref = parsed.HeadCommit.ID
		hook.StatusRef = parsed.HeadCommit.ID
		hook.URL = parsed.Compare
	}

	hook.Target = strings.TrimPrefix(hook.Target, "refs/heads/")

	return hook, nil
}

func prepareAmazonSNSService(token, owner, repo, snsTopic, awsKey, awsSecret, awsRegion string) error {
	client, err := getClientForToken(token)
	if err != nil {
		return err
	}

	hookID, err := findExistingAmazonSNSHookID(client, owner, repo)
	if hookID == 0 || err != nil {
		err = createAmazonSNSHook(client, owner, repo, snsTopic, awsKey, awsSecret, awsRegion)
	} else {
		err = editAmazonSNSHook(client, owner, repo, snsTopic, awsKey, awsSecret, awsRegion, hookID)
	}

	return err
}

func findExistingAmazonSNSHookID(client *github.Client, owner, repo string) (int, error) {
	listOptions := github.ListOptions{Page: 1, PerPage: 100}

	for {
		hooks, res, err := client.Repositories.ListHooks(context.Background(), owner, repo, &listOptions)
		if err != nil {
			return 0, err
		}
		for _, hook := range hooks {
			if hook.Name != nil && *hook.Name == "amazonsns" && hook.ID != nil {
				return *hook.ID, nil
			}
		}
		if res.NextPage == 0 {
			break
		}
		listOptions.Page = res.NextPage
	}

	return 0, nil
}

func createAmazonSNSHook(client *github.Client, owner, repo, snsTopic, awsKey, awsSecret, awsRegion string) error {
	hook, _, err := client.Repositories.CreateHook(context.Background(), owner, repo, githubAmazonSNSHookStruct(snsTopic, awsKey, awsSecret, awsRegion))

	return detectHookError(hook, err)
}

func editAmazonSNSHook(client *github.Client, owner, repo, snsTopic, awsKey, awsSecret, awsRegion string, hookID int) error {
	hook, _, err := client.Repositories.EditHook(context.Background(), owner, repo, hookID, githubAmazonSNSHookStruct(snsTopic, awsKey, awsSecret, awsRegion))

	return detectHookError(hook, err)
}

func githubAmazonSNSHookStruct(snsTopic, awsKey, awsSecret, awsRegion string) *github.Hook {
	name := "amazonsns"
	active := true
	return &github.Hook{
		Name:   &name,
		Events: []string{"push", "pull_request"},
		Active: &active,
		Config: map[string]interface{}{
			"sns_topic":  snsTopic,
			"aws_key":    awsKey,
			"aws_secret": awsSecret,
			"sns_region": awsRegion,
		},
	}
}

func detectHookError(hook *github.Hook, err error) error {
	if err != nil {
		return err
	}

	if hook == nil {
		return fmt.Errorf("github client returned nil for hook")
	}

	if hook.ID == nil {
		return fmt.Errorf("github client returned nil for hook id")
	}

	return nil
}
