//Package grim is the "GitHub Responder In MediaMath". We liked the acronym and awkwardly filled in the details to fit it. In short, it is a task runner that is triggered by GitHub push/pull request hooks that is intended as a much simpler and easy-to-use build server than the more modular alternatives (eg. Jenkins).
//grim provides the library functions to support this use case.
//grimd is a daemon process that uses the grim library.
package grim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
	"log"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Instance models the state of a configured Grim instance.
type Instance struct {
	configRoot *string
	queue      *sqsQueue
}

// SetConfigRoot sets the base path of the configuration directory and clears any previously read config values from memory.
func (i *Instance) SetConfigRoot(path string) {
	i.configRoot = &path
	i.queue = nil
}

// PrepareGrimQueue creates or reuses the Amazon SQS queue named in the config.
func (i *Instance) PrepareGrimQueue() error {
	configRoot := getEffectiveConfigRoot(i.configRoot)

	config, err := getEffectiveGlobalConfig(configRoot)
	if err != nil {
		return fatalGrimErrorf("error while reading config: %v", err)
	}

	queue, err := prepareSQSQueue(config.awsKey, config.awsSecret, config.awsRegion, config.grimQueueName)
	if err != nil {
		return fatalGrimErrorf("error preparing queue: %v", err)
	}

	i.queue = queue

	return nil
}

// PrepareRepos discovers all repos that are configured then sets up SNS and GitHub.
// It is an error to call this without calling PrepareGrimQueue first.
func (i *Instance) PrepareRepos() error {
	if err := i.checkGrimQueue(); err != nil {
		return err
	}

	configRoot := getEffectiveConfigRoot(i.configRoot)

	config, err := getEffectiveGlobalConfig(configRoot)
	if err != nil {
		return fatalGrimErrorf("error while reading config: %v", err)
	}

	repos := getAllConfiguredRepos(configRoot)

	var topicARNs []string
	for _, repo := range repos {
		localConfig, err := getEffectiveConfig(configRoot, repo.owner, repo.name)
		if err != nil {
			return fatalGrimErrorf("Error with config for %s/%s. %v", repo.owner, repo.name, err)
		}

		snsTopicARN, err := prepareSNSTopic(config.awsKey, config.awsSecret, config.awsRegion, localConfig.snsTopicName)
		if err != nil {
			return fatalGrimErrorf("error creating SNS Topic %s for %s/%s topic: %v", localConfig.snsTopicName, repo.owner, repo.name, err)
		}

		err = prepareSubscription(config.awsKey, config.awsSecret, config.awsRegion, snsTopicARN, i.queue.ARN)
		if err != nil {
			return fatalGrimErrorf("error subscribing Grim queue %q to SNS topic %q: %v", i.queue.ARN, snsTopicARN, err)
		}

		err = prepareAmazonSNSService(localConfig.gitHubToken, repo.owner, repo.name, snsTopicARN, config.awsKey, config.awsSecret, config.awsRegion)
		if err != nil {
			return fatalGrimErrorf("error creating configuring GitHub AmazonSNS service: %v", err)
		}

		topicARNs = append(topicARNs, snsTopicARN)
	}

	err = setPolicy(config.awsKey, config.awsSecret, config.awsRegion, i.queue.ARN, i.queue.URL, topicARNs)
	if err != nil {
		return fatalGrimErrorf("error setting policy for Grim queue %q with topics %v: %v", i.queue.ARN, topicARNs, err)
	}

	return nil
}

