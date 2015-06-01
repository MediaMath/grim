package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
	"path/filepath"
)

func build(token, configRoot, workspaceRoot, resultRoot, clonePath, owner, repo, ref string, extraEnv []string) (*executeResult, string, error) {

	workspacePath, err := prepareWorkspace(token, workspaceRoot, clonePath, owner, repo, ref)
	if err != nil {
		return nil, workspacePath, fmt.Errorf("failed to prepare workspace: %v", err)
	}

	buildScriptPath := findBuildScript(configRoot, workspacePath, clonePath, owner, repo)
	if buildScriptPath == "" {
		return nil, workspacePath, fmt.Errorf("unable to find a build script to run; see README.md for more information")
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("CLONE_PATH=%v", clonePath))
	env = append(env, extraEnv...)
	result, err := execute(env, workspacePath, buildScriptPath)
	if err != nil {
		return nil, workspacePath, err
	}

	if result.ExitCode == 0 {
		os.RemoveAll(workspacePath)
	}

	return result, workspacePath, nil
}

func findBuildScript(configRoot, workspacePath, clonePath, owner, repo string) string {
	configBuildScript := filepath.Join(configRoot, owner, repo, buildScriptName)
	if fileExists(configBuildScript) {
		return configBuildScript
	}

	repoBuildScript := filepath.Join(workspacePath, clonePath, repoBuildScriptName)
	if fileExists(repoBuildScript) {
		return repoBuildScript
	}

	hiddenRepoBuildScript := filepath.Join(workspacePath, clonePath, repoHiddenBuildScriptName)
	if fileExists(hiddenRepoBuildScript) {
		return hiddenRepoBuildScript
	}

	return ""
}
