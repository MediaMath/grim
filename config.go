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
	defaultColorForSuccess    = colorForSuccess()
	defaultColorForFailure    = colorForFailure()
	defaultColorForError      = colorForError()
	defaultColorForPending    = colorForPending()
	defaultTemplateForFailure = templateForFailureandError("Failure during")
	defaultHipChatVersion     = 1
)

type configMap map[string]interface{}

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

func colorForSuccess() *string {
	c := string(ColorGreen)
	return &c
}

func colorForFailure() *string {
	c := string(ColorRed)
	return &c
}

func colorForError() *string {
	c := string(ColorGray)
	return &c
}

func colorForPending() *string {
	c := string(ColorYellow)
	return &c
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
	f, _ := val.(float64)
	i := int(f)

	if i == 0 {
		for _, i = range ints {
			if i != 0 {
				break
			}
		}
	}

	return i
}
