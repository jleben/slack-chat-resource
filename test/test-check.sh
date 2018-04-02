#! /bin/bash
type=$1
request=$2

if [[ -z $type || -z $request ]]; then
    echo "Required arguments: <resource type> <request file>"
    exit 1
fi

cat "$request" | docker run --rm -i jakobleben/slack-$type-resource /opt/resource/check
