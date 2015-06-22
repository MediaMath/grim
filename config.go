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
	"strings"
)

var (
	defaultGrimQueueName      = "grim-queue"
	defaultConfigRoot         = "/etc/grim"
	defaultResultRoot         = "/var/log/grim"
	defaultWorkspaceRoot      = "/var/tmp/grim"
	defaultTimeOutSeconds	  = 30
	configFileName            = "config.json"
	buildScriptName           = "build.sh"
	repoBuildScriptName       = "grim_build.sh"
	repoHiddenBuildScriptName = ".grim_build.sh"
)

type config struct {
	GrimQueueName     *string
	ResultRoot        *string
	WorkspaceRoot     *string
	AWSRegion         *string
	AWSKey            *string
	AWSSecret         *string
	GitHubToken       *string
	PathToCloneIn     *string
	SnsTopicName      *string
	HipChatRoom       *string
	HipChatToken      *string
	GrimServerID      *string
	PendingTemplate   *string
	ErrorTemplate     *string
	SuccessTemplate   *string
	FailureTemplate   *string
	Timeout           int
	UsernameWhitelist []string
}

type effectiveConfig struct {
	grimQueueName     string
	resultRoot        string
	workspaceRoot     string
	awsRegion         string
	awsKey            string
	awsSecret         string
	gitHubToken       string
	pathToCloneIn     string
	snsTopicName      string
	hipChatRoom       string
	hipChatToken      string
	grimServerID      string
	origServerID      string
	truncateID        string
	pendingTemplate   string
	errorTemplate     string
	successTemplate   string
	failureTemplate   string
	timeout           int
	usernameWhitelist []string
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
		ec := buildGlobalEffectiveConfig(global)

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
			ec := buildLocalEffectiveConfig(buildGlobalEffectiveConfig(global), local, owner, repo)

			if err = validateLocalEffectiveConfig(ec); err != nil {
				return nil, err
			}

			if err = validateEffectiveConfig(ec); err == nil {
				return &ec, nil
			}

		}
	}

	return nil, err
}

func buildGlobalEffectiveConfig(global *config) effectiveConfig {
	origServerID := firstNonEmptyStringPtr(global.GrimServerID, global.GrimQueueName, &defaultGrimQueueName)
	truncatedServerID := origServerID
	truncateID := ""

	if len(truncatedServerID) > 15 {
		truncatedServerID = fmt.Sprintf("%s", origServerID[:15])

		truncateID = "GrimServerID"
		if !stringPtrNotEmpty(global.GrimServerID) {
			truncateID = "GrimQueueName"
		}
	}

	return effectiveConfig{
		grimQueueName:   firstNonEmptyStringPtr(global.GrimQueueName, &defaultGrimQueueName),
		resultRoot:      firstNonEmptyStringPtr(global.ResultRoot, &defaultResultRoot),
		workspaceRoot:   firstNonEmptyStringPtr(global.WorkspaceRoot, &defaultWorkspaceRoot),
		grimServerID:    truncatedServerID,
		origServerID:    origServerID,
		truncateID:      truncateID,
		awsRegion:       firstNonEmptyStringPtr(global.AWSRegion),
		awsKey:          firstNonEmptyStringPtr(global.AWSKey),
		awsSecret:       firstNonEmptyStringPtr(global.AWSSecret),
		gitHubToken:     firstNonEmptyStringPtr(global.GitHubToken),
		hipChatRoom:     firstNonEmptyStringPtr(global.HipChatRoom),
		hipChatToken:    firstNonEmptyStringPtr(global.HipChatToken),
		pendingTemplate: firstNonEmptyStringPtr(global.PendingTemplate, templateForStart()),
		errorTemplate:   firstNonEmptyStringPtr(global.ErrorTemplate, templateForFailureandError("Error during")),
		failureTemplate: firstNonEmptyStringPtr(global.FailureTemplate, templateForFailureandError("Failure during")),
		successTemplate: firstNonEmptyStringPtr(global.SuccessTemplate, templateForSuccess()),
		timeout:         firstNonZeroInt(global.Timeout, defaultTimeOutSeconds),
	}
}

func buildLocalEffectiveConfig(global effectiveConfig, local *config, owner, repo string) effectiveConfig {
	return effectiveConfig{
		grimQueueName:     global.grimQueueName,
		resultRoot:        global.resultRoot,
		workspaceRoot:     global.workspaceRoot,
		awsRegion:         global.awsRegion,
		awsKey:            global.awsKey,
		awsSecret:         global.awsSecret,
		grimServerID:      global.grimServerID,
		gitHubToken:       firstNonEmptyStringPtr(local.GitHubToken, &global.gitHubToken),
		hipChatRoom:       firstNonEmptyStringPtr(local.HipChatRoom, &global.hipChatRoom),
		hipChatToken:      firstNonEmptyStringPtr(local.HipChatToken, &global.hipChatToken),
		pendingTemplate:   firstNonEmptyStringPtr(local.PendingTemplate, &global.pendingTemplate),
		successTemplate:   firstNonEmptyStringPtr(local.SuccessTemplate, &global.successTemplate),
		errorTemplate:     firstNonEmptyStringPtr(local.ErrorTemplate, &global.errorTemplate),
		failureTemplate:   firstNonEmptyStringPtr(local.FailureTemplate, &global.failureTemplate),
		pathToCloneIn:     firstNonEmptyStringPtr(local.PathToCloneIn),
		snsTopicName:      firstNonEmptyStringPtr(local.SnsTopicName, defaultTopicName(owner, repo)),
		timeout:           firstNonZeroInt(local.Timeout, global.timeout),
		usernameWhitelist: local.UsernameWhitelist,
	}
}

func defaultTopicName(owner, repo string) *string {
	snsTopicName := fmt.Sprintf("grim-%v-%v-repo-topic", owner, repo)
	return &snsTopicName
}

func templateForStart() *string {
	s := fmt.Sprintf("Starting build of {{.Owner}}/{{.Repo}} initiated by a {{.EventName}} to {{.Target}} by {{.UserName}}")
	return &s
}

func templateForSuccess() *string {
	s := fmt.Sprintf("Success after build of {{.Owner}}/{{.Repo}} initiated by a {{.EventName}} to {{.Target}} by {{.UserName}} ({{.Workspace}})")
	return &s
}

func templateForFailureandError(preamble string) *string {
	s := fmt.Sprintf("%s build of {{.Owner}}/{{.Repo}} initiated by a {{.EventName}} to {{.Target}} by {{.UserName}} ({{.LogDir}})", preamble)
	return &s
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

func validateLocalEffectiveConfig(ec effectiveConfig) error {
	if ec.snsTopicName == "" {
		return fmt.Errorf("Must have a Sns topic name!")
	}

	if strings.Contains(ec.snsTopicName, ".") {
		return fmt.Errorf("Cannot have . in sns topic name.  Default topic names can be set in the build config file using the SnsTopicName parameter.")
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

func firstNonZeroInt(ints ...int) int {
	for _, value := range ints {
		if value > 0 {
			return value
		}
	}
	return 0
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

func (ec *effectiveConfig) usernameCanBuild(username string) (allowed bool) {
	wlLen := len(ec.usernameWhitelist)

	if ec.usernameWhitelist == nil || wlLen == 0 {
		allowed = true
	} else {
		for i := 0; i < wlLen; i++ {
			if ec.usernameWhitelist[i] == username {
				allowed = true
				break
			}
		}
	}

	return
}
