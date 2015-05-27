package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io/ioutil"
)

func prepareWorkspace(token, workspaceRoot, clonePath, owner, repo, ref string) (string, error) {
	workspacePath, err := createWorkspaceDirectory(workspaceRoot, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %v", err)
	}

	_, err = cloneRepo(token, workspacePath, clonePath, owner, repo, ref)
	if err != nil {
		return "", fmt.Errorf("failed to download repo archive: %v", err)
	}

	return workspacePath, nil
}

func createWorkspaceDirectory(workspaceRoot, owner, repo string) (string, error) {
	workspaceParent := makeTree(workspaceRoot, owner, repo)

	workspacePath, err := ioutil.TempDir(workspaceParent, "")
	if err != nil {
		return "", err
	}

	return workspacePath, nil
}
