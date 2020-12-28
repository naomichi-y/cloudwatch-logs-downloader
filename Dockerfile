FROM golang:1.15-alpine

RUN apk update && apk add --no-cache git
RUN go get -u github.com/aws/aws-sdk-go

WORKDIR /go/src/app
COPY . .
