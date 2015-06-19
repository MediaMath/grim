package grim

import (
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"bytes"
	"strings"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func TestTruncatedGrimServerID(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Lshortfile)

	tempDir, err := ioutil.TempDir("", "TestTruncatedGrimServerID")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	defer os.RemoveAll(tempDir)

	// GrimQueueName is set to be 20 chars long, which should get truncated
	configJS := `{"GrimQueueName":"12345678901234567890","AWSRegion":"empty","AWSKey":"empty","AWSSecret":"empty"}`
	ioutil.WriteFile(filepath.Join(tempDir, "config.json"), []byte(configJS), 0644)

	g := &Instance{
		configRoot : &tempDir,
		queue      : nil,
	}

	g.PrepareGrimQueue(logger)
	message := fmt.Sprintf("%v", &buf)

	if !strings.Contains(message, buildTruncatedMessage("GrimQueueName")) {
		t.Errorf("Failed to log truncation of grimServerID")
	}
}

func TestTimeOutConfig(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping prepare test in short mode.")
	}

	tempDir, err := ioutil.TempDir("", "TestTimeOut")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	defer os.RemoveAll(tempDir)

	configJS := `{"Timeout":4,"AWSRegion":"empty","AWSKey":"empty","AWSSecret":"empty"}`
	ioutil.WriteFile(filepath.Join(tempDir, "config.json"), []byte(configJS), 0644)

	config, err := getEffectiveGlobalConfig(tempDir)
	if err != nil {
		t.Errorf("|%v|", err)
	}
	config.resultRoot = tempDir
	if config.timeout == defaultTimeOutSeconds {
		t.Errorf("Failed to use non default timeout time")
	}

	err = doWaitAction(config, testOwner, testRepo, 10)
	errMess := fmt.Sprintf("%v",err)
	if errMess != "Build Timeout" {
		t.Errorf("Failed to timeout")
	}

	err = doWaitAction(config, testOwner, testRepo, 2)
	if err != nil {
		t.Errorf("Failed to not timeout")
	}
}

func doWaitAction(config *effectiveConfig, owner, repo string, wait int) error {
	return onHook("not-used", config, hookEvent{Owner: owner, Repo: repo}, nil, func(r string, resultPath string, c *effectiveConfig, h hookEvent, s string) (*executeResult, string, error) {
		time.Sleep(time.Duration(wait) * time.Second)
		return &executeResult{ExitCode: 0}, "", nil
	})
}

func TestBuildRef(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping prepare test in short mode.")
	}

	owner := "MediaMath"
	repo := "grim"
	ref := "test" //special grim branch
	clonePath := "go/src/github.com/MediaMath/grim"

	temp, _ := ioutil.TempDir("", "TestBuildRef")

	configRoot := filepath.Join(temp, "config")
	os.MkdirAll(filepath.Join(configRoot, owner, repo), 0700)

	rr := filepath.Join(temp, "results")
	ws := filepath.Join(temp, "ws")
	bogus := "bogus"

	grimConfig := &config{
		ResultRoot:    &rr,
		WorkspaceRoot: &ws,
		AWSRegion:     &bogus,
		AWSKey:        &bogus,
		AWSSecret:     &bogus,
	}

	configJs, _ := json.Marshal(grimConfig)
	ioutil.WriteFile(filepath.Join(configRoot, "config.json"), configJs, 0644)

	localConfig := &config{PathToCloneIn: &clonePath}
	localJs, _ := json.Marshal(localConfig)

	ioutil.WriteFile(filepath.Join(configRoot, owner, repo, "config.json"), localJs, 0644)
	var g Instance
	g.SetConfigRoot(configRoot)

	logfile, err := os.OpenFile(filepath.Join(temp, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}

	logger := log.New(logfile, "", log.Ldate|log.Ltime)
	buildErr := g.BuildRef(owner, repo, ref, logger)
	logfile.Close()

	if buildErr != nil {
		t.Errorf("%v: %v", temp, buildErr)
	}

	if !t.Failed() {
		os.RemoveAll(temp)
	}
}

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

	onHook("not-used", &effectiveConfig{resultRoot: tempDir}, hook, nil, func(r string, resultPath string, c *effectiveConfig, h hookEvent, s string) (*executeResult, string, error) {
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
	return onHook("not-used", &effectiveConfig{resultRoot: tempDir}, hookEvent{Owner: owner, Repo: repo}, nil, func(r string, resultPath string, c *effectiveConfig, h hookEvent, s string) (*executeResult, string, error) {
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
