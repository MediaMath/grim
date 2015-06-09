package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "fmt"

func prepareWorkspace(token, workspaceRoot, clonePath, owner, repo, ref, basename string) (string, error) {
	return inWorkspaceDirectory(workspaceRoot, owner, repo, basename, func(workspacePath string) error {
		_, err := cloneRepo(token, workspacePath, clonePath, owner, repo, ref)
		return err
	})
}

func inWorkspaceDirectory(root, owner, repo, basename string, action func(string) error) (string, error) {
	workspacePath, err := makeTree(root, owner, repo, basename)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %v", err)
	}

	if err := action(workspacePath); err != nil {
		return "", fmt.Errorf("failed to download repo archive: %v", err)
	}

	return workspacePath, nil
}
