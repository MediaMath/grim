package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
	"path/filepath"
)

func build(token, configRoot, workspaceRoot, resultRoot, clonePath, owner, repo, ref string, extraEnv []string) (*executeResult, string, string, error) {

	workspacePath, err := prepareWorkspace(token, workspaceRoot, clonePath, owner, repo, ref)
	if err != nil {
		return nil, workspacePath, "", fmt.Errorf("failed to prepare workspace: %v", err)
	}

	resultRootPath, err := findLogDirPath(resultRoot, owner, repo)
	if err != nil {
		return nil, "", resultRootPath, fmt.Errorf("failed to get resultRootPath: %v", err)
	}

	buildScriptPath := findBuildScript(configRoot, workspacePath, clonePath, owner, repo)
	if buildScriptPath == "" {
		return nil, workspacePath, "", fmt.Errorf("unable to find a build script to run; see README.md for more information")
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("CLONE_PATH=%v", clonePath))
	env = append(env, extraEnv...)
	result, err := execute(env, workspacePath, buildScriptPath)
	if err != nil {
		return nil, workspacePath, resultRootPath, err
	}

	if result.ExitCode == 0 {
		os.RemoveAll(workspacePath) //TODO: remove resultrootpath
	}

	return result, workspacePath, resultRootPath, nil
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

//Workspace - /var/tmp/grim/MediaMath/Keryx/244159680.
//log directory - /var/log/grim/MediaMath/Keryx/1432776632064635784
func findLogDirPath(resultRootPath, owner, repo string) (string, error) {
	logPath := makeTreeNoCreate(resultRootPath, owner, repo)

	if !fileExistsAndIsDirectory(logPath) {
		return "", nil
	}
	return logPath, nil
}
