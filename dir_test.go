package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var pathsNames = &pathNames{} // a global variable to store file paths of workspace and result

/**
Test on the consistency of timestamp of workspace and result.
*/
func TestOnWorkSpaceAndResultNameConsistency(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "results-dir-consistencyCheck")
	defer os.RemoveAll(tempDir)
	//trigger a build to have file paths of result and workspace
	err := builtForHook(tempDir, "MediaMath", "grim", 0)
	if err != nil {
		t.Fatalf("%v", err)
	}

	isMatched, err := pathsNames.isConsistent()
	if err != nil {
		t.Fatal(err.Error())
	}
	if !isMatched {
		t.Fatalf("inconsistent dir name")
	}
}

type pathNames struct {
	workspacePath string
	resultPath    string
}

func (pn *pathNames) isConsistent() (bool, error) {
	workspacePaths := strings.Split(pn.workspacePath, "/")
	resultPaths := strings.Split(pn.resultPath, "/")
	a := workspacePaths[len(workspacePaths)-1]
	b := resultPaths[len(resultPaths)-1]

	if len(a) == 0 {
		return false, fmt.Errorf("empty workspacePaths name ")
	}

	if len(b) == 0 {
		return false, fmt.Errorf("empty resultPaths name ")
	}

	if len(a) != len(b) || !strings.EqualFold(a, b) {
		return false, fmt.Errorf("inconsistent timestamp workspace:" + a + " and resultpath:" + b)
	}

	return true, nil
}

func builtForHook(tempDir, owner, repo string, exitCode int) error {
	return onHookBuild("not-used", &effectiveConfig{resultRoot: tempDir, workspaceRoot: tempDir}, hookEvent{Owner: owner, Repo: repo}, nil, stubBuild)
}

func stubBuild(configRoot string, resultPath string, config *effectiveConfig, hook hookEvent, basename string) (*executeResult, string, error) {
	pathsNames.resultPath = resultPath
	return built(config.gitHubToken, configRoot, config.workspaceRoot, resultPath, config.pathToCloneIn, hook.Owner, hook.Repo, hook.Ref, hook.env(), basename)
}

func built(token, configRoot, workspaceRoot, resultPath, clonePath, owner, repo, ref string, extraEnv []string, basename string) (*executeResult, string, error) {
	ws := &testWorkSpaceBuilder{"!@#", &executeResult{ExitCode: 0}}
	return grimBuild(ws, resultPath, basename)
}

type testWorkSpaceBuilder struct {
	StubbuildScriptPath string
	StubbuildResult     *executeResult
}

func (tb *testWorkSpaceBuilder) PrepareWorkspace(basename string) (string, error) {
	workSpacePath, err := makeTree(basename)
	pathsNames.workspacePath = workSpacePath
	return workSpacePath, err
}

func (tb *testWorkSpaceBuilder) FindBuildScript(workspacePath string) (string, error) {
	return ".", nil
}

func (tb *testWorkSpaceBuilder) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return &executeResult{ExitCode: 0}, nil
}