// BuildNextInGrimQueue creates or reuses an SQS queue as a source of work.
func (i *Instance) BuildNextInGrimQueue(logger *log.Logger) error {
	if err := i.checkGrimQueue(); err != nil {
		return err
	}

	configRoot := getEffectiveConfigRoot(i.configRoot)

	globalConfig, err := getEffectiveGlobalConfig(configRoot)
	if err != nil {
		return grimErrorf("error while reading config: %v", err)
	}

	message, err := getNextMessage(globalConfig.awsKey, globalConfig.awsSecret, globalConfig.awsRegion, i.queue.URL)
	if err != nil {
		return grimErrorf("error retrieving message from Grim queue %q: %v", i.queue.URL, err)
	}

	if message != "" {
		hook, err := extractHookEvent(message)
		if err != nil {
			return grimErrorf("error extracting hook from message: %v", err)
		}

		if !(hook.EventName == "push" || hook.EventName == "pull_request" && (hook.Action == "opened" || hook.Action == "reopened" || hook.Action == "synchronize")) {
			return nil
		}

		if hook.EventName == "pull_request" {
			sha, err := pollForMergeCommitSha(globalConfig.gitHubToken, hook.Owner, hook.Repo, hook.PrNumber)
			if err != nil {
				return grimErrorf("error getting merge commit sha: %v", err)
			} else if sha == "" {
				return grimErrorf("error getting merge commit sha: field empty")
			}
			hook.Ref = sha
		}

		localConfig, err := getEffectiveConfig(configRoot, hook.Owner, hook.Repo)
		if err != nil {
			return grimErrorf("error while reading config: %v", err)
		}

		return buildForHook(configRoot, localConfig, *hook, logger)
	}

	return nil
}

// BuildRef builds a git ref immediately.
func (i *Instance) BuildRef(owner, repo, ref string, logger *log.Logger) error {
	configRoot := getEffectiveConfigRoot(i.configRoot)

	config, err := getEffectiveConfig(configRoot, owner, repo)
	if err != nil {
		return fatalGrimErrorf("error while reading config: %v", err)
	}

	return buildForHook(configRoot, config, hookEvent{
		Owner: owner,
		Repo:  repo,
		Ref:   ref,
	}, logger)
}

func buildOnHook(configRoot string, resultPath string, config *effectiveConfig, hook hookEvent) (*executeResult, string, error) {
	return build(config.gitHubToken, configRoot, config.workspaceRoot, resultPath, config.pathToCloneIn, hook.Owner, hook.Repo, hook.Ref, hook.env())
}

//logger coming from method BuildRef is a different logger than that coming from BuildNextInGrimQueue
//unsure if that is a problem
func buildForHook(configRoot string, config *effectiveConfig, hook hookEvent, logger *log.Logger) error {
	return onHook(configRoot, config, hook, buildOnHook, logger)
}

type hookAction func(string, string, *effectiveConfig, hookEvent) (*executeResult, string, error)

func writeHookEvent(resultPath string, hook hookEvent) error {
	hookFile := filepath.Join(resultPath, "hook.json")
	hookBytes, marshalErr := json.Marshal(&hook)
	if marshalErr != nil {
		return marshalErr
	}

	ioutil.WriteFile(hookFile, hookBytes, 0644)
	return nil
}

func onHook(configRoot string, config *effectiveConfig, hook hookEvent, action hookAction, logger *log.Logger) error {
	basename := fmt.Sprintf("%v", time.Now().UnixNano())
	resultPath := makeTree(config.resultRoot, hook.Owner, hook.Repo, basename)

	// TODO: do something with this err too!
	writeHookEvent(resultPath, hook)

	// TODO: do something with the err
	message, _ := notify(config, hook, "", GrimPending)
	if logger != nil {
		logger.Print("grim pending: " + message + "\n")
	}

	result, ws, err := action(configRoot, resultPath, config, hook)
	if err != nil {
		message, _ = notify(config, hook, ws, GrimError)
		if logger != nil {
			logger.Print("grim error: " + message + "\n")
		}
		return fatalGrimErrorf("error during %v: %v", hook.Describe(), err)
	}

	var notifyError error
	if result.ExitCode == 0 {
		message, notifyError = notify(config, hook, ws, GrimSuccess)
		if logger != nil {
			logger.Print("grim success: " + message + "\n")
		}
	} else {
		message, notifyError = notify(config, hook, ws, GrimFailure)
		if logger != nil {
			logger.Print("grim failure: " + message + "\n")
		}
	}

	return notifyError
}

func (i *Instance) checkGrimQueue() error {
	if i.queue == nil {
		return fatalGrimErrorf("the Grim queue must be prepared first")
	}

	return nil
}
