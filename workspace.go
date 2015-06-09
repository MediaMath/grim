package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"path/filepath"
)

func prepareWorkspace(token, workspaceRoot, clonePath, owner, repo, ref, basename string) (string, error) {
	workspacePath, err := createWorkspaceDirectory(workspaceRoot, owner, repo, basename)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %v", err)
	}

	_, err = cloneRepo(token, workspacePath, clonePath, owner, repo, ref)
	if err != nil {
		return "", fmt.Errorf("failed to download repo archive: %v", err)
	}

	return workspacePath, nil
}

func createWorkspaceDirectory(workspaceRoot, owner, repo, basename string) (string, error) {
	workspaceParent, err := makeTree(workspaceRoot, owner, repo)
	if err != nil {
		return "", err
	}

	workspacePath := filepath.Join(workspaceParent, basename)
	return workspacePath, nil
}
