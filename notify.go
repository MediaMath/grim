package grim

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"text/template"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type grimNotification interface {
	GithubRefStatus() refStatusState
	HipchatNotification(context *grimNotificationContext, config localConfig) (string, string, error)
}

type standardGrimNotification struct {
	githubState     refStatusState
	getHipchatColor func(localConfig) string
	getTemplate     func(localConfig) string
}

//GrimPending is the notification used for pending builds.
var GrimPending = &standardGrimNotification{
	RSPending,
	func(c localConfig) string { return c.pendingColor() },
	func(c localConfig) string { return c.pendingTemplate() },
}

//GrimError is the notification used for builds that cannot be run correctly.
var GrimError = &standardGrimNotification{
	RSError,
	func(c localConfig) string { return c.errorColor() },
	func(c localConfig) string { return c.errorTemplate() },
}

//GrimFailure is the notification used when builds fail.
var GrimFailure = &standardGrimNotification{
	RSFailure,
	func(c localConfig) string { return c.failureColor() },
	func(c localConfig) string { return c.failureTemplate() },
}

//GrimSuccess is the notification used when builds succeed.
var GrimSuccess = &standardGrimNotification{
	RSSuccess,
	func(c localConfig) string { return c.successColor() },
	func(c localConfig) string { return c.successTemplate() },
}

func (s *standardGrimNotification) GithubRefStatus() refStatusState {
	return s.githubState
}

func (s *standardGrimNotification) HipchatNotification(context *grimNotificationContext, config localConfig) (string, string, error) {
	message, err := context.render(s.getTemplate(config))
	return message, s.getHipchatColor(config), err
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
	return &grimNotificationContext{hook.Owner, hook.Repo, hook.EventName, hook.Target, hook.UserName, ws, logDir}
}

func notify(config localConfig, hook hookEvent, ws, logDir string, notification grimNotification, logger *log.Logger) error {
	if hook.EventName != "push" && hook.EventName != "pull_request" {
		return nil
	}

	ghErr := setRefStatus(config.gitHubToken(), hook.Owner, hook.Repo, hook.StatusRef, notification.GithubRefStatus(), "", "")

	context := buildContext(hook, ws, logDir)
	message, color, err := notification.HipchatNotification(context, config)
	logger.Print(message)

	if config.hipChatToken() != "" && config.hipChatRoom() != "" {
		if err != nil {
			logger.Printf("Hipchat: Error while rendering message: %v", err)
			return err
		}

		switch config.hipChatVersion() {
		case 1:
			err = sendMessageToRoom(config.hipChatToken(), config.hipChatRoom(), config.grimServerID(), message, color)
		case 2:
			err = sendMessageToRoom2(config.hipChatToken(), config.hipChatRoom(), config.grimServerID(), message, color)
		default:
			err = errors.New("invalid or unsupported hipchat version")
		}

		if err != nil {
			logger.Printf("Hipchat: Error while sending message to room: %v", err)
			return err
		}
	} else {
		logger.Print("HipChat: config.hipChatToken and config.hitChatRoom not set")
	}

	return ghErr
}
