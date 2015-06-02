package grim

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
var testOwner = "MediaMath"
var testRepo = "grim"

func TestOnActionFailure(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-failure")
	defer os.RemoveAll(tempDir)

	doNothingAction(tempDir, testOwner, testRepo, 123, nil)

	if err := resultsDirectoryExists(tempDir, testOwner, testRepo); err != nil {
		t.Errorf("|%v|", err)
	}

}

func TestOnActionError(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-error")
	defer os.RemoveAll(tempDir)

	doNothingAction(tempDir, testOwner, testRepo, 0, fmt.Errorf("Bad Bad thing happened"))

	if err := resultsDirectoryExists(tempDir, testOwner, testRepo); err != nil {
		t.Errorf("|%v|", err)
	}
}

func TestResultsDirectoryCreatedInOnHook(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-success")
	defer os.RemoveAll(tempDir)

	doNothingAction(tempDir, testOwner, testRepo, 0, nil)

	if err := resultsDirectoryExists(tempDir, testOwner, testRepo); err != nil {
		t.Errorf("|%v|", err)
	}
}

func doNothingAction(tempDir, owner, repo string, exitCode int, returnedErr error) error {
	return onHook("not-used", &effectiveConfig{resultRoot: tempDir}, hookEvent{owner: owner, repo: repo}, func(r string, resultPath string, c *effectiveConfig, h hookEvent) (*executeResult, string, error) {
		return &executeResult{ExitCode: exitCode}, "", returnedErr
	})
}

func resultsDirectoryExists(tempDir, owner, repo string) error {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return err
	}

	var fileNames []string
	for _, stat := range files {
		fileNames = append(fileNames, stat.Name())
	}

	repoResults := filepath.Join(tempDir, owner, repo)

	if _, err := os.Stat(repoResults); os.IsNotExist(err) {
		return fmt.Errorf("%s was not created: %s", repoResults, fileNames)
	}

	baseFiles, err := ioutil.ReadDir(repoResults)
	if len(baseFiles) != 1 {
		return fmt.Errorf("Did not create base name in repo results")
	}

	return nil
}
