package grim

import (
	"bytes"
	"fmt"
	"text/template"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type grimNotification interface {
	GithubRefStatus() refStatusState
	HipchatNotification(context *grimNotificationContext, config *effectiveConfig) (string, messageColor, error)
}

type standardGrimNotification struct {
	githubState  refStatusState
	hipchatColor messageColor
	getTemplate  func(*effectiveConfig) string
}

var GrimPending = &standardGrimNotification{RSPending, ColorYellow, func(c *effectiveConfig) string { return c.pendingTemplate }}
var GrimError = &standardGrimNotification{RSError, ColorGray, func(c *effectiveConfig) string { return c.errorTemplate }}
var GrimFailure = &standardGrimNotification{RSFailure, ColorRed, func(c *effectiveConfig) string { return c.failureTemplate }}
var GrimSuccess = &standardGrimNotification{RSSuccess, ColorGreen, func(c *effectiveConfig) string { return c.successTemplate }}

func (s *standardGrimNotification) GithubRefStatus() refStatusState {
	return s.githubState
}

func (s *standardGrimNotification) HipchatNotification(context *grimNotificationContext, config *effectiveConfig) (string, messageColor, error) {
	message, err := context.render(s.getTemplate(config))
	return message, s.hipchatColor, err
}

type grimNotificationContext struct {
	Owner     string
	Repo      string
	EventName string
	Target    string
	UserName  string
	Workspace string
	LogDir    string
}

func (c *grimNotificationContext) render(templateString string) (string, error) {
	template, tempErr := template.New("msg").Parse(templateString)
	if tempErr != nil {
		return "", fmt.Errorf("Error parsing notification template: %v", tempErr)
	}

	var doc bytes.Buffer
	if tempErr = template.Execute(&doc, c); tempErr != nil {
		return "", fmt.Errorf("Error applying template: %v", tempErr)
	}

	return doc.String(), nil
}

func buildContext(hook hookEvent, ws, logDir string) *grimNotificationContext {
	return &grimNotificationContext{hook.owner, hook.repo, hook.eventName, hook.target, hook.userName, ws, logDir}
}

func notify(config *effectiveConfig, hook hookEvent, ws string, logDir string, notification grimNotification) error {
	if hook.eventName != "push" && hook.eventName != "pull_request" {
		return nil
	}

	ghErr := setRefStatus(config.gitHubToken, hook.owner, hook.repo, hook.statusRef, notification.GithubRefStatus(), "", "")

	if config.hipChatToken != "" && config.hipChatRoom != "" {
		context := buildContext(hook, ws, logDir)
		message, color, err := notification.HipchatNotification(context, config)
		if err != nil {
			return err
		}

		err = sendMessageToRoom(config.hipChatToken, config.hipChatRoom, config.grimServerID, message, color)
		if err != nil {
			return err
		}
	}

	return ghErr
}
