
FROM golang:alpine as builder
RUN apk update
RUN apk add git
RUN go get github.com/nlopes/slack
# FIXME: go package location
COPY ./protocol /go/src/github.com/jleben/trigger-resource/protocol/
COPY ./check /go/src/github.com/jleben/trigger-resource/check/
COPY ./in /go/src/github.com/jleben/trigger-resource/in/
COPY ./out /go/src/github.com/jleben/trigger-resource/out/
RUN go build -o /assets/check github.com/jleben/trigger-resource/check
RUN go build -o /assets/in github.com/jleben/trigger-resource/in
RUN go build -o /assets/out github.com/jleben/trigger-resource/out

FROM alpine as resource
RUN apk update
RUN apk add ca-certificates
COPY --from=builder /assets /opt/resource

FROM resource
