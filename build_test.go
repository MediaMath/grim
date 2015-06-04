package grim

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
type testBuilder struct {
	workspaceErr    error
	workspaceResult string
	buildScriptErr  error
	buildScriptPath string
	buildErr        error
	buildResult     *executeResult
}

func (tb *testBuilder) PrepareWorkspace() (string, error) {
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
	grimBuild(tb, resultPath)

	_, err := ioutil.ReadFile(resultPath + "/build.txt")
	if err != nil {
		t.Errorf(fmt.Sprintf("Error in building status file: %v", err))
	}
}

func TestOnPrepareWorkspaceFailure(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "prepare-workspace-failure")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{workspaceErr: errors.New(""), workspaceResult: "@$@$"}
	grimBuild(tb, resultPath)

	buildFile, _ := ioutil.ReadFile(resultPath + "/build.txt")
	if !strings.Contains(string(buildFile), "failed to prepare workspace @$@$") {
		t.Errorf("Failed to log error in preparing workspace")
	}
}

func TestOnBuildScriptFailure(t *testing.T) {
	resultPath, _ := ioutil.TempDir("", "builds-script-failure")
	defer os.RemoveAll(resultPath)

	tb := &testBuilder{buildScriptErr: errors.New("&^&^")}
	grimBuild(tb, resultPath)

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
	grimBuild(tb, resultPath)

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
	grimBuild(tb, resultPath)

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
	grimBuild(tb, resultPath)

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
