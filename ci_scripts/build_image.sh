#!/bin/bash -e

#FIXME: should resolve absolute path to script's parent dir insead
declare -x BASEPATH=$PWD
declare -x GO111MODULE=on

if [ -d $BASEPATH/ci_scripts/staging ]; then
  rm -rf $BASEPATH/ci_scripts/staging
fi

mkdir -p $BASEPATH/ci_scripts/staging

go test ./...

declare -x GOOS=linux
cd $BASEPATH/watcher
go build

mv watcher $BASEPATH/ci_scripts/staging
cp $BASEPATH/config/sampleconfig.yaml $BASEPATH/ci_scripts/staging

cd $BASEPATH/ci_scripts
docker build .
