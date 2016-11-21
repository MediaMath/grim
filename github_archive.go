package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func cloneRepo(token string, workspacePath string, clonePath string, owner string, repo string, ref string, timeOut time.Duration) (string, error) {
	log.Printf("downloading repo %v/%v@%v to %v with timeout %v", owner, repo, ref, workspacePath, timeOut)
	archive, err := downloadRepo(token, owner, repo, ref, workspacePath)
	if err != nil {
		return "", err
	}

	log.Printf("download of repo %v/%v@%v to %v complete", owner, repo, ref, workspacePath)

	log.Printf("unarchiving repo from %v to %v with timeout %v", archive, filepath.Join(workspacePath, clonePath), timeOut)
	finalPath, err := unarchiveRepo(archive, workspacePath, clonePath, timeOut)
	if err == nil {
		log.Printf("unarchiving of repo from %v to %v complete", archive, finalPath)
	}

	return finalPath, err
}

func unarchiveRepo(file, workspacePath, clonePath string, timeOut time.Duration) (string, error) {
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
	result, err := execute(nil, workspacePath, tarPath, timeOut, "-xvf", file, "-C", finalName, "--strip-components=1")

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
