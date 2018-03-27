# Concourse CI test

- [Concourse CI test](#concourse-ci-test)
  - [Deploy Concourse CI with local Docker](#deploy-concourse-ci-with-local-docker)
    - [Concourse config](#concourse-config)
    - [Run Docker compose](#run-docker-compose)
  - [Install and configure fly CLI](#install-and-configure-fly-cli)
  - [Fly CLI bash completion](#fly-cli-bash-completion)
  - [Deploy a pipeline](#deploy-a-pipeline)
    - [Hello world pipeline](#hello-world-pipeline)
    - [Go API pipeline](#go-api-pipeline)
  - [Test locally](#test-locally)

## Deploy Concourse CI with local Docker

The compose file will deploy:
- a Docker registry server
- database server for concourse
- concourse server

### Concourse config
Before running the script, export your IP
```
export CONCOURSE_EXTERNAL_URL="10.33.101.57"
```

Generate the needed keys as per the Concourse [install guide](https://concourse.ci/docker-repository.html) in a directory called "keys".
```
mkdir -p keys/web keys/worker

ssh-keygen -t rsa -f ./keys/web/tsa_host_key -N ''
ssh-keygen -t rsa -f ./keys/web/session_signing_key -N ''

ssh-keygen -t rsa -f ./keys/worker/worker_key -N ''

cp ./keys/worker/worker_key.pub ./keys/web/authorized_worker_keys
cp ./keys/web/tsa_host_key.pub ./keys/worker
```

### Run Docker compose

Run the compose script
```
docker-compose up
```

The GUI should now be available on the address you've set in `CONCOURSE_EXTERNAL_URL` via your browser.
Login: concourse
Password: changeme

## Install and configure fly CLI

Concourse work with a heavy client called `fly`. Download the compiled binary from [the official website](https://concourse.ci/downloads.html).
You may need to Unzip it. Place it in a directory on you PATH:
```
cp path/to/fly /usr/local/bin/fly
```

Connect to your backend
```
fly -t lite login -c <CONCOURSE_EXTERNAL_URL>
```

E.g
```
fly -t lite login -c http://10.39.9.94:8080
```

This will save the target backend and attach it to the id "lite". You can change this name as well.

## Fly CLI bash completion

Add the following content into `/etc/bash_completion.d/fly`
```
#compdef fly
autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit
_fly_bash_autocomplete() {
    args=("${COMP_WORDS[@]:1:$COMP_CWORD}")
    # Only split on newlines
    local IFS=$'\n'
    # Call completion (note that the first element of COMP_WORDS is
    # the executable itself)
    COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
    return 0
}
complete -F _fly_bash_autocomplete fly
```

If zhh is used instead of bash, update the zshrc with following commands:
```
autoload bashcompinit
bashcompinit
source /etc/bash_completion.d/fly
```

## Deploy a pipeline

### Hello world pipeline

This part can be retrieved on t[he official website](https://concourse-ci.org/hello-world.html).

Set the pipeline
```
fly -t lite set-pipeline -p hello-world -c hello.yml
```

Unpause
```
fly -t lite unpause-pipeline --pipeline hello-world
```

### Go API pipeline

This pipeline will:
- detect change on a git repo
- test the project
- build the project
- push into a Swarm cluster the API

Set the pipeline
```
fly -t lite set-pipeline -p go-hello -c go_hello_pipeline.yml
```

Unpause
```
fly -t lite unpause-pipeline --pipeline go-hello
```


## Test locally

Concourse allows you to test the code of the CI locally before pushing the pipeline to the server.

Local execution of the testing job
```
fly -t lite execute --config ci/01-task-tests.yml --input source-code-hello-world=.
```

