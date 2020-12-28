package main

import (
  "log"
  "time"
  "os"
  "fmt"
  "encoding/json"
  "flag"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var file string

func write(data string) {
  log.Print("Write results...")

  fp, err := os.OpenFile(file, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)

  if err != nil {
    log.Fatal(err)
  }

  defer fp.Close()
  fmt.Fprintln(fp, data)
}

func log_events(svc *cloudwatchlogs.CloudWatchLogs, g string, s string, b int64, f int64) {
  log.Print("Searching log events...")

  var nextToken *string = nil
  var list []map[string]string

  for {
    resp, err := svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
      LogGroupName: aws.String(g),
      LogStreamName: aws.String(s),
      StartTime: aws.Int64(b),
      EndTime: aws.Int64(f),
      NextToken: nextToken,
    })

    if err != nil {
      log.Fatal(err)
    }

    for _, v := range resp.Events {
      r := map[string]string{
        "Message": aws.StringValue(v.Message),
        "Timestamp": time.Unix(aws.Int64Value(v.Timestamp) / 1000, 0).String(),
      }
      list = append(list, r)
      log.Print(r)
    }

    if aws.StringValue(nextToken) == aws.StringValue(resp.NextForwardToken) {
      break
    }

    nextToken = resp.NextForwardToken
  }

  data, err := json.MarshalIndent(list, "", "  ")

  if err != nil {
    log.Fatal(err)
  }

  write(string(data))
}

func main() {
  file = "./dist/result_" + time.Now().Format("2006010230405") + ".log"

  var group = flag.String("group", "", "Log group")
  var start = flag.String("start", "", "Start date")
  var end = flag.String("end", "", "End date")

  flag.Parse()

  sess := session.Must(session.NewSession())
  svc := cloudwatchlogs.New(
    sess,
    aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")),
  )

  var nextToken *string = nil
  var layout = "2006-01-02 15:04:05"

  s, err := time.Parse(layout, *start)

  if err != nil {
    log.Fatal(err)
  }

  e, err := time.Parse(layout, *end)

  if err != nil {
    log.Fatal(err)
  }

  startUnixtime := int64(s.Unix() * 1000)
  endUnixtime := int64(e.Unix() * 1000)

  log.Print("Searching log groups...")

  for {
    resp, err := svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
      LogGroupName: aws.String(*group),
      LogStreamNamePrefix: aws.String("app/app"),
      NextToken: nextToken,
    })

    if err != nil {
      log.Fatal(err)
    }

    for _, v := range resp.LogStreams {
      if aws.Int64Value(v.CreationTime) <= startUnixtime && endUnixtime <= aws.Int64Value(v.LastEventTimestamp) {
        log.Print("Found log group: " + aws.StringValue(v.Arn))
        log_events(svc, *group, *v.LogStreamName, startUnixtime, endUnixtime)
      }
    }

    nextToken = resp.NextToken

    if nextToken == nil {
      break
    }
  }

  log.Print("Generated log file: " + file)
}
