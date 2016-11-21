package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func prepareSNSTopic(key, secret, region, topic string) (string, error) {
	session := getSession(key, secret, region)

	topicARN, err := findExistingTopicARN(session, topic)
	if err != nil || topicARN == "" {
		topicARN, err = createTopic(session, topic)
	}

	if err != nil {
		return "", err
	}

	return topicARN, nil
}

func prepareSubscription(key, secret, region, topicARN, queueARN string) error {
	session := getSession(key, secret, region)

	subARN, err := findSubscription(session, topicARN, queueARN)
	if err != nil {
		return err
	}

	if subARN == "" {
		subARN, err = createSubscription(session, topicARN, queueARN)
	}

	if subARN == "" {
		return fmt.Errorf("failed to create subscription")
	}

	return err
}

func createSubscription(session *session.Session, topicARN, queueARN string) (string, error) {
	svc := sns.New(session)

	params := &sns.SubscribeInput{
		Protocol: aws.String("sqs"),
		TopicArn: aws.String(topicARN),
		Endpoint: aws.String(queueARN),
	}

	resp, err := svc.Subscribe(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while creating subscription to SNS topic: %v %v", awserr.Code(), awserr.Message())
	} else if err != nil {
		return "", fmt.Errorf("error while creating subscription to SNS topic: %v", err)
	} else if resp == nil || resp.SubscriptionArn == nil {
		return "", fmt.Errorf("error while creating subscription to SNS topic")
	}

	return *resp.SubscriptionArn, nil
}

func findSubscription(session *session.Session, topicARN, queueARN string) (string, error) {
	svc := sns.New(session)

	params := &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicARN),
	}

	for {
		resp, err := svc.ListSubscriptionsByTopic(params)
		if awserr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("aws error while listing subscriptions to SNS topic: %v %v", awserr.Code(), awserr.Message())
		} else if err != nil {
			return "", fmt.Errorf("error while listing subscriptions to SNS topic: %v", err)
		} else if resp == nil || resp.Subscriptions == nil {
			break
		}

		for _, sub := range resp.Subscriptions {
			if sub.Endpoint != nil && *sub.Endpoint == queueARN && sub.Protocol != nil && *sub.Protocol == "sqs" && sub.SubscriptionArn != nil {
				return *sub.SubscriptionArn, nil
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

func createTopic(session *session.Session, topic string) (string, error) {
	svc := sns.New(session)

	params := &sns.CreateTopicInput{
		Name: aws.String(topic),
	}

	resp, err := svc.CreateTopic(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while creating SNS topic: %v %v", awserr.Code(), awserr.Message())
	} else if err != nil {
		return "", fmt.Errorf("error while creating SNS topic: %v", err)
	} else if resp == nil || resp.TopicArn == nil {
		return "", nil
	}

	return *resp.TopicArn, nil
}

func findExistingTopicARN(session *session.Session, topic string) (string, error) {
	svc := sns.New(session)

	params := &sns.ListTopicsInput{
		NextToken: nil,
	}

	for {
		resp, err := svc.ListTopics(params)
		if awserr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("aws error while listing SNS topics: %v %v", awserr.Code(), awserr.Message())
		} else if err != nil {
			return "", fmt.Errorf("error while listing SNS topics: %v", err)
		} else if resp == nil || resp.Topics == nil {
			break
		}

		for _, topicPtr := range resp.Topics {
			if topicPtr != nil && topicPtr.TopicArn != nil && strings.HasSuffix(*topicPtr.TopicArn, topic) {
				return *topicPtr.TopicArn, nil
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
