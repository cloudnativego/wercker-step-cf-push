#! /bin/bash

pwd
echo "---------------"
ls -la
echo "---------------"

pushd `dirname $0` > /dev/null
SCRIPTPATH=`pwd -P`
popd > /dev/null
eval "$SCRIPTPATH/cf-push-step $@";
