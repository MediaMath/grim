package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getSession(key, secret, region string) *session.Session {
	creds := credentials.NewStaticCredentials(key, secret, "")
	return session.New(&aws.Config{Credentials: creds, Region: &region})
}

func getAccountIDFromARN(arn string) string {
	ps := strings.Split(arn, ":")
	if len(ps) > 5 {
		return ps[4]
	}

	return ""
}
