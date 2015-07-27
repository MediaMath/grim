package grim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type acceptsLocalConfig interface {
	configRoot() string
	owner() string
	repo() string
	setLocalConfig(localConfig)
	setError(error)
}

var errNilLocalConfig = fmt.Errorf("local config was nil")

func localConfigReaderProcess(reqs chan acceptsLocalConfig) {
	for req := range reqs {
		localConfigReader(req)
	}
}

func localConfigReader(req acceptsLocalConfig) {
	lc, err := readLocalConfig(req.configRoot(), req.owner(), req.repo())
	if err != nil {
		req.setError(err)
	} else if lc.local == nil {
		req.setError(errNilLocalConfig)
	} else {
		req.setLocalConfig(lc)
	}
}

func readLocalConfig(configRoot, owner, repo string) (lc localConfig, err error) {
	var bs []byte
	global, err := readGlobalConfig(configRoot)
	if err == nil {
		local := make(configMap)

		bs, err = ioutil.ReadFile(filepath.Join(configRoot, owner, repo, configFileName))
		if err == nil {
			err = json.Unmarshal(bs, &local)
			if err == nil {
				lc = localConfig{owner, repo, local, global}
			}
		}
	}

	return
}

type localConfig struct {
	owner, repo string
	local       configMap
	global      globalConfig
}

func (lc localConfig) errors() (errs []error) {
	snsTopicName := lc.snsTopicName()
	if snsTopicName == "" {
		errs = append(errs, fmt.Errorf("Must have a Sns topic name!"))
	} else if strings.Contains(snsTopicName, ".") {
		errs = append(errs, fmt.Errorf("Cannot have . in sns topic name.  Default topic names can be set in the build config file using the SnsTopicName parameter."))
	}
	return
}

func (lc localConfig) warnings() (errs []error) {
	return
}

func (lc localConfig) grimQueueName() string {
	return lc.global.grimQueueName()
}

func (lc localConfig) resultRoot() string {
	return lc.global.resultRoot()
}

func (lc localConfig) workspaceRoot() string {
	return lc.global.workspaceRoot()
}

func (lc localConfig) awsRegion() string {
	return lc.global.awsRegion()
}

func (lc localConfig) awsKey() string {
	return lc.global.awsKey()
}

func (lc localConfig) awsSecret() string {
	return lc.global.awsSecret()
}

func (lc localConfig) gitHubToken() string {
	return readStringWithDefaults(lc.local, "GitHubToken", lc.global.gitHubToken())
}

func (lc localConfig) pathToCloneIn() string {
	return readStringWithDefaults(lc.local, "PathToCloneIn")
}

func (lc localConfig) snsTopicName() string {
	return readStringWithDefaults(lc.local, "SNSTopicName", *defaultTopicName(lc.owner, lc.repo))
}

func (lc localConfig) hipChatRoom() string {
	return readStringWithDefaults(lc.local, "HipChatRoom", lc.global.hipChatRoom())
}

func (lc localConfig) hipChatToken() string {
	return readStringWithDefaults(lc.local, "HipChatToken", lc.global.hipChatToken())
}

func (lc localConfig) grimServerID() string {
	return lc.global.grimServerID()
}

func (lc localConfig) pendingTemplate() string {
	return readStringWithDefaults(lc.local, "PendingTemplate", lc.global.pendingTemplate())
}

func (lc localConfig) errorTemplate() string {
	return readStringWithDefaults(lc.local, "ErrorTemplate", lc.global.errorTemplate())
}

func (lc localConfig) successTemplate() string {
	return readStringWithDefaults(lc.local, "SuccessTemplate", lc.global.successTemplate())
}

func (lc localConfig) failureTemplate() string {
	return readStringWithDefaults(lc.local, "FailureTemplate", lc.global.failureTemplate())
}

func (lc localConfig) timeout() (to time.Duration) {
	val := readIntWithDefaults(lc.local, "Timeout")

	if val > 0 {
		to = time.Duration(val) * time.Second
	} else {
		to = lc.global.timeout()
	}

	return
}

func (lc localConfig) usernameWhitelist() []string {
	val, _ := lc.local["UsernameWhitelist"]
	iSlice, _ := val.([]interface{})
	var wl []string
	for _, entry := range iSlice {
		entryStr, _ := entry.(string)
		wl = append(wl, entryStr)
	}
	return wl
}

func defaultTopicName(owner, repo string) *string {
	snsTopicName := fmt.Sprintf("grim-%v-%v-repo-topic", owner, repo)
	return &snsTopicName
}
