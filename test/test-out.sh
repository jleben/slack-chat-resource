#! /bin/bash
request=$1

cat "$request" | docker run --rm -i \
-e BUILD_NAME=mybuild \
-e BUILD_JOB_NAME=myjob \
-e BUILD_PIPELINE_NAME=mypipe \
-e BUILD_TEAM_NAME=myteam \
-e ATC_EXTERNAL_URL="https://example.com" \
-v "$(pwd)/src:/tmp/resource/src" jakobleben/slack-request-resource /opt/resource/out /tmp/resource/src
