#! /bin/bash
image=$1
request=$2
cat "$request" | docker run --rm -i -v "$(pwd)/out:/tmp/resource/out" "$image" /opt/resource/in /tmp/resource/out
