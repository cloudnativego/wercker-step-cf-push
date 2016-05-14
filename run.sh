#! /bin/bash

pwd
echo "---------------"
ls -la
echo "---------------"
eval "./cf-push-step $@";
