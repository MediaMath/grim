package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	defaultGrimQueueName      = "grim-queue"
	defaultConfigRoot         = "/etc/grim"
	defaultResultRoot         = "/var/log/grim"
	defaultWorkspaceRoot      = "/var/tmp/grim"
	defaultTimeout            = 5 * time.Minute
	configFileName            = "config.json"
	buildScriptName           = "build.sh"
	repoBuildScriptName       = "grim_build.sh"
	repoHiddenBuildScriptName = ".grim_build.sh"
	defaultTemplateForStart   = templateForStart()
	defaultTemplateForError   = templateForFailureandError("Error during")
	defaultTemplateForSuccess = templateForSuccess()
	defaultTemplateForFailure = templateForFailureandError("Failure during")
)

type configMap map[string]interface{}

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

//BuildTimeout is a "safe" way to get the timeout configured for builds.  It will ensure non-zero build timeouts.
type BuildTimeout interface {
	BuildTimeout() time.Duration
}

func (ec *effectiveConfig) BuildTimeout() time.Duration {
	seconds := ec.timeout

	if seconds > 0 {
		return time.Second * time.Duration(seconds)
	}

	return defaultTimeout
}

func getEffectiveConfigRoot(configRootPtr *string) string {
	if configRootPtr == nil || *configRootPtr == "" {
		return defaultConfigRoot
	}

	return *configRootPtr
}

type repo struct {
	owner, name string
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

func getEffectiveGlobalConfig(configRoot string) (ec *effectiveConfig, err error) {
	gc, err := readGlobalConfig(configRoot)
	if err == nil {
		errs := gc.errors()

		if len(errs) == 0 {
			truncateID := ""

			if gc.grimServerID() != gc.rawGrimServerID() {
				truncateID = "GrimServerID"
				if _, ok := gc["GrimServerID"]; !ok {
					truncateID = "GrimQueueName"
				}
			}

			ec = &effectiveConfig{
				grimQueueName:     gc.grimQueueName(),
				resultRoot:        gc.resultRoot(),
				workspaceRoot:     gc.workspaceRoot(),
				awsRegion:         gc.awsRegion(),
				awsKey:            gc.awsKey(),
				awsSecret:         gc.awsSecret(),
				gitHubToken:       gc.gitHubToken(),
				snsTopicName:      gc.snsTopicName(),
				hipChatRoom:       gc.hipChatRoom(),
				hipChatToken:      gc.hipChatToken(),
				grimServerID:      gc.grimServerID(),
				origServerID:      gc.rawGrimServerID(),
				truncateID:        truncateID,
				pendingTemplate:   gc.pendingTemplate(),
				errorTemplate:     gc.errorTemplate(),
				successTemplate:   gc.successTemplate(),
				failureTemplate:   gc.failureTemplate(),
				timeout:           10,
				usernameWhitelist: []string{},
			}
		} else {
			err = errs[0]
		}
	}

	return ec, err
}

func getEffectiveConfig(configRoot, owner, repo string) (ec *effectiveConfig, err error) {
	lc, err := readLocalConfig(configRoot, owner, repo)
	if err == nil {
		errs := lc.errors()

		if len(errs) == 0 {
			ec = &effectiveConfig{
				grimQueueName:     lc.grimQueueName(),
				resultRoot:        lc.resultRoot(),
				workspaceRoot:     lc.workspaceRoot(),
				awsRegion:         lc.awsRegion(),
				awsKey:            lc.awsKey(),
				awsSecret:         lc.awsSecret(),
				gitHubToken:       lc.gitHubToken(),
				pathToCloneIn:     lc.pathToCloneIn(),
				snsTopicName:      lc.snsTopicName(),
				hipChatRoom:       lc.hipChatRoom(),
				hipChatToken:      lc.hipChatToken(),
				grimServerID:      lc.grimServerID(),
				pendingTemplate:   lc.pendingTemplate(),
				errorTemplate:     lc.errorTemplate(),
				successTemplate:   lc.successTemplate(),
				failureTemplate:   lc.failureTemplate(),
				timeout:           int(lc.timeout().Seconds()),
				usernameWhitelist: lc.usernameWhitelist(),
			}
		} else {
			err = errs[0]
		}
	}

	return
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

func readStringWithDefaults(m map[string]interface{}, key string, strs ...string) string {
	val, _ := m[key]
	str, _ := val.(string)

	if str == "" {
		for _, str = range strs {
			if str != "" {
				break
			}
		}
	}

	return str
}

func readIntWithDefaults(m map[string]interface{}, key string, ints ...int) int {
	val, _ := m[key]
	i, _ := val.(int)

	if i == 0 {
		for _, i = range ints {
			if i != 0 {
				break
			}
		}
	}

	return i
}
