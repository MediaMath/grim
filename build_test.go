package grim

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

const testBuildTimeOut = time.Second * time.Duration(10)

func TestPreparePublicRepo(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping prepare test in short mode.")
	}

	f, err := ioutil.TempDir("", "prepare-public-repo-test")
	if err != nil {
		t.Errorf("|%v|", err)
	}
	basename := getTimeStamp()
	clonePath := filepath.Join("foo", "bar", "baz")

	builder := &workspaceBuilder{
		workspaceRoot: f,
		clonePath:     clonePath,
		//public repo shouldnt require token
		token:      "",
		configRoot: "",
		owner:      "MediaMath",
		repo:       "part",
		ref:        "eb78552e86dfead7f6506e6d35ae5db9fc078403",
		extraEnv:   []string{},
		timeOut:    testBuildTimeOut,
	}

	ws, err := builder.PrepareWorkspace(basename)
	if err != nil {
		t.Errorf("|%v|", err)
		t.FailNow()
	}

	files, err := ioutil.ReadDir(filepath.Join(ws, clonePath))
	if err != nil {
		t.Errorf("|%v|", err)
	}

	if len(files) != 16 {
		t.Errorf("Directory %s had %v files.", ws, len(files))
	}

	if !t.Failed() {
		//remove directory only on success so you can diagnose failure
		os.RemoveAll(f)
	}
}

type testBuilder struct {
	workspaceErr    error
	workspaceResult string
	buildScriptErr  error
	buildScriptPath string
	buildErr        error
	buildResult     *executeResult
}

func (tb *testBuilder) PrepareWorkspace(basename string) (string, error) {
	return tb.workspaceResult, tb.workspaceErr
}
func (tb *testBuilder) FindBuildScript(workspacePath string) (string, error) {
	return tb.buildScriptPath, tb.buildScriptErr
}
func (tb *testBuilder) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return tb.buildResult, tb.buildErr
}

func TestOnBuildStatusFileError(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "build-status-file-error")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{workspaceErr: errors.New(""), workspaceResult: "@$@$"}
	grimBuild(tb, resultPath, "")

	_, err := ioutil.ReadFile(resultPath + "/build.txt")
	if err != nil {
		t.Errorf(fmt.Sprintf("Error in building status file: %v", err))
	}
}

func TestOnPrepareWorkspaceFailure(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "prepare-workspace-failure")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{workspaceErr: errors.New(""), workspaceResult: "@$@$"}
	grimBuild(tb, resultPath, "")

	buildFile, _ := ioutil.ReadFile(resultPath + "/build.txt")
	if !strings.Contains(string(buildFile), "failed to prepare workspace @$@$") {
		t.Errorf("Failed to log error in preparing workspace")
	}
}

func TestOnBuildScriptFailure(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "builds-script-failure")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{buildScriptErr: errors.New("&^&^")}
	grimBuild(tb, resultPath, "")

	buildFile, _ := ioutil.ReadFile(resultPath + "/build.txt")
	buildText := string(buildFile)

	if !strings.Contains(buildText, "workspace created") {
		t.Errorf("Failed to log sucessful workspace creation")
	}

	if !strings.Contains(buildText, "&^&^") {
		t.Errorf("Failed to log error in FindBuildScript")
	}
}

func TestOnRunBuildScriptError(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "builds-script-error")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{buildScriptPath: "!@#", buildErr: errors.New("^%$")} //buildResult: &executeResult{ExitCode: 0}}
	grimBuild(tb, resultPath, "")

	buildFile, _ := ioutil.ReadFile(resultPath + "/build.txt")
	buildText := string(buildFile)

	if !strings.Contains(buildText, "workspace created") {
		t.Errorf("Failed to log successful workspace creation")
	}

	if !strings.Contains(buildText, "build script found !@#") {
		t.Errorf("Failed to log successful build start")
	}
	if !strings.Contains(buildText, "build started ...") {
		t.Errorf("Failed to log build start")
	}

	if !strings.Contains(buildText, "build error ^%$") {
		t.Errorf("Failed to log build error")
	}
}

func TestOnRunBuildScriptSuccess(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "build-script-success")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{buildScriptPath: "!@#", buildResult: &executeResult{ExitCode: 0}}
	grimBuild(tb, resultPath, "")

	buildFile, _ := ioutil.ReadFile(resultPath + "/build.txt")
	buildText := string(buildFile)

	if !strings.Contains(buildText, "workspace created") {
		t.Errorf("Failed to log successful workspace creation")
	}

	if !strings.Contains(buildText, "build script found !@#") {
		t.Errorf("Failed to log successful build start")
	}
	if !strings.Contains(buildText, "build started ...") {
		t.Errorf("Failed to log build start")
	}

	if !strings.Contains(buildText, "build success") {
		t.Errorf("Failed to log build success")
	}
}

func TestOnRunBuildScriptFailure(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "build-script-error")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{buildScriptPath: "!@#", buildResult: &executeResult{ExitCode: 123123123}}
	grimBuild(tb, resultPath, "")

	buildFile, _ := ioutil.ReadFile(resultPath + "/build.txt")
	buildText := string(buildFile)

	if !strings.Contains(buildText, "workspace created") {
		t.Errorf("Failed to log successful workspace creation")
	}

	if !strings.Contains(buildText, "build script found !@#") {
		t.Errorf("Failed to log successful build start")
	}
	if !strings.Contains(buildText, "build started ...") {
		t.Errorf("Failed to log build start")
	}

	if !strings.Contains(buildText, "build failed") || !strings.Contains(buildText, "123123123") {
		t.Errorf("Failed to log build failure")
	}
}
