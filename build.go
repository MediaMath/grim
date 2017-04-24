package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type grimBuilder interface {
	PrepareWorkspace(basename string) (string, error)
	FindBuildScript(workspacePath string) (string, error)
	RunBuildScript(workspacePath, buildScript string, outputChan chan string) (*executeResult, error)
}

func (ws *workspaceBuilder) PrepareWorkspace(basename string) (string, error) {
	workspacePath, err := makeTree(ws.workspaceRoot, ws.owner, ws.repo, basename)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %v", err)
	}

	_, err = cloneRepo(ws.token, workspacePath, ws.clonePath, ws.owner, ws.repo, ws.ref, ws.timeout)
	if err != nil {
		return "", fmt.Errorf("failed to download repo archive: %v", err)
	}

	return workspacePath, nil
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

	return executeWithOutputChan(outputChan, env, workspacePath, buildScript, ws.timeout)
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
	timeout       time.Duration
}

func grimBuild(builder grimBuilder, resultPath, basename string) (*executeResult, string, error) {

	status, err := buildStatusFile(resultPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create build status file: %v", err)
	}
	defer status.Close()

	statusLogger := log.New(status, "", log.Ldate|log.Ltime)

	workspacePath, err := builder.PrepareWorkspace(basename)
	if err != nil {
		statusLogger.Printf("failed to prepare workspace %s %v\n", workspacePath, err)
		return nil, workspacePath, fmt.Errorf("failed to prepare workspace: %v", err)
	}
	statusLogger.Printf("workspace created %s\n", workspacePath)

	buildScriptPath, err := builder.FindBuildScript(workspacePath)
	if err != nil {
		statusLogger.Printf("%v\n", err)
		return nil, workspacePath, err
	}
	statusLogger.Printf("build script found %s\n", buildScriptPath)

	outputChan := make(chan string)
	go writeOutput(resultPath, outputChan)

	statusLogger.Println("build started ...")
	result, err := builder.RunBuildScript(workspacePath, buildScriptPath, outputChan)
	if err != nil {
		statusLogger.Printf("build error %v\n", err)
		return nil, workspacePath, err
	}

	if result.ExitCode == 0 {
		statusLogger.Println("build success")
		os.RemoveAll(workspacePath)
	} else {
		statusLogger.Printf("build failed %v\n", result.ExitCode)
	}

	err = appendResult(resultPath, *result)
	if err != nil {
		return result, workspacePath, fatalGrimErrorf("error while storing result: %v", err)
	}

	statusLogger.Println("build done")
	return result, workspacePath, nil
}

func build(token, configRoot, workspaceRoot, resultPath, clonePath, owner, repo, ref string, extraEnv []string, basename string, timeout time.Duration) (*executeResult, string, error) {
	ws := &workspaceBuilder{workspaceRoot, clonePath, token, configRoot, owner, repo, ref, extraEnv, timeout}
	return grimBuild(ws, resultPath, basename)
}
