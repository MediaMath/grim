package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/service/sns"
)

func prepareSNSTopic(key, secret, region, topic string) (string, error) {
	config := getConfig(key, secret, region)

	topicARN, err := findExistingTopicARN(config, topic)
	if err != nil {
		topicARN, err = createTopic(config, topic)
	}

	if err != nil {
		return "", err
	}

	return topicARN, nil
}

func prepareSubscription(key, secret, region, topicARN, queueARN string) error {
	config := getConfig(key, secret, region)

	subARN, err := findSubscription(config, topicARN, queueARN)
	if err != nil {
		return err
	}

	if subARN == "" {
		subARN, err = createSubscription(config, topicARN, queueARN)
	}

	if subARN == "" {
		return fmt.Errorf("failed to create subscription")
	}

	return err
}

func createSubscription(config *aws.Config, topicARN, queueARN string) (string, error) {
	svc := sns.New(config)

	params := &sns.SubscribeInput{
		Protocol: aws.String("sqs"),
		TopicARN: aws.String(topicARN),
		Endpoint: aws.String(queueARN),
	}

	resp, err := svc.Subscribe(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while creating subscription to SNS topic: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return "", fmt.Errorf("error while creating subscription to SNS topic: %v", err)
	} else if resp == nil || resp.SubscriptionARN == nil {
		return "", fmt.Errorf("error while creating subscription to SNS topic")
	}

	return *resp.SubscriptionARN, nil
}

func findSubscription(config *aws.Config, topicARN, queueARN string) (string, error) {
	svc := sns.New(config)

	params := &sns.ListSubscriptionsByTopicInput{
		TopicARN: aws.String(topicARN),
	}

	for {
		resp, err := svc.ListSubscriptionsByTopic(params)
		if awserr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("aws error while listing subscriptions to SNS topic: %v %v", awserr.Code, awserr.Message)
		} else if err != nil {
			return "", fmt.Errorf("error while listing subscriptions to SNS topic: %v", err)
		} else if resp == nil || resp.Subscriptions == nil {
			break
		}

		for _, sub := range resp.Subscriptions {
			if sub.Endpoint != nil && *sub.Endpoint == queueARN && sub.Protocol != nil && *sub.Protocol == "sqs" && sub.SubscriptionARN != nil {
				return *sub.SubscriptionARN, nil
			}
		}

		if resp.NextToken != nil {
			params.NextToken = resp.NextToken
		} else {
			break
		}
	}

	return "", nil
}

func createTopic(config *aws.Config, topic string) (string, error) {
	svc := sns.New(config)

	params := &sns.CreateTopicInput{
		Name: aws.String(topic),
	}

	resp, err := svc.CreateTopic(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while creating SNS topic: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return "", fmt.Errorf("error while creating SNS topic: %v", err)
	} else if resp == nil || resp.TopicARN == nil {
		return "", nil
	}

	return *resp.TopicARN, nil
}

func findExistingTopicARN(config *aws.Config, topic string) (string, error) {
	svc := sns.New(nil)

	params := &sns.ListTopicsInput{
		NextToken: nil,
	}

	for {
		resp, err := svc.ListTopics(params)
		if awserr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("aws error while listing SNS topics: %v %v", awserr.Code, awserr.Message)
		} else if err != nil {
			return "", fmt.Errorf("error while listing SNS topics: %v", err)
		} else if resp == nil || resp.Topics == nil {
			break
		}

		for _, topicPtr := range resp.Topics {
			if topicPtr != nil && topicPtr.TopicARN != nil && strings.HasSuffix(*topicPtr.TopicARN, topic) {
				return *topicPtr.TopicARN, nil
			}
		}

		if resp.NextToken != nil {
			params.NextToken = resp.NextToken
		} else {
			break
		}
	}

	return "", nil
}
