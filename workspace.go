package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func prepareWorkspace(configRoot, workspaceRoot, owner, repo string, env []string) (string, error) {
	prepareScriptPath := findPrepareScript(configRoot, owner, repo)

	workspacePath, err := createWorkspaceDirectory(workspaceRoot, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %v", err)
	}

	if prepareScriptPath != "" {
		result, err := execute(prepareScriptPath, workspacePath, env)
		if err != nil {
			return "", fmt.Errorf("prepare script exited with err: %v", err)
		}

		if result.ExitCode != 0 {
			return "", fmt.Errorf("prepare script exited with code: %v", result.ExitCode)
		}
	}

	return workspacePath, nil
}

func findPrepareScript(configRoot, owner, repo string) string {
	repoPrepareScript := filepath.Join(configRoot, owner, repo, prepareScriptName)
	if fileExists(repoPrepareScript) {
		return repoPrepareScript
	}

	globalPrepareScript := filepath.Join(configRoot, prepareScriptName)
	if fileExists(globalPrepareScript) {
		return globalPrepareScript
	}

	return ""
}

func createWorkspaceDirectory(workspaceRoot, owner, repo string) (string, error) {
	workspaceParent := makeTree(workspaceRoot, owner, repo)

	workspacePath, err := ioutil.TempDir(workspaceParent, "")
	if err != nil {
		return "", err
	}

	return workspacePath, nil
}
