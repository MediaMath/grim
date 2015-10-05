#!/bin/bash

set -eu

. /opt/golang/go1.4.2/bin/go_env.sh

export GOPATH="$(pwd)/go"
export PATH="$GOPATH/bin:$PATH"

cd "./$CLONE_PATH"

go get ./...
go get github.com/golang/lint/golint

if [ "$GH_EVENT_NAME" == "push" -a "$GH_TARGET" == "master" ]; then
	#on merge of master publish to release artifactory repo
	REPOSITORY=libs-release-global make clean check publish packer
elif [ "$GH_EVENT_NAME" == "pull_request" -a "$GH_TARGET" == "master" ]; then
	#on pull requests publish to staging repo, allows for end to end testing with automation
	REPOSITORY=libs-staging-global make clean check publish
else 
	#otherwise just build it
	make clean check grimd
fi
