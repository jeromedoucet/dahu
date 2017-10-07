#!/usr/bin/env bash

# script embeded in the test image to perform tests on
# the Dahu run module.

if [ -z "$REPO_URL" ]; then
    echo "Need to set REPO_URL"
    exit 1
fi

if [ "$STATUS" == "failure" ]
then
    echo "Failure"
    exit 1
fi

if [ "$STATUS" == "timeout" ]
then
    sleep 10
fi

echo "Success"
exit 0
