#!/bin/bash -e

#FIXME: should resolve absolute path to script's parent dir insead
declare -x BASEPATH=$PWD

if [ -d $BASEPATH/ci_scripts/staging ]; then
  rm -rf $BASEPATH/ci_scripts/staging
fi

mkdir -p $BASEPATH/ci_scripts/staging

declare -x GOPATH=${HOME}/go:$BASEPATH
cd $BASEPATH/src/config
go test

cd $BASEPATH/src/watcher
declare -x GOOS=linux
go build

mv watcher $BASEPATH/ci_scripts/staging
cp $BASEPATH/src/config/sampleconfig.yaml $BASEPATH/ci_scripts/staging

cd $BASEPATH/ci_scripts
docker build .
