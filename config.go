package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	defaultGrimQueueName      = "grim-queue"
	defaultConfigRoot         = "/etc/grim"
	defaultResultRoot         = "/var/log/grim"
	defaultWorkspaceRoot      = "/var/tmp/grim"
	configFileName            = "config.json"
	prepareScriptName         = "prepare.sh"
	buildScriptName           = "build.sh"
	repoBuildScriptName       = "grim_build.sh"
	repoHiddenBuildScriptName = ".grim_build.sh"
)

type config struct {
	GrimQueueName *string
	ResultRoot    *string
	WorkspaceRoot *string
	AWSRegion     *string
	AWSKey        *string
	AWSSecret     *string
	GitHubToken   *string
	PathToCloneIn *string
	HipChatRoom   *string
	HipChatToken  *string
}

type effectiveConfig struct {
	grimQueueName string
	resultRoot    string
	workspaceRoot string
	awsRegion     string
	awsKey        string
	awsSecret     string
	gitHubToken   string
	pathToCloneIn string
	hipChatRoom   string
	hipChatToken  string
}

type repo struct {
	owner, name string
}

func getEffectiveConfigRoot(configRootPtr *string) string {
	if stringPtrNotEmpty(configRootPtr) {
		return *configRootPtr
	}

	return defaultConfigRoot
}

func getAllConfiguredRepos(configRoot string) []repo {
	var repos []repo

	repoPattern := filepath.Join(configRoot, "*/*")
	matches, err := filepath.Glob(repoPattern)
	if err != nil {
		return repos
	}

	for _, match := range matches {
		fi, fiErr := os.Stat(match)
		if fiErr != nil || !fi.Mode().IsDir() {
			continue
		}

		rel, relErr := filepath.Rel(configRoot, match)
		if relErr != nil {
			continue
		}

		owner, name := filepath.Split(rel)
		if owner != "" && name != "" {
			repos = append(repos, repo{filepath.Clean(owner), filepath.Clean(name)})
		}
	}

	return repos
}

func getEffectiveGlobalConfig(configRoot string) (*effectiveConfig, error) {
	var err error
	var global *config

	if global, err = loadGlobalConfig(configRoot); err == nil {
		ec := effectiveConfig{
			grimQueueName: firstNonEmptyStringPtr(global.GrimQueueName, &defaultGrimQueueName),
			resultRoot:    firstNonEmptyStringPtr(global.ResultRoot, &defaultResultRoot),
			workspaceRoot: firstNonEmptyStringPtr(global.WorkspaceRoot, &defaultWorkspaceRoot),
			awsRegion:     firstNonEmptyStringPtr(global.AWSRegion),
			awsKey:        firstNonEmptyStringPtr(global.AWSKey),
			awsSecret:     firstNonEmptyStringPtr(global.AWSSecret),
			gitHubToken:   firstNonEmptyStringPtr(global.GitHubToken),
			hipChatRoom:   firstNonEmptyStringPtr(global.HipChatRoom),
			hipChatToken:  firstNonEmptyStringPtr(global.HipChatToken),
		}

		if err = validateEffectiveConfig(ec); err == nil {
			return &ec, nil
		}
	}

	return nil, err
}

func getEffectiveConfig(configRoot, owner, repo string) (*effectiveConfig, error) {
	var err error
	var global, local *config

	if global, err = loadGlobalConfig(configRoot); err == nil {
		if local, err = loadLocalConfig(configRoot, owner, repo); err == nil {
			ec := effectiveConfig{
				grimQueueName: firstNonEmptyStringPtr(global.GrimQueueName, &defaultGrimQueueName),
				resultRoot:    firstNonEmptyStringPtr(global.ResultRoot, &defaultResultRoot),
				workspaceRoot: firstNonEmptyStringPtr(global.WorkspaceRoot, &defaultWorkspaceRoot),
				awsRegion:     firstNonEmptyStringPtr(global.AWSRegion),
				awsKey:        firstNonEmptyStringPtr(global.AWSKey),
				awsSecret:     firstNonEmptyStringPtr(global.AWSSecret),
				gitHubToken:   firstNonEmptyStringPtr(local.GitHubToken, global.GitHubToken),
				pathToCloneIn: firstNonEmptyStringPtr(local.PathToCloneIn),
				hipChatRoom:   firstNonEmptyStringPtr(local.HipChatRoom, global.HipChatRoom),
				hipChatToken:  firstNonEmptyStringPtr(local.HipChatToken, global.HipChatToken),
			}

			if err = validateEffectiveConfig(ec); err == nil {
				return &ec, nil
			}
		}
	}

	return nil, err
}

func validateEffectiveConfig(ec effectiveConfig) error {
	if ec.awsRegion == "" {
		return fmt.Errorf("AWS region is required")
	} else if ec.awsKey == "" {
		return fmt.Errorf("AWS key is required")
	} else if ec.awsSecret == "" {
		return fmt.Errorf("AWS secret is required")
	}

	return nil
}

func loadGlobalConfig(configRoot string) (*config, error) {
	return loadConfig(filepath.Join(configRoot, configFileName))
}

func loadLocalConfig(configRoot, owner, repo string) (*config, error) {
	return loadConfig(filepath.Join(configRoot, owner, repo, configFileName))
}

func loadConfig(path string) (*config, error) {
	c := new(config)

	if bs, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bs, c); err != nil {
		return nil, err
	}

	return c, nil
}

func firstNonEmptyStringPtr(strPtrs ...*string) string {
	for _, strPtr := range strPtrs {
		if stringPtrNotEmpty(strPtr) {
			return *strPtr
		}
	}

	return ""
}

func stringPtrNotEmpty(strPtr *string) bool {
	return strPtr != nil && *strPtr != ""
}
