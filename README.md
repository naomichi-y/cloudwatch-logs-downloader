# cloudwatch-logs-downloader

[![Go Report Card](https://goreportcard.com/badge/github.com/naomichi-y/cloudwatch-logs-downloader)](https://goreportcard.com/report/github.com/naomichi-y/cloudwatch-logs-downloader)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Download logs for specified period from Amazon CloudWatch Logs and output them to a file in JSON format.

## Setup

```bash
# Set AWS credential
$ cp .env .env.default

$ docker build -t cld .
```

## Usage

```bash
$ docker run --rm -it --env-file=.env -v ${PWD}:/go/src/app cld -group={group}
```

|Argumnet|Required|Description|Default|
|---|---|---|---|
|`-group`|Yes|Log group name||
|`-prefix`||Prefix name when searching log groups||
|`-start`||Log stream event search start date and time (UTC)|Start time will be 10 minutes before current time|
|`-end`||Log stream event search end date and time (UTC)|Current time|
|`-pattern`||Regular expressions for events to be extracted||

## Execution sample

```bash
$ docker run --rm -it --env-file=.env -v ${PWD}:/go/src/app cld -group=ecs/production-log -start="2020-12-27 15:59:00" -end="2020-12-27 15:59:59"

2020/12/29 07:05:42 Write results...
2020/12/29 07:05:42 Generated log file: ./dist/result_2020122970533.log

$ cat ./dist/result_2020122970533.log
[
  {
    "IngestionTime": "2020-12-27 15:59:22 +0000 UTC",
    "LogStream": "app/app/5302adeb527b42a0acbb11ac3444d98f",
    "Message": "Foo",
    "Timestamp": "2020-12-27 15:59:21 +0000 UTC"
  },
  {
    "IngestionTime": "2020-12-27 15:59:22 +0000 UTC",
    "LogStream": "app/app/5302adeb527b42a0acbb11ac3444d98f",
    "Message": "Bar",
    "Timestamp": "2020-12-27 15:59:21 +0000 UTC"
  },
  {
    "IngestionTime": "2020-12-27 15:59:22 +0000 UTC",
    "LogStream": "app/app/5302adeb527b42a0acbb11ac3444d98f",
    "Message": "Baz",
    "Timestamp": "2020-12-27 15:59:21 +0000 UTC"
  },
  ...
]
```
