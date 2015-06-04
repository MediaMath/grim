package grim

import (
	"fmt"
	"testing"
	"io/ioutil"
	"strings"
	"os"
)
// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//status code for generating different versions of workspace generator
const (
	FailWhenPrepareWorkspace string = "failed to prepare workspace  failed to create workspace directory"
	FailWhenFindBuildScript string = "unable to find a build script to run; see README.md for more information"
	FailWhenRunBuildScript string = "build error in running process"
	NoFail string = "no fails after run build script"
	DoneWithPrepareWorkspace string = "PrepareWorkspace...finished"
	DoneWithFindBuildScript string = "FindBuildScript...finished"
	DoneWithRunBuildScript string = "RunBuildScript...finished"
)

//path to write log of workspace builder
var resultPath = "./"

//a workspace builder that has no fail when calling PrepareWorkspace(), FindBuildScript() and RunBuildScript()
type workSpaceBuilderNoFail  struct {
	workspaceBuilder
}

func (ws *workSpaceBuilderNoFail) PrepareWorkspace() (string, error) {
	return DoneWithPrepareWorkspace, nil
}

func (ws *workSpaceBuilderNoFail) FindBuildScript(worksFindBuildScriptpacePath string) (string, error) {
	return DoneWithFindBuildScript, nil
}

func (ws *workSpaceBuilderNoFail) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return &executeResult{ExitCode:0}, nil
}

//a workspace builder that will fail when calling PrepareWorkspace()
type workspaceBuilderFailWhenPrepareWorkSpace struct {
	workSpaceBuilderNoFail
}

func (ws *workspaceBuilderFailWhenPrepareWorkSpace) PrepareWorkspace() (string, error) {
	return "", fmt.Errorf(FailWhenPrepareWorkspace)
}

//a workspace builder that will fail when calling FindBuildScript()
type workspaceBuilderFailWhenFindBuildScript struct {
	workSpaceBuilderNoFail
}

func (ws *workspaceBuilderFailWhenFindBuildScript) FindBuildScript(workspacePath string) (string, error) {
	return "", fmt.Errorf(FailWhenFindBuildScript)

}

//a workspace builder that will fail when calling RunBuildScript()
type workspaceBuilderFailWhenRunBuildScript struct {
	workSpaceBuilderNoFail
}

func (ws *workspaceBuilderFailWhenRunBuildScript) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return nil, fmt.Errorf(FailWhenRunBuildScript)
}

func workspaceBuilderGeneator(whenToFail string) (grimBuilder, error) {
	switch(whenToFail){
	case FailWhenPrepareWorkspace:
		return &workspaceBuilderFailWhenPrepareWorkSpace{workSpaceBuilderNoFail{workspaceBuilder{}}}, nil

	case FailWhenFindBuildScript:
		return &workspaceBuilderFailWhenFindBuildScript{workSpaceBuilderNoFail{workspaceBuilder{}}}, nil

	case FailWhenRunBuildScript:
		return &workspaceBuilderFailWhenRunBuildScript{workSpaceBuilderNoFail{workspaceBuilder{}}}, nil

	case NoFail:
		return &workSpaceBuilderNoFail{workspaceBuilder{}}, nil

	default:
		return nil, fmt.Errorf("failed to generate workspace builder")
	}
}

func TestFailWhenPrepareWorkspace(t *testing.T) {
	os.Remove(resultPath+"/build.txt")
	wb, _ := workspaceBuilderGeneator(FailWhenPrepareWorkspace)
	grimBuild(wb, resultPath)
	content, _ := ioutil.ReadFile(resultPath+"/build.txt")
	if !strings.Contains(strings.TrimSpace(string(content)), FailWhenPrepareWorkspace) {
		t.Error("something wrong when calling PrepareWorkspace(), there should nbe an error here")
	}
}

func TestFailWhenFindBuildScript(t *testing.T) {
	os.RemoveAll(resultPath+"/build.txt")
	wb, _ := workspaceBuilderGeneator(FailWhenFindBuildScript)
	grimBuild(wb, resultPath)
	content, _ := ioutil.ReadFile(resultPath+"/build.txt")

	if !strings.Contains(string(content), DoneWithPrepareWorkspace) {
		t.Error("workspace is not created when calling FindBuildScript()")
	}

	if !strings.Contains(string(content), FailWhenFindBuildScript) {
		t.Error("something wrong when calling FindBuildScript(), there should be an error here")
	}
}

func TestFailWhenWhenRunBuildScript(t *testing.T) {
	os.RemoveAll(resultPath+"/build.txt")
	wb, _ := workspaceBuilderGeneator(FailWhenRunBuildScript)
	grimBuild(wb, resultPath)
	content, _ := ioutil.ReadFile(resultPath+"/build.txt")

	if !strings.Contains(string(content), DoneWithPrepareWorkspace) {
		t.Error("workspace is not created when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), DoneWithFindBuildScript) {
		t.Error("run script is not found when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), "build started ...") {
		t.Error("build is not started when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), FailWhenRunBuildScript) {
		t.Error("something wrong when calling FindBuildScript(), there should be an error here")
	}
}

func TestWhenNoFailsinWorkSpaceBuilder(t *testing.T) {
	os.RemoveAll(resultPath+"/build.txt")
	wb, _ := workspaceBuilderGeneator(NoFail)
	grimBuild(wb, resultPath)
	content, _ := ioutil.ReadFile(resultPath+"/build.txt")

	if !strings.Contains(string(content), DoneWithPrepareWorkspace) {
		t.Error("workspace is not created when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), DoneWithFindBuildScript) {
		t.Error("run script is not found when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), "build started ...") {
		t.Error("build is not started when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), "build success") {
		t.Error("build is not success when calling RunBuildScript()")
	}

	if !strings.Contains(string(content), "build done") {
		t.Error("build is not done when calling RunBuildScript()")
	}

}
