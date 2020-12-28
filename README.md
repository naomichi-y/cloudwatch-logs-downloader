# cloudwatch-logs-downloader

Download logs for specified period from Amazon CloudWatch Logs and output them to a file in JSON format.

## Setup

```bash
# Set credential
$ cp .env .env.default

$ docker-compose build
```

## Execution sample

```bash
$ docker-compose up -d
$ docker-compose exec app go run main.go -group=ecs/production-log -start="2020-12-27 15:59:00" -end="2020-12-27 15:59:59"

2020/12/28 06:44:13 Generated log file: ./dist/result_2020122864407.log
```

* group: Log group name
* start: Filter start date and time (UTC)
* end: Filter end date and time (UTC)
