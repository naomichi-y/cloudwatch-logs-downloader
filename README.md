# cloudwatch-logs-downloader

Download logs for specified period from Amazon CloudWatch Logs and output them to a file in JSON format.

## Setup

```bash
# Set AWS credential
$ cp .env .env.default

$ docker build -t cld .
```

## Usage

```bash
$ docker run --rm -it --env-file=.env -v ${PWD}:/go/src/app cld go run main.go -group={group}
```

|Argumnet|Required|Description|Default|
|---|---|---|---|
|`-group`|Yes|Log group name||
|`-prefix`||Prefix name when searching log groups||
|`-start`||Filter start date and time (UTC)|Start time will be 10 minutes before current time|
|`-end`||Filter end date and time (UTC)|Current time|

## Execution sample

```bash
$ docker run --rm -it --env-file=.env -v ${PWD}:/go/src/app cld go run main.go -group=ecs/production-log -start="2020-12-27 15:59:00" -end="2020-12-27 15:59:59"

2020/12/28 06:44:13 Generated log file: ./dist/result_2020122864407.log
```
