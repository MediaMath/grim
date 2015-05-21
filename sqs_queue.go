package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"fmt"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/service/sqs"
)

var policyFormat = `{
	"Version": "2008-10-17",
	"Id": "grim-policy",
	"Statement": [
	  {
			"Sid": "1",
			"Effect": "Allow",
			"Principal": {
				"AWS": "*"
			},
			"Action": "SQS:*",
			"Resource": "%v",
			"Condition": {
				"ArnEquals": {
					"aws:SourceArn": %v
				}
			}
		}
	]
}`

type sqsQueue struct {
	URL string
	ARN string
}

func prepareSQSQueue(key, secret, region, queue string) (*sqsQueue, error) {
	config := getConfig(key, secret, region)

	queueURL, err := getQueueURLByName(config, queue)
	if err != nil {
		queueURL, err = createQueue(config, queue)
	}

	if err != nil {
		return nil, err
	}

	queueARN, err := getARNForQueueURL(config, queueURL)
	if err != nil {
		return nil, err
	}

	return &sqsQueue{queueURL, queueARN}, nil
}

func getNextMessage(key, secret, region, queueURL string) (string, error) {
	config := getConfig(key, secret, region)

	message, err := getMessage(config, queueURL)
	if err != nil {
		return "", err
	} else if message == nil || message.ReceiptHandle == nil {
		return "", nil
	}

	err = deleteMessage(config, queueURL, *message.ReceiptHandle)
	if err != nil {
		return "", err
	}

	if message.Body == nil {
		return "", nil
	}

	return *message.Body, nil
}

func setPolicy(key, secret, region, queueARN, queueURL string, topicARNs []string) error {
	svc := sqs.New(getConfig(key, secret, region))

	bs, err := json.Marshal(topicARNs)
	if err != nil {
		return fmt.Errorf("error while creating policy for SQS queue: %v", err)
	}

	policy := fmt.Sprintf(policyFormat, queueARN, string(bs))

	params := &sqs.SetQueueAttributesInput{
		Attributes: &map[string]*string{
			"Policy": aws.String(policy),
		},
		QueueURL: aws.String(queueURL),
	}

	_, err = svc.SetQueueAttributes(params)
	if awserr, ok := err.(awserr.Error); ok {
		return fmt.Errorf("aws error while setting policy for SQS queue: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return fmt.Errorf("error while setting policy for SQS queue: %v", err)
	}

	return nil
}

func getQueueURLByName(config *aws.Config, queue string) (string, error) {
	svc := sqs.New(config)

	params := &sqs.GetQueueURLInput{
		QueueName: aws.String(queue),
	}

	resp, err := svc.GetQueueURL(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while getting URL for SQS queue: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return "", fmt.Errorf("error while getting URL for SQS queue: %v", err)
	} else if resp == nil || resp.QueueURL == nil {
		return "", nil
	}

	return *resp.QueueURL, nil
}

func getARNForQueueURL(config *aws.Config, queueURL string) (string, error) {
	svc := sqs.New(config)

	arnKey := "QueueArn"

	params := &sqs.GetQueueAttributesInput{
		QueueURL: aws.String(string(queueURL)),
		AttributeNames: []*string{
			aws.String(arnKey),
		},
	}

	resp, err := svc.GetQueueAttributes(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while getting ARN for SQS queue: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return "", fmt.Errorf("error while getting ARN for SQS queue: %v", err)
	} else if resp == nil || resp.Attributes == nil {
		return "", nil
	}

	atts := *resp.Attributes

	arnPtr, ok := atts[arnKey]
	if !ok || arnPtr == nil {
		return "", nil
	}

	return *arnPtr, nil
}

func createQueue(config *aws.Config, queue string) (string, error) {
	svc := sqs.New(config)

	params := &sqs.CreateQueueInput{
		QueueName: aws.String(queue),
		Attributes: &map[string]*string{
			"ReceiveMessageWaitTimeSeconds": aws.String("5"),
		},
	}

	resp, err := svc.CreateQueue(params)
	if awserr, ok := err.(awserr.Error); ok {
		return "", fmt.Errorf("aws error while creating SQS queue: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return "", fmt.Errorf("error while creating SQS queue: %v", err)
	} else if resp == nil || resp.QueueURL == nil {
		return "", nil
	}

	return *resp.QueueURL, nil
}

func getMessage(config *aws.Config, queueURL string) (*sqs.Message, error) {
	svc := sqs.New(config)

	params := &sqs.ReceiveMessageInput{
		QueueURL:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Long(1),
	}

	resp, err := svc.ReceiveMessage(params)
	if awserr, ok := err.(awserr.Error); ok {
		return nil, fmt.Errorf("aws error while receiving message from SQS: %v %v", awserr.Code, awserr.Message)
	} else if err != nil {
		return nil, fmt.Errorf("error while receiving message from SQS: %v", err)
	} else if resp == nil || len(resp.Messages) == 0 {
		return nil, nil
	}

	return resp.Messages[0], nil
}

func deleteMessage(config *aws.Config, queueURL string, receiptHandle string) error {
	svc := sqs.New(config)

	params := &sqs.DeleteMessageInput{
		QueueURL:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := svc.DeleteMessage(params)
	if awserr, ok := err.(awserr.Error); ok {
		return fmt.Errorf("aws error while deleting message from SQS: %v %v", awserr.Code(), awserr.Message())
	} else if err != nil {
		return fmt.Errorf("error while deleting message from SQS: %v", err)
	}

	return nil
}
