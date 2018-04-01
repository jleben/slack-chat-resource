#! /bin/bash
image=$1
request=$2
cat "$request" | docker run --rm -i "$image" /opt/resource/check
