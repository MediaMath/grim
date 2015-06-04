package grim

import (
//"encoding/json"
//"io/ioutil"
//"os"
//"path/filepath"
	"fmt"
	"testing"
	"io/ioutil"
	"strings"
	"os"
	"time"
)
// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//status code for generating different versions of workspace generator
const (
	FailWhenPrepareWorkspace string = "failed to prepare workspace  failed to create workspace directory"
	FailWhenFindBuildScript string = "unable to find a build script to run; see README.md for more information"
	FailWhenRunBuildScript string = "build error in running process"
	NoFail = "no fails after run build script"
	DoneWithPrepareWorkspace string = "PrepareWorkspace...finished"
	DoneWithFindBuildScript string = "PrepareWorkspace...finished"
	DoneWithRunBuildScript string = "PrepareWorkspace...finished"
)

//path to write log of workspace builder
var resultPath = "./"
//a full functional  workspace builder
var workSpaceBuilderWhenNoFunctionFails grimBuilder = &WorkSpaceBuilderWhenNoFunctionFails{workspaceBuilder{"", "", "", "", "", "", "", nil}}

type WorkSpaceBuilderWhenNoFunctionFails  struct {
	workspaceBuilder
}

func (ws *WorkSpaceBuilderWhenNoFunctionFails) PrepareWorkspace() (string, error) {
	return DoneWithPrepareWorkspace, nil
}

func (ws *WorkSpaceBuilderWhenNoFunctionFails) FindBuildScript(workspacePath string) (string, error) {
	return DoneWithFindBuildScript, nil
}

func (ws *WorkSpaceBuilderWhenNoFunctionFails) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return &executeResult{time.Now(), time.Now(), time.Nanosecond, time.Nanosecond, nil, 0, ""}, nil
}

//a workspace builder which will fail in the preparing stage
type workspaceBuilderFailWhenPrepareWorkSpace struct {
	workspaceBuilder
}

func (ws *workspaceBuilderFailWhenPrepareWorkSpace) PrepareWorkspace() (string, error) {
	return "", fmt.Errorf(FailWhenPrepareWorkspace)
}

func (ws *workspaceBuilderFailWhenPrepareWorkSpace) FindBuildScript(workspacePath string) (string, error) {
	return workSpaceBuilderWhenNoFunctionFails.FindBuildScript(workspacePath)
}

func (ws *workspaceBuilderFailWhenPrepareWorkSpace) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return workSpaceBuilderWhenNoFunctionFails.RunBuildScript(workspacePath, buildScript, outputChan)
}
//a workspace builder which will fail when finding build script
type workspaceBuilderFailWhenFindBuildScript struct {
	workspaceBuilder
}

func (ws *workspaceBuilderFailWhenFindBuildScript) PrepareWorkspace() (string, error) {
	return workSpaceBuilderWhenNoFunctionFails.PrepareWorkspace()
}

func (ws *workspaceBuilderFailWhenFindBuildScript) FindBuildScript(workspacePath string) (string, error) {
	return "", fmt.Errorf(FailWhenFindBuildScript)

}

func (ws *workspaceBuilderFailWhenFindBuildScript) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return workSpaceBuilderWhenNoFunctionFails.RunBuildScript(workspacePath, buildScript, outputChan)
}
//a workspace builder which will fail when running script
type workspaceBuilderFailWhenRunBuildScript struct {
	workspaceBuilder
}

func (ws *workspaceBuilderFailWhenRunBuildScript) PrepareWorkspace() (string, error) {
	return workSpaceBuilderWhenNoFunctionFails.PrepareWorkspace()
}

func (ws *workspaceBuilderFailWhenRunBuildScript) FindBuildScript(workspacePath string) (string, error) {
	return workSpaceBuilderWhenNoFunctionFails.FindBuildScript(workspacePath)
}

func (ws *workspaceBuilderFailWhenRunBuildScript) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	return nil, fmt.Errorf(FailWhenRunBuildScript)
}

func workspaceBuilderGeneator(whenToFail string) (grimBuilder, error) {
	switch(whenToFail){
	case FailWhenPrepareWorkspace:
		return &workspaceBuilderFailWhenPrepareWorkSpace{workspaceBuilder{"", "", "", "", "", "", "", nil}}, nil

	case FailWhenFindBuildScript:
		return &workspaceBuilderFailWhenFindBuildScript{workspaceBuilder{"", "", "", "", "", "", "", nil}}, nil

	case FailWhenRunBuildScript:
		return &workspaceBuilderFailWhenRunBuildScript{workspaceBuilder{"", "", "", "", "", "", "", nil}}, nil
	case NoFail:
		return workSpaceBuilderWhenNoFunctionFails, nil
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

	if !strings.Contains(string(content), DoneWithRunBuildScript) {
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

func TestWhenWhenNoFailsinWorkSpaceBuilder(t *testing.T) {
	os.RemoveAll(resultPath+"/build.txt")
	wb, _ := workspaceBuilderGeneator(NoFail)
	grimBuild(wb, resultPath)
	content, _ := ioutil.ReadFile(resultPath+"/build.txt")

	if !strings.Contains(string(content), DoneWithRunBuildScript) {
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
