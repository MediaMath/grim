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

var errNilLocalConfig = fmt.Errorf("local config was nil")

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
		errs = append(errs, fmt.Errorf("Cannot have . in sns topic name [ %s ].  Default topic names can be set in the build config file using the SnsTopicName parameter.", snsTopicName))
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

func (lc localConfig) hipChatVersion() int {
	return readIntWithDefaults(lc.local, "HipChatVersion", lc.global.hipChatVersion())
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

func (lc localConfig) successColor() string {
	return readStringWithDefaults(lc.local, "SuccessColor", lc.global.successColor())
}

func (lc localConfig) errorColor() string {
	return readStringWithDefaults(lc.local, "ErrorColor", lc.global.errorColor())
}

func (lc localConfig) failureColor() string {
	return readStringWithDefaults(lc.local, "FailureColor", lc.global.failureColor())
}

func (lc localConfig) pendingColor() string {
	return readStringWithDefaults(lc.local, "PendingColor", lc.global.pendingColor())
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

func (lc localConfig) usernameCanBuild(username string) (allowed bool) {
	whitelist := lc.usernameWhitelist()

	wlLen := len(whitelist)

	if whitelist == nil || wlLen == 0 {
		allowed = true
	} else {
		for i := 0; i < wlLen; i++ {
			if whitelist[i] == username {
				allowed = true
				break
			}
		}
	}

	return
}

func defaultTopicName(owner, repo string) *string {
	snsTopicName := fmt.Sprintf("grim-%v-%v-repo-topic", owner, repo)
	return &snsTopicName
}
