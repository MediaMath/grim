package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func cloneRepo(token, workspacePath, clonePath, owner, repo, ref string) (string, error) {
	archive, err := downloadRepo(token, owner, repo, ref, workspacePath)
	if err != nil {
		return "", err
	}

	// var ws = new(workspaceBuilder)
	// conf, err := getEffectiveConfig(ws.configRoot, ws.owner, ws.repo)
	// if err != nil {
	// 	return "", err
	// }
	return unarchiveRepo(archive, workspacePath, clonePath)
}

func unarchiveRepo(file, workspacePath, clonePath string) (string, error) {
	tarPath, err := exec.LookPath("tar")
	if err != nil {
		return "", err
	}

	finalName := filepath.Join(workspacePath, clonePath)
	if mkErr := os.MkdirAll(finalName, 0700); mkErr != nil {
		return "", fmt.Errorf("Could not make path %s: %v", finalName, mkErr)
	}

	//extracts the folder into the finalName directory pulling off the top level folder
	//will break if github starts returning a different tar format
	result, err := execute(nil, workspacePath, tarPath, "-xvf", file, "-C", finalName, "--strip-components=1")

	if err != nil {
		return "", err
	}

	if result.ExitCode != 0 {
		return "", fmt.Errorf("extract archive failed: %v %v", result.ExitCode, strings.TrimSpace(result.Output))
	}

	return finalName, nil
}

func downloadRepo(token, owner, repo, ref string, location string) (string, error) {
	client, err := getClientForToken(token)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("repos/%s/%s/tarball/%s", owner, repo, ref)
	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	temp, err := ioutil.TempFile(location, "download")
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req, temp)
	temp.Close()
	if err != nil {
		return "", err
	}
	_, params, err := mime.ParseMediaType(resp.Response.Header["Content-Disposition"][0])
	if err != nil {
		return "", err
	}

	fileName := params["filename"]
	downloaded := filepath.Join(location, fileName)

	os.Rename(temp.Name(), downloaded)

	return downloaded, nil
}
