#! /bin/bash
image=$1
request=$2

cat "$request" | docker run --rm -i \
-e BUILD_NAME=mybuild \
-e BUILD_JOB_NAME=myjob \
-e BUILD_PIPELINE_NAME=mypipe \
-e BUILD_TEAM_NAME=myteam \
-e ATC_EXTERNAL_URL="https://example.com" \
-v "$(pwd)/out:/tmp/resource/out" "$image" /opt/resource/in /tmp/resource/out
