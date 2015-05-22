package grim

import "testing"

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func TestGlobalEffectiveFailureTemplate(t *testing.T) {
	global := &config{FailureTemplate: s("template")}

	if ec := buildGlobalEffectiveConfig(global); ec.failureTemplate != "template" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.failureTemplate != ps(templateFor("Failure during")) {
		t.Errorf("No defaulting %v", ec)
	}
}

func TestGlobalEffectiveSuccessTemplate(t *testing.T) {
	global := &config{SuccessTemplate: s("template")}

	if ec := buildGlobalEffectiveConfig(global); ec.successTemplate != "template" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.successTemplate != ps(templateFor("Success after")) {
		t.Errorf("No defaulting %v", ec)
	}
}

func TestGlobalEffectiveErrorTemplate(t *testing.T) {
	global := &config{ErrorTemplate: s("template")}

	if ec := buildGlobalEffectiveConfig(global); ec.errorTemplate != "template" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.errorTemplate != ps(templateFor("Error during")) {
		t.Errorf("No defaulting %v", ec)
	}
}

func TestGlobalEffectivePendingTemplate(t *testing.T) {
	global := &config{PendingTemplate: s("template")}

	if ec := buildGlobalEffectiveConfig(global); ec.pendingTemplate != "template" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.pendingTemplate != ps(templateFor("Starting")) {
		t.Errorf("No defaulting %v", ec)
	}
}

