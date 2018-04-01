
FROM golang:alpine as builder
RUN apk update
RUN apk add git
RUN go get github.com/nlopes/slack
COPY ./protocol /go/src/github.com/jleben/trigger-resource/protocol/
COPY ./in /go/src/github.com/jleben/trigger-resource/in/
COPY ./check /go/src/github.com/jleben/trigger-resource/check/
RUN go build -o /assets/in github.com/jleben/trigger-resource/in
RUN go build -o /assets/check github.com/jleben/trigger-resource/check

FROM alpine as resource
RUN apk update
RUN apk add ca-certificates
COPY --from=builder /assets /opt/resource

FROM resource
