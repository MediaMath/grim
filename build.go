package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
	"path/filepath"
)

type grimBuilder interface {
	PrepareWorkspace() (string, error)
	FindBuildScript(workspacePath string) (string, error)
	RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error)
}

func (ws *workspaceBuilder) PrepareWorkspace() (string, error) {
	return prepareWorkspace(ws.token, ws.workspaceRoot, ws.clonePath, ws.owner, ws.repo, ws.ref)
}

func (ws *workspaceBuilder) FindBuildScript(workspacePath string) (string, error) {
	configBuildScript := filepath.Join(ws.configRoot, ws.owner, ws.repo, buildScriptName)
	if fileExists(configBuildScript) {
		return configBuildScript, nil
	}

	repoBuildScript := filepath.Join(workspacePath, ws.clonePath, repoBuildScriptName)
	if fileExists(repoBuildScript) {
		return repoBuildScript, nil
	}

	hiddenRepoBuildScript := filepath.Join(workspacePath, ws.clonePath, repoHiddenBuildScriptName)
	if fileExists(hiddenRepoBuildScript) {
		return hiddenRepoBuildScript, nil
	}

	return "", fmt.Errorf("unable to find a build script to run; see README.md for more information")
}

func (ws *workspaceBuilder) RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error) {
	env := os.Environ()
	env = append(env, fmt.Sprintf("CLONE_PATH=%v", ws.clonePath))
	env = append(env, ws.extraEnv...)

	return executeWithOutputChan(outputChan, env, workspacePath, buildScript)
}

type workspaceBuilder struct {
	workspaceRoot string
	clonePath     string
	token         string
	configRoot    string
	owner         string
	repo          string
	ref           string
	extraEnv      []string
}

func grimBuild(builder grimBuilder, resultPath string) (*executeResult, string, error) {

	status, err := buildStatusFile(resultPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create build status file: %v", err)
	}
	defer status.Close()

	workspacePath, err := builder.PrepareWorkspace()
	if err != nil {
		status.WriteString(fmt.Sprintf("failed to prepare workspace %s %v\n", workspacePath, err))
		return nil, workspacePath, fmt.Errorf("failed to prepare workspace: %v", err)
	}
	status.WriteString(fmt.Sprintf("workspace created %s\n", workspacePath))

	buildScriptPath, err := builder.FindBuildScript(workspacePath)
	if err != nil {
		status.WriteString(fmt.Sprintf("%v\n", err))
		return nil, workspacePath, err
	}
	status.WriteString(fmt.Sprintf("build script found %s\n", buildScriptPath))

	outputChan := make(chan string)
	go writeOutput(resultPath, outputChan)

	status.WriteString(fmt.Sprintf("build started ...\n"))
	result, err := builder.RunBuildScript(workspacePath, buildScriptPath, outputChan)
	if err != nil {
		status.WriteString(fmt.Sprintf("build error %v\n", err))
		return nil, workspacePath, err
	}

	if result.ExitCode == 0 {
		status.WriteString(fmt.Sprintf("build success\n"))
		os.RemoveAll(workspacePath)
	} else {
		status.WriteString(fmt.Sprintf("build failed %v\n", result.ExitCode))
	}

	err = appendResult(resultPath, *result)
	if err != nil {
		return result, workspacePath, fatalGrimErrorf("error while storing result: %v", err)
	}

	status.WriteString(fmt.Sprintf("build done\n"))
	return result, workspacePath, nil
}

func build(token, configRoot, workspaceRoot, resultPath, clonePath, owner, repo, ref string, extraEnv []string) (*executeResult, string, error) {
	ws := &workspaceBuilder{workspaceRoot, clonePath, token, configRoot, owner, repo, ref, extraEnv}
	return grimBuild(ws, resultPath)
}
