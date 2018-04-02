
FROM golang:alpine as build-env
RUN apk update
RUN apk add git
RUN go get github.com/nlopes/slack
COPY ./utils /go/src/github.com/jleben/slack-chat-resource/utils
RUN go build -o /assets/utils github.com/jleben/slack-chat-resource/utils

FROM build-env as build-post
#COPY ./post/check /go/src/github.com/jleben/slack-chat-resource/post/check/
#COPY ./post/in /go/src/github.com/jleben/slack-chat-resource/post/in/
COPY ./post/out /go/src/github.com/jleben/slack-chat-resource/post/out/
#RUN go build -o /assets/check github.com/jleben/slack-chat-resource/post/check
#RUN go build -o /assets/in github.com/jleben/slack-chat-resource/post/in
RUN go build -o /assets/out github.com/jleben/slack-chat-resource/post/out

FROM alpine
RUN apk update
RUN apk add ca-certificates
COPY --from=build-post /assets /opt/resource
