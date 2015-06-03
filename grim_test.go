package grim

import (
	"encoding/json"
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

	if _, err := resultsDirectoryExists(tempDir, testOwner, testRepo); err != nil {
		t.Errorf("|%v|", err)
	}

}

func TestOnActionError(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-error")
	defer os.RemoveAll(tempDir)

	doNothingAction(tempDir, testOwner, testRepo, 0, fmt.Errorf("Bad Bad thing happened"))

	if _, err := resultsDirectoryExists(tempDir, testOwner, testRepo); err != nil {
		t.Errorf("|%v|", err)
	}
}

func TestResultsDirectoryCreatedInOnHook(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-success")
	defer os.RemoveAll(tempDir)

	doNothingAction(tempDir, testOwner, testRepo, 0, nil)

	if _, err := resultsDirectoryExists(tempDir, testOwner, testRepo); err != nil {
		t.Errorf("|%v|", err)
	}
}

func TestHookGetsLogged(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-success")
	defer os.RemoveAll(tempDir)

	hook := hookEvent{Owner: testOwner, Repo: testRepo, StatusRef: "fooooooooooooooooooo"}

	onHook("not-used", &effectiveConfig{resultRoot: tempDir}, hook, func(r string, resultPath string, c *effectiveConfig, h hookEvent) (*executeResult, string, error) {
		return &executeResult{ExitCode: 0}, "", nil
	})

	results, _ := resultsDirectoryExists(tempDir, testOwner, testRepo)
	hookFile := filepath.Join(results, "hook.json")

	if _, err := os.Stat(hookFile); os.IsNotExist(err) {
		t.Errorf("%s was not created.", hookFile)
	}

	jsonHookFile, readerr := ioutil.ReadFile(hookFile)
	if readerr != nil {
		t.Errorf("Error reading file %v", readerr)
	}

	var parsed hookEvent
	parseErr := json.Unmarshal(jsonHookFile, &parsed)
	if parseErr != nil {
		t.Errorf("Error parsing: %v", parseErr)
	}

	if hook.Owner != parsed.Owner || hook.Repo != parsed.Repo || hook.StatusRef != parsed.StatusRef {
		t.Errorf("Did not match:\n%v\n%v", hook, parsed)
	}

}

func doNothingAction(tempDir, owner, repo string, exitCode int, returnedErr error) error {
	return onHook("not-used", &effectiveConfig{resultRoot: tempDir}, hookEvent{Owner: owner, Repo: repo}, func(r string, resultPath string, c *effectiveConfig, h hookEvent) (*executeResult, string, error) {
		return &executeResult{ExitCode: exitCode}, "", returnedErr
	})
}

func resultsDirectoryExists(tempDir, owner, repo string) (string, error) {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return "", err
	}

	var fileNames []string
	for _, stat := range files {
		fileNames = append(fileNames, stat.Name())
	}

	repoResults := filepath.Join(tempDir, owner, repo)

	if _, err := os.Stat(repoResults); os.IsNotExist(err) {
		return "", fmt.Errorf("%s was not created: %s", repoResults, fileNames)
	}

	baseFiles, err := ioutil.ReadDir(repoResults)
	if len(baseFiles) != 1 {
		return "", fmt.Errorf("Did not create base name in repo results")
	}

	return filepath.Join(repoResults, baseFiles[0].Name()), nil
}
