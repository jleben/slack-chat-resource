
all: read-resource post-resource

read-resource:
	docker build -t jakobleben/slack-read-resource -f read/Dockerfile .

post-resource:
	docker build -t jakobleben/slack-post-resource -f post/Dockerfile .