func TestGlobalEffectiveGrimServerId(t *testing.T) {
	global := &config{GrimServerID: s("id")}

	if ec := buildGlobalEffectiveConfig(global); ec.grimServerID != "id" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	noidButQueue := &config{GrimQueueName: s("q")}

	if ec := buildGlobalEffectiveConfig(noidButQueue); ec.grimServerID != "q" {
		t.Errorf("No defaulting to q name %v", ec.workspaceRoot)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.grimServerID != "grim-queue" {
		t.Errorf("No defaulting to queue name default %v", ec)
	}
}

func TestGlobalEffectiveWorkSpaceRoot(t *testing.T) {
	global := &config{WorkspaceRoot: s("ws")}

	if ec := buildGlobalEffectiveConfig(global); ec.workspaceRoot != "ws" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.workspaceRoot != "/var/tmp/grim" {
		t.Errorf("No defaulting %v", ec)
	}
}

func TestGlobalEffectiveResultRoot(t *testing.T) {
	global := &config{ResultRoot: s("result")}

	if ec := buildGlobalEffectiveConfig(global); ec.resultRoot != "result" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.resultRoot != "/var/log/grim" {
		t.Errorf("No defaulting %v", ec)
	}
}

func TstGlobalEffectiveGrimQueueName(t *testing.T) {
	global := &config{GrimQueueName: s("queue")}

	if ec := buildGlobalEffectiveConfig(global); ec.grimQueueName != "queue" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.grimQueueName != "grim-queue" {
		t.Errorf("No defaulting %v", ec)
	}
}

func TestGlobalEffectiveConfigNoDefaults(t *testing.T) {
	global := &config{
		AWSRegion:    s("region"),
		AWSKey:       s("key"),
		AWSSecret:    s("secret"),
		GitHubToken:  s("ghtoken"),
		HipChatRoom:  s("hcRoom"),
		HipChatToken: s("hcToken")}

	if ec := buildGlobalEffectiveConfig(global); ec.awsRegion != "region" ||
		ec.awsKey != "key" ||
		ec.awsSecret != "secret" ||
		ec.gitHubToken != "ghtoken" ||
		ec.hipChatRoom != "hcRoom" ||
		ec.hipChatToken != "hcToken" {
		t.Errorf("Did not set effective correctly %v", ec)
	}

	none := &config{}

	if ec := buildGlobalEffectiveConfig(none); ec.awsRegion != "" ||
		ec.awsKey != "" ||
		ec.awsSecret != "" ||
		ec.gitHubToken != "" ||
		ec.hipChatRoom != "" ||
		ec.hipChatToken != "" {
		t.Errorf("Defaulted undefaultable vals %v", ec)
	}
}

func TestLocalEffectiveConfigPathIsNotGlobal(t *testing.T) {
	global := effectiveConfig{pathToCloneIn: "global"}
	has := config{PathToCloneIn: s("local")}
	none := config{}

	if ec := buildLocalEffectiveConfig(global, &has); ec.pathToCloneIn != "local" {
		t.Errorf("local didnt exists %v", ec)
	}

	if ec := buildLocalEffectiveConfig(global, &none); ec.pathToCloneIn != "" {
		t.Errorf("had global path %v", ec)
	}
}

func TestLocalEffectiveConfigDoesOverwriteGlobals(t *testing.T) {
	global := effectiveConfig{
		pendingTemplate: "global",
		errorTemplate:   "global",
		successTemplate: "global",
		failureTemplate: "global",
		gitHubToken:     "global",
		pathToCloneIn:   "global",
		hipChatRoom:     "global",
		hipChatToken:    "global"}

	has := config{
		PendingTemplate: s("local"),
		ErrorTemplate:   s("local"),
		SuccessTemplate: s("local"),
		FailureTemplate: s("local"),
		GitHubToken:     s("local"),
		PathToCloneIn:   s("local"),
		HipChatRoom:     s("local"),
		HipChatToken:    s("local")}

	none := config{}

	if ec := buildLocalEffectiveConfig(global, &has); ec.gitHubToken != "local" ||
		ec.pendingTemplate != "local" ||
		ec.errorTemplate != "local" ||
		ec.successTemplate != "local" ||
		ec.failureTemplate != "local" ||
		ec.hipChatRoom != "local" ||
		ec.hipChatToken != "local" {
		t.Errorf("local did not overwrite global %v", ec)
	}

	if ec := buildLocalEffectiveConfig(global, &none); ec.gitHubToken != "global" ||
		ec.pendingTemplate != "global" ||
		ec.errorTemplate != "global" ||
		ec.successTemplate != "global" ||
		ec.failureTemplate != "global" ||
		ec.hipChatRoom != "global" ||
		ec.hipChatToken != "global" {
		t.Errorf("global does not back stop local %v", ec)
	}

}

func TestLocalEffectiveConfigDoesntOverwriteGlobals(t *testing.T) {
	global := effectiveConfig{
		grimQueueName: "global.grimQueueName",
		resultRoot:    "global.resultRoot",
		workspaceRoot: "global.workspaceRoot",
		awsRegion:     "global.awsRegion",
		awsKey:        "global.awsKey",
		awsSecret:     "global.awsSecret",
		grimServerID:  "global.grimServerID"}

	local := config{
		GrimQueueName: s("local.grimQueueName"),
		ResultRoot:    s("local.resultRoot"),
		WorkspaceRoot: s("local.workspaceRoot"),
		AWSRegion:     s("local.awsRegion"),
		AWSKey:        s("local.awsKey"),
		AWSSecret:     s("local.awsSecret"),
		GrimServerID:  s("local.grimServerID")}

	ec := buildLocalEffectiveConfig(global, &local)

	if ec.grimQueueName != "global.grimQueueName" ||
		ec.resultRoot != "global.resultRoot" ||
		ec.workspaceRoot != "global.workspaceRoot" ||
		ec.awsRegion != "global.awsRegion" ||
		ec.awsKey != "global.awsKey" ||
		ec.awsSecret != "global.awsSecret" ||
		ec.grimServerID != "global.grimServerID" {
		t.Errorf("local overwrote global. %v", ec)
	}
}

func TestValidateEffectiveConfig(t *testing.T) {
	if err := validateEffectiveConfig(effectiveConfig{}); err == nil {
		t.Errorf("validated with no credentials")
	}

	if err := validateEffectiveConfig(effectiveConfig{awsRegion: "reg", awsKey: "key"}); err == nil {
		t.Errorf("validated with no secret")
	}

	if err := validateEffectiveConfig(effectiveConfig{awsSecret: "secret", awsKey: "key"}); err == nil {
		t.Errorf("validated with no region")
	}

	if err := validateEffectiveConfig(effectiveConfig{awsSecret: "secret", awsRegion: "region"}); err == nil {
		t.Errorf("validated with no key")
	}

	if err := validateEffectiveConfig(effectiveConfig{awsSecret: "secret", awsRegion: "region", awsKey: "key"}); err != nil {
		t.Errorf("didnt validate with all credentials")
	}
}

func TestLoadGlobalConfig(t *testing.T) {
	config, err := loadGlobalConfig("./test_config")
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if ps(config.GrimServerID) != "def-serverid" {
		t.Errorf("Didn't match:\n%v", config)
	}
}

func TestLoadRepoConfig(t *testing.T) {
	config, err := loadLocalConfig("./test_config", "MediaMath", "foo")
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if ps(config.PathToCloneIn) != "go/src/github.com/MediaMath/foo" {
		t.Errorf("Didn't match:\n%v", config)
	}
}

func TestLoadConfig(t *testing.T) {
	config, err := loadConfig("./test_config/config.json")
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if ps(config.HipChatRoom) != "def-hcroom" {
		t.Errorf("Didn't load correctly")
	}
}

func ps(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func s(s string) *string {
	return &s
}

var testEmpty = ""
var testNotEmpty = "foo"

func TestFirstNonEmptyPtrIsEmptyStringInEmptyCase(t *testing.T) {
	if firstNonEmptyStringPtr() != "" {
		t.Errorf("Didnt default")
	}
}

func TestFirstNonEmptyPtr(t *testing.T) {

	if firstNonEmptyStringPtr(&testNotEmpty, nil, &testEmpty) != testNotEmpty {
		t.Errorf("didn't find it first")
	}

	if firstNonEmptyStringPtr(nil, &testNotEmpty, nil, &testEmpty) != testNotEmpty {
		t.Errorf("didn't find it with nil in front ")
	}

	if firstNonEmptyStringPtr(&testEmpty, &testNotEmpty, nil, &testEmpty) != testNotEmpty {
		t.Errorf("didn't find it with empty in front")
	}

	var secondEmpty = "second"
	if firstNonEmptyStringPtr(&testEmpty, &testNotEmpty, &secondEmpty) != testNotEmpty {
		t.Errorf("didn't find it with other string")
	}

}

func TestStringPtrNotEmpty(t *testing.T) {
	if stringPtrNotEmpty(nil) {
		t.Errorf("Thinks nil is not empty")
	}

	var testEmpty = ""
	if stringPtrNotEmpty(&testEmpty) {
		t.Errorf("Thinks empty is not empty")
	}

	if !stringPtrNotEmpty(&testNotEmpty) {
		t.Errorf("Thinks not empty is")
	}
}
