package grim

import (
	"fmt"
	"testing"
	"log"
	"strings"
	"bytes"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
var testContext = &grimNotificationContext{
	Owner:     "rain",
	Repo:      "spain",
	EventName: "falls",
	Target:    "plain",
	UserName:  "mainly",
	Workspace: "boogey/nights",
	LogDir:    "once/again/where/it/rains",
}

var testConfig = &effectiveConfig{
	pendingTemplate: "pending {{.Owner}}",
	errorTemplate:   "error {{.Repo}}",
	failureTemplate: "failure {{.Target}}",
	successTemplate: "success {{.UserName}}"}

var testHook = hookEvent{
	Owner: "MediaMath",
	Repo: "grim",
	EventName: "push",
}

func TestLoggingHipChatIncompleteSetup(t *testing.T){
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Lshortfile)

	notify(testConfig, testHook, "", GrimPending, logger)
	message := fmt.Sprintf("%v", &buf)

	if !strings.Contains(message, "pending MediaMath") {
		t.Errorf("Failed to log message")
	}

	if !strings.Contains(message, "HipChat: config.hipChatToken and config.hitChatRoom not set") {
		t.Errorf("Failed to log that token and room from config are not set")
	}
}

func TestLoggingHipChatErrorCreatingMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Lshortfile)

	testConfigWithHC := &effectiveConfig{
		pendingTemplate: "pending {{.NOPE}}",
		hipChatToken: "NOT_EMPTY",
		hipChatRoom: "NON_EMPTY",
	}

	notify(testConfigWithHC, testHook, "", GrimPending, logger)
	message := fmt.Sprintf("%v", &buf)

	if !strings.Contains(message, "Hipchat: Error while rendering message") {
		t.Errorf("Failed to log error in creating to room")
	}
}

func TestLoggingHipChatErrorSendingMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Lshortfile)

	testConfigWithHC := &effectiveConfig{
		pendingTemplate: "pending {{.Owner}}",
		hipChatToken: "NOT_EMPTY",
		hipChatRoom: "NON_EMPTY",
	}

	notify(testConfigWithHC, testHook, "", GrimPending, logger)
	message := fmt.Sprintf("%v", &buf)

	if !strings.Contains(message, "pending MediaMath") {
		t.Errorf("Failed to log message")
	}

	if !strings.Contains(message, "Hipchat: Error while sending message to room") {
		t.Errorf("Failed to log error in sending to room")
	}
}

func TestPending(t *testing.T) {
	if err := compareNotification(GrimPending, RSPending, ColorYellow, "pending rain"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestError(t *testing.T) {
	if err := compareNotification(GrimError, RSError, ColorGray, "error spain"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestFailure(t *testing.T) {
	if err := compareNotification(GrimFailure, RSFailure, ColorRed, "failure plain"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestSuccess(t *testing.T) {
	if err := compareNotification(GrimSuccess, RSSuccess, ColorGreen, "success mainly"); err != nil {
		t.Errorf("%v", err)
	}
}

func compareNotification(n *standardGrimNotification, state refStatusState, color messageColor, message string) error {
	if n.GithubRefStatus() != state {
		return fmt.Errorf("Github: %v", n)
	}

	msg, color, err := n.HipchatNotification(testContext, testConfig)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}

	if msg != message || color != color {
		return fmt.Errorf("%v %v", message, color)
	}

	return nil
}

func TestContextRender(t *testing.T) {

	str, err := testContext.render("The {{.Owner}} in {{.Repo}} {{.EventName}} {{.UserName}} on the {{.Target}} {{.Workspace}}!")
	errStr, err := testContext.render("The {{.Owner}} in {{.Repo}} {{.EventName}} {{.UserName}} on the {{.Target}} {{.LogDir}}!")

	if err != nil {
		t.Errorf("error %v", err)
	}

	if str != "The rain in spain falls mainly on the plain boogey/nights!" {
		t.Errorf("Didn't match %v", str)
	}

	if errStr != "The rain in spain falls mainly on the plain once/again/where/it/rains!" {
		t.Errorf("Didn't match %v", errStr)
	}
}
