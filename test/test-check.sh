#! /bin/bash
request=$1
cat "$request" | docker run --rm -i jakobleben/slack-request-resource /opt/resource/check
