#!/usr/bin/env bash

# -e make it so the entire script fails if a single command fails
# -x means that each command should be printed as it's run
set -e -x

# debug
ls -l $PWD

echo "Create go workspace"
cwd=$(pwd)
mkdir ${cwd}/go
mkdir ${cwd}/go/pkg
mkdir ${cwd}/go/src
mkdir ${cwd}/go/bin

cp -r source-code-hello-world/hello_go_api ${cwd}/go/src

export GOPATH="${cwd}/go"
cd ${cwd}/go/src/hello_go_api
# get libs
go get .

# run a test server
go run main.go &

sleep 5
echo "Go server launched"

# test
go test