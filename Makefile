.PHONY:	grimd publish test check clean run cover part ansible

# Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

TIMESTAMP := $(shell date +"%s")
BUILD_TIME := $(shell date +"%Y%m%d.%H%M%S")
ARTIFACTORY_HOST = artifactory.mediamath.com
SHELL := /bin/bash

VERSION = $(strip $(TIMESTAMP))
ifndef REPOSITORY 
	REPOSITORY = libs-staging-global
endif

LDFLAGS = -ldflags "-X main.version $(VERSION)-$(BUILD_TIME)"

ifdef VERBOSE
	TEST_VERBOSITY=-v
else
	TEST_VERBOSITY=
endif

grimd: 
	go get ./...
	go install $(LDFLAGS) github.com/MediaMath/grim/grimd

tmp/grimd-$(VERSION).zip: grimd | tmp 
	export PATH=$$PATH:$${GOPATH//://bin:}/bin; zip -r -j $@ $$(which grimd)

test:
	go test $(TEST_VERBOSITY) ./...

part: 
	go get -a github.com/MediaMath/part

publish: part tmp/grimd-$(VERSION).zip
	part -verbose -credentials=$(HOME)/.ivy2/credentials/$(ARTIFACTORY_HOST) -h="https://$(ARTIFACTORY_HOST)/artifactory" -r=$(REPOSITORY) -g=com.mediamath.grim -a=grimd -v=$(VERSION) tmp/grimd-$(VERSION).zip

cover: tmp
	cvr -o=tmp/coverage -short ./...

clean:
	go clean ./...
	rm -rf grimd/grim tmp/*

tmp:
	mkdir tmp

check: test
	go vet ./...
	golint ./...

ansible:
	cd ansible && ansible-playbook -i inventory site.xml
