package grim

import (
	"strings"
	"testing"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func TestGlobalEffectiveFailureTemplate(t *testing.T) {
	gc := globalConfig{"FailureTemplate": "template"}

	if gc.failureTemplate() != "template" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.failureTemplate() != *templateForFailureandError("Failure during") {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveFailureColor(t *testing.T) {
	gc := globalConfig{"FailureColor": "purple"}

	if gc.failureColor() != "purple" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.failureColor() != *colorForFailure() {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveSuccessTemplate(t *testing.T) {
	gc := globalConfig{"SuccessTemplate": "template"}

	if gc.successTemplate() != "template" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.successTemplate() != *templateForSuccess() {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveSuccessColor(t *testing.T) {
	gc := globalConfig{"SuccessColor": "purple"}

	if gc.successColor() != "purple" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.successColor() != *colorForSuccess() {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveErrorTemplate(t *testing.T) {
	gc := globalConfig{"ErrorTemplate": "template"}

	if gc.errorTemplate() != "template" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.errorTemplate() != *templateForFailureandError("Error during") {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveErrorColor(t *testing.T) {
	gc := globalConfig{"ErrorColor": "purple"}

	if gc.errorColor() != "purple" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.errorColor() != *colorForError() {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectivePendingColor(t *testing.T) {
	gc := globalConfig{"PendingColor": "purple"}

	if gc.pendingColor() != "purple" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.pendingColor() != *colorForPending() {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectivePendingTemplate(t *testing.T) {
	gc := globalConfig{"PendingTemplate": "template"}

	if gc.pendingTemplate() != "template" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.pendingTemplate() != *templateForStart() {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveGrimServerId(t *testing.T) {
	gc := globalConfig{"GrimServerID": "id"}

	if gc.grimServerID() != "id" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	noidButQueue := globalConfig{"GrimQueueName": "q"}

	if noidButQueue.grimServerID() != "q" {
		t.Errorf("No defaulting to q name %v", noidButQueue.grimServerID())
	}

	none := globalConfig{}

	if none.grimServerID() != "grim-queue" {
		t.Errorf("No defaulting to queue name default %v", none.grimServerID())
	}
}

func TestGlobalEffectiveWorkSpaceRoot(t *testing.T) {
	gc := globalConfig{"WorkspaceRoot": "ws"}

	if gc.workspaceRoot() != "ws" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.workspaceRoot() != "/var/tmp/grim" {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveResultRoot(t *testing.T) {
	gc := globalConfig{"ResultRoot": "result"}

	if gc.resultRoot() != "result" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.resultRoot() != "/var/log/grim" {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveGrimQueueName(t *testing.T) {
	gc := globalConfig{"GrimQueueName": "queue"}

	if gc.grimQueueName() != "queue" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.grimQueueName() != "grim-queue" {
		t.Errorf("No defaulting %v", none)
	}
}

func TestGlobalEffectiveConfigNoDefaults(t *testing.T) {
	gc := globalConfig{
		"AWSRegion":    "region",
		"AWSKey":       "key",
		"AWSSecret":    "secret",
		"GitHubToken":  "ghtoken",
		"HipChatRoom":  "hcRoom",
		"HipChatToken": "hcToken",
	}

	if gc.awsRegion() != "region" ||
		gc.awsKey() != "key" ||
		gc.awsSecret() != "secret" ||
		gc.gitHubToken() != "ghtoken" ||
		gc.hipChatRoom() != "hcRoom" ||
		gc.hipChatToken() != "hcToken" {
		t.Errorf("Did not set effective correctly %v", gc)
	}

	none := globalConfig{}

	if none.awsRegion() != "" ||
		none.awsKey() != "" ||
		none.awsSecret() != "" ||
		none.gitHubToken() != "" ||
		none.hipChatRoom() != "" ||
		none.hipChatToken() != "" {
		t.Errorf("Defaulted undefaultable vals %v", none)
	}
}

func TestLocalEffectiveConfigSnsTopic(t *testing.T) {
	gc := globalConfig{"SNSTopicName": "global"}
	has := localConfig{"foo", "bar", configMap{"SNSTopicName": "local"}, gc}
	none := localConfig{"foo", "bar", configMap{}, gc}

	if has.snsTopicName() != "local" {
		t.Errorf("local didnt exists %v", has)
	}

	if none.snsTopicName() != "grim-foo-bar-repo-topic" {
		t.Errorf("didnt default %v", none)
	}
}

func TestLocalEffectiveConfigDoesOverwriteGlobals(t *testing.T) {
	gc := globalConfig{
		"PendingTemplate": "global",
		"ErrorTemplate":   "global",
		"SuccessTemplate": "global",
		"FailureTemplate": "global",
		"GitHubToken":     "global",
		"PathToCloneIn":   "global",
		"HipChatRoom":     "global",
		"HipChatToken":    "global",
	}

	has := localConfig{"foo", "bar", configMap{
		"PendingTemplate": "local",
		"ErrorTemplate":   "local",
		"SuccessTemplate": "local",
		"FailureTemplate": "local",
		"GitHubToken":     "local",
		"PathToCloneIn":   "local",
		"HipChatRoom":     "local",
		"HipChatToken":    "local",
	}, gc}

	none := localConfig{"foo", "bar", configMap{}, gc}

	if has.gitHubToken() != "local" ||
		has.pendingTemplate() != "local" ||
		has.errorTemplate() != "local" ||
		has.successTemplate() != "local" ||
		has.failureTemplate() != "local" ||
		has.hipChatRoom() != "local" ||
		has.hipChatToken() != "local" {
		t.Errorf("local did not overwrite global %v", has)
	}

	if none.gitHubToken() != "global" ||
		none.pendingTemplate() != "global" ||
		none.errorTemplate() != "global" ||
		none.successTemplate() != "global" ||
		none.failureTemplate() != "global" ||
		none.hipChatRoom() != "global" ||
		none.hipChatToken() != "global" {
		t.Errorf("global does not back stop local %v", none)
	}

}

func TestLocalEffectiveConfigDoesntOverwriteGlobals(t *testing.T) {
	gc := globalConfig{
		"GrimQueueName": "global.grimQueueName",
		"ResultRoot":    "global.resultRoot",
		"WorkspaceRoot": "global.workspaceRoot",
		"AWSRegion":     "global.awsRegion",
		"AWSKey":        "global.awsKey",
		"AWSSecret":     "global.awsSecret",
		"GrimServerID":  "grimServerID",
	}

	local := localConfig{"foo", "bar", configMap{
		"GrimQueueName": "local.grimQueueName",
		"ResultRoot":    "local.resultRoot",
		"WorkspaceRoot": "local.workspaceRoot",
		"AWSRegion":     "local.awsRegion",
		"AWSKey":        "local.awsKey",
		"AWSSecret":     "local.awsSecret",
		"GrimServerID":  "local.grimServerID",
	}, gc}

	if local.grimQueueName() != "global.grimQueueName" ||
		local.resultRoot() != "global.resultRoot" ||
		local.workspaceRoot() != "global.workspaceRoot" ||
		local.awsRegion() != "global.awsRegion" ||
		local.awsKey() != "global.awsKey" ||
		local.awsSecret() != "global.awsSecret" ||
		local.grimServerID() != "grimServerID" {
		t.Errorf("local overwrote global. %v", local)
	}
}

func TestValidateLocalEffectiveConfig(t *testing.T) {
	snsTopicName := "foo.go"
	errs := localConfig{local: configMap{"SNSTopicName": snsTopicName}}.errors()
	if len(errs) == 0 {
		t.Errorf("validated with period in name")
	}

	if !strings.Contains(errs[0].Error(), "[ "+snsTopicName+" ]") {
		t.Errorf("error message should mention SNSTopicName %v", errs[0].Error())
	}
}

func TestValidateEffectiveConfig(t *testing.T) {
	checks := []struct {
		gc             globalConfig
		shouldValidate bool
	}{
		{globalConfig{}, false},
		{globalConfig{"AWSRegion": "reg", "AWSKey": "key"}, false},
		{globalConfig{"AWSSecret": "secret", "AWSRegion": "region"}, false},
		{globalConfig{"AWSSecret": "secret", "AWSRegion": "region", "AWSKey": "key"}, true},
	}
	for _, check := range checks {
		errs := check.gc.errors()
		errsEmpty := len(errs) == 0

		if errsEmpty && !check.shouldValidate {
			t.Errorf("invalid config didn't fail validation: %v", check.gc)
		} else if !errsEmpty && check.shouldValidate {
			t.Errorf("valid config failed validation: %v", check.gc)
		}
	}
}

func TestLoadGlobalConfig(t *testing.T) {
	ec, err := readGlobalConfig("./test_data/config_test")
	if err != nil {
		t.Fatalf("|%v|", err)
	}

	if ec.grimServerID() != "def-serverid" {
		t.Errorf("Didn't match:\n%v", ec)
	}

	if ec.successColor() != "green" {
		t.Errorf("Didn't match: \n%v", ec)
	}

	if ec.failureColor() != "red" {
		t.Errorf("Didn't match: \n%v", ec)
	}

	if ec.errorColor() != "gray" {
		t.Errorf("Didn't match: \n%v", ec)
	}

	if ec.pendingColor() != "yellow" {
		t.Errorf("Didn't match: \n%v", ec)
	}
}

func TestLoadRepoConfig(t *testing.T) {
	ec, err := readLocalConfig("./test_data/config_test", "MediaMath", "foo")
	if err != nil {
		t.Fatalf("|%v|", err)
	}

	if ec.pathToCloneIn() != "go/src/github.com/MediaMath/foo" {
		t.Errorf("Didn't match:\n%v", ec)
	}

	if ec.successColor() != "green" {
		t.Errorf("Didn't match: \n%v", ec)
	}

	if ec.failureColor() != "red" {
		t.Errorf("Didn't match: \n%v", ec)
	}

	if ec.errorColor() != "gray" {
		t.Errorf("Didn't match: \n%v", ec)
	}

	if ec.pendingColor() != "yellow" {
		t.Errorf("Didn't match: \n%v", ec)
	}

}

func TestLoadConfig(t *testing.T) {
	ec, err := readGlobalConfig("./test_data/config_test")
	if err != nil {
		t.Fatalf("|%v|", err)
	}

	if ec.hipChatRoom() != "def-hcroom" {
		t.Errorf("Didn't load correctly")
	}
}

func TestBHandMMCanBuildByDefault(t *testing.T) {
	config, err := readLocalConfig("./test_data/config_test", "MediaMath", "foo")
	if err != nil {
		t.Fatalf("|%v|", err)
	}

	if !config.usernameCanBuild("bhand-mm") {
		t.Fatal("bhand-mm should be able to build")
	}
}

func TestBHandMMCanBuild(t *testing.T) {
	config, err := readLocalConfig("./test_data/config_test", "MediaMath", "bar")
	if err != nil {
		t.Fatalf("|%v|", err)
	}

	if !config.usernameCanBuild("bhand-mm") {
		t.Fatal("bhand-mm should be able to build")
	}
}

func TestKKlipschCantBuild(t *testing.T) {
	config, err := readLocalConfig("./test_data/config_test", "MediaMath", "bar")
	if err != nil {
		t.Fatalf("|%v|", err)
	}

	if config.usernameCanBuild("kklipsch") {
		t.Fatal("kklipsch should not be able to build")
	}
}
