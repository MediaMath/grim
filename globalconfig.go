package grim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

var errNilGlobalConfig = fmt.Errorf("global config was nil")

func readGlobalConfig(configRoot string) (gc globalConfig, err error) {
	gc = make(globalConfig)

	bs, err := ioutil.ReadFile(filepath.Join(configRoot, configFileName))
	if err == nil {
		err = json.Unmarshal(bs, &gc)
	}

	return
}

type globalConfig configMap

func (gc globalConfig) errors() (errs []error) {
	if gc.awsRegion() == "" {
		errs = append(errs, fmt.Errorf("AWS region is required"))
	}

	if gc.awsKey() == "" {
		errs = append(errs, fmt.Errorf("AWS key is required"))
	}

	if gc.awsSecret() == "" {
		errs = append(errs, fmt.Errorf("AWS secret is required"))
	}

	return
}

func (gc globalConfig) warnings() (errs []error) {
	rawSID := gc.rawGrimServerID()
	actualSID := gc.grimServerID()
	queueName := gc.grimQueueName()
	instructions := `you can override this by setting "GrimServerID" in your config`

	if rawSID == "" {
		if queueName == "" {
			errs = append(errs, fmt.Errorf(`using default %q as server id; %v`, defaultGrimQueueName, instructions))
		} else {
			errs = append(errs, fmt.Errorf("using queue name %q as server id; %v", queueName, instructions))
		}
	}

	if rawSID != actualSID {
		errs = append(errs, fmt.Errorf("truncating server id from %q to %q", rawSID, actualSID))
	}

	return
}

func (gc globalConfig) grimQueueName() string {
	return readStringWithDefaults(gc, "GrimQueueName", defaultGrimQueueName)
}

func (gc globalConfig) resultRoot() string {
	return readStringWithDefaults(gc, "ResultRoot", defaultResultRoot)
}

func (gc globalConfig) workspaceRoot() string {
	return readStringWithDefaults(gc, "WorkspaceRoot", defaultWorkspaceRoot)
}

func (gc globalConfig) awsRegion() string {
	return readStringWithDefaults(gc, "AWSRegion")
}

func (gc globalConfig) awsKey() string {
	return readStringWithDefaults(gc, "AWSKey")
}

func (gc globalConfig) awsSecret() string {
	return readStringWithDefaults(gc, "AWSSecret")
}

func (gc globalConfig) gitHubToken() string {
	return readStringWithDefaults(gc, "GitHubToken")
}

func (gc globalConfig) snsTopicName() string {
	return readStringWithDefaults(gc, "SNSTopicName")
}

func (gc globalConfig) hipChatRoom() string {
	return readStringWithDefaults(gc, "HipChatRoom")
}

func (gc globalConfig) hipChatToken() string {
	return readStringWithDefaults(gc, "HipChatToken")
}

func (gc globalConfig) hipChatVersion() int {
	return readIntWithDefaults(gc, "HipChatVersion", defaultHipChatVersion)
}

func (gc globalConfig) grimServerID() string {
	sid := gc.rawGrimServerID()
	if len(sid) > 15 {
		sid = sid[:15]
	}
	return sid
}

func (gc globalConfig) rawGrimServerID() string {
	return readStringWithDefaults(gc, "GrimServerID", gc.grimQueueName(), defaultGrimQueueName)
}

func (gc globalConfig) grimServerIDSource() string {
	if _, ok := gc["GrimServerID"]; ok {
		return "GrimServerID"
	}

	if _, ok := gc["GrimQueueName"]; ok {
		return "GrimQueueName"
	}

	return ""
}

func (gc globalConfig) grimServerIDWasTruncated() bool {
	return gc.grimServerID() != gc.rawGrimServerID()
}

func (gc globalConfig) pendingTemplate() string {
	return readStringWithDefaults(gc, "PendingTemplate", *defaultTemplateForStart)
}

func (gc globalConfig) errorTemplate() string {
	return readStringWithDefaults(gc, "ErrorTemplate", *defaultTemplateForError)
}

func (gc globalConfig) successTemplate() string {
	return readStringWithDefaults(gc, "SuccessTemplate", *defaultTemplateForSuccess)
}

func (gc globalConfig) successColor() string {
	return readStringWithDefaults(gc, "SuccessColor", *defaultColorForSuccess)
}

func (gc globalConfig) errorColor() string {
	return readStringWithDefaults(gc, "ErrorColor", *defaultColorForError)
}

func (gc globalConfig) failureColor() string {
	return readStringWithDefaults(gc, "FailureColor", *defaultColorForFailure)
}

func (gc globalConfig) pendingColor() string {
	return readStringWithDefaults(gc, "PendingColor", *defaultColorForPending)
}

func (gc globalConfig) failureTemplate() string {
	return readStringWithDefaults(gc, "FailureTemplate", *defaultTemplateForFailure)
}

func (gc globalConfig) timeout() (to time.Duration) {
	val := readIntWithDefaults(gc, "Timeout")

	if val > 0 {
		to = time.Duration(val) * time.Second
	} else {
		to = defaultTimeout
	}

	return
}
