# Concourse CI test

- [Concourse CI test](#concourse-ci-test)
  - [Enable remote API for Dockerd](#enable-remote-api-for-dockerd)
  - [Deploy Concourse CI with local Docker](#deploy-concourse-ci-with-local-docker)
    - [Concourse config](#concourse-config)
    - [Run Docker compose](#run-docker-compose)
  - [Install and configure fly CLI](#install-and-configure-fly-cli)
  - [Fly CLI bash completion](#fly-cli-bash-completion)
  - [Push required Docker images](#push-required-docker-images)
  - [Deploy a pipeline](#deploy-a-pipeline)
    - [Hello world pipeline](#hello-world-pipeline)
    - [Go API pipeline](#go-api-pipeline)
  - [Test locally](#test-locally)

## Enable remote API for Dockerd

For this demo, we will need to discuss with the local docker daemon in order to deploy the built project.

Create /etc/systemd/system/docker.service.d/startup_options.conf
```
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H tcp://0.0.0.0:2376
```
Then reload the docker daemon
```
sudo systemctl daemon-reload
sudo systemctl restart docker.service
```

Check
```
curl -X GET http://127.0.0.1:2376/v1.24/containers/json
```

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

## Push required Docker images

In the pipeline, one of our task is based on a own made docker image. We need to make this image available from our private registry.

First step, build the image locally
```
cd ci/deployment
docker build --force-rm=true -f hello.deployment.dockerfile -t hello-deployment .
```

Tag and push into the private registry that has been deployed with docker compose previously
```
docker tag hello-deployment 10.33.101.57:5000/hello-deployment
docker push 10.33.101.57:5000/hello-deployment
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

Local execution of the deployment
```
fly -t lite execute --config ci/02-task-deploy.yml --input source-code-hello-world=.
```
