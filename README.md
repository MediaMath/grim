#Grim: "No It Doesn't Do That"

Grim is the "GitHub Responder In MediaMath".  We liked the acronym and awkwardly filled in the details to fit it.  In short, it is a task runner that is triggered by GitHub push/pull request hooks that is intended as a much simpler and easy-to-use build server than the more modular alternatives (eg. Jenkins).

On start up, Grim will:

1. Create (or reuse) an Amazon SQS queue with the name specified in its config file as `GrimQueueName`
2. Detect which GitHub repositories it is configured to work with
3. Create (or reuse) an Amazon SNS topic for each repository
4. Configure each created topic to push to the Grim queue
5. Configure each repositories' AmazonSNS service to push hook updates (`push`, `pull_request`) to the topic

![Grimd data flow](docs/grimd.png "An example including 3 Grimd's (one in EC2 and two MacBookPros) and two repositories.")

_An example including two repositories watched by three Grimd's (one in EC2 and two MacBookPros)._

Each GitHub repo can push to exactly one SNS topic.  Multiple SQS queues can subscribe to one topic and multiple Grim instances can read from the same SQS queue.  If a Grim instance isn't configured to respond to the repo specified in the hook it silently ignores the event.

##Installation

### 1. Get grimd

```bash
wget https://artifactory.mediamath.com/artifactory/libs-release-global/com/mediamath/grim/grimd/[RELEASE]/grimd-[RELEASE].zip
```

### 2. Give grimd user the ability to clone from the github repos you are interested in.

The user that is running your installed grimd process needs to be able to clone the github repos it is responding to.  To that end it will need to be able to run git and should have access to the github repos (either via private key, or because they are all public).

### 3. Global Configuration

Grim tries to honor the conventional [Unix filesystem hierarchy](http://en.wikipedia.org/wiki/Unix_filesystem#Conventional_directory_layout) as much as possible.  Configuration files are by default found in `/etc/grim`.  You may override that default by specifying `--config-root [some dir]`, more briefly `-c [some dir]` or by setting the `GRIM_CONFIG_ROOT` environment variable.  Inside that directory there is expected to be a `config.json` that specifies the other paths used as well as global defaults.  Here is an example:

```
{
	"GrimQueueName": "grim-queue",
	"ResultRoot": "/var/log/grim",
	"WorkspaceRoot": "/var/tmp/grim",
	"AWSRegion": "us-east-1",
	"AWSKey": "xxxx",
	"AWSSecret": "xxxx",
	"GitHubToken": "xxxx",
	"HipChatToken": "xxxx"
}
```

If you don't configure `GrimQueueName`, `ResultRoot` or `WorkspaceRoot` Grim will use default values.  The AWS credentials supplied must be able to create and modify SNS topics and SQS queues.

#### Required GitHub token scopes

* `write:repo_hook` to be able to create/edit repository hooks
* `repo:status` to be able set commit statuses

#### Prepare.sh

There is also expected to be an executable script in the configuration root called `prepare.sh` that will take clones a repository in a path and checks out a ref:

```bash
#!/bin/bash

env | sort > ~/grimd.log

set -eu

FULL_CLONE_PATH=$(pwd)"/$CLONE_PATH"
mkdir -p "$FULL_CLONE_PATH"

git clone "ssh://git@github.com/$GH_OWNER/$GH_REPO.git" "$FULL_CLONE_PATH"
cd "$FULL_CLONE_PATH"

if [ "$GH_EVENT_NAME" == "pull_request" ]; then
	git fetch origin "refs/pull/$GH_PR_NUMBER/head:pull_branch"
	git checkout pull_branch
else
	git checkout "$GH_REF"
fi
```

It will be run with its working directory being the workspace that was created for this build.  You can also override the global prepare script by placing a `prepare.sh` in the configuration directory for a specific repo.

The environment variables available to this script are documented [here](#environment-variables).

### 4. Repository Configuration

In order for Grim to respond to GitHub events it needs subdirectories to be made in the configuration root.  Inside those subdirectories should be a `config.json` and optionally a `build.sh`.  Here is an example directory structure:

```
/etc/grim
/etc/grim/config.json
/etc/grim/MediaMath
/etc/grim/MediaMath/grim
/etc/grim/MediaMath/grim/config.json
/etc/grim/MediaMath/grim/build.sh
```

The file `config.json` can have an empty JSON object or have the following fields:

```
{
	"GitHubToken": "xxxx",
	"HipChatToken": "xxxx"
	"HipChatRoom": "xxxx",
	"PathToCloneIn": "/go/src/github.com/MediaMath/grim"
}
```

The GitHub and HipChat tokens will override the global ones if present.  The HipChat room is optional and if present will indicate that status messages will go to that room.  The field `PathToCloneIn` is relative to the workspace that was created for this build.

#### Build script location

Grim will look for a build script first in the configuration directory for the repo as `build.sh` and failing that in the root of the cloned repo as either `.grim_build.sh` or `grim_build.sh`.

The environment variables available to this script are documented [here](#environment-variables).

### Environment Variables 
```
CLONE_PATH= the path relative to the workspace to clone the repo in
GH_EVENT_NAME= either 'push', 'pull_request' or '' (for manual builds)
GH_ACTION= the sub action of a pull request (eg. 'opened', 'closed', or 'reopened', 'synchronize') or blank for other event types
GH_USER_NAME= the user initiating the event
GH_OWNER= the owner part of a repo (eg. 'MediaMath')
GH_REPO= the name of a repo (eg. 'grim')
GH_TARGET= the branch that a commit was merged to
GH_REF= the ref to build
GH_STATUS_REF= the ref to set the status of
GH_URL= the GitHub URL to find the changes at
```
