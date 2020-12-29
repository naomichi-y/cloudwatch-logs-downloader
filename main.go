package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var file string

func write(data string) {
	log.Print("Write results...")

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer fp.Close()
	fmt.Fprintln(fp, data)
}

func log_events(service *cloudwatchlogs.CloudWatchLogs, group string, stream string, start int64, end int64) {
	log.Print("Searching log events...")

	var token *string = nil
	var result []map[string]string

	for {
		resp, err := service.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(group),
			LogStreamName: aws.String(stream),
			StartTime:     aws.Int64(start),
			EndTime:       aws.Int64(end),
			NextToken:     token,
		})

		if err != nil {
			log.Fatal(err)
		}

		for _, v := range resp.Events {
			r := map[string]string{
				"LogStream": stream,
				"Message":   aws.StringValue(v.Message),
				"Timestamp": time.Unix(aws.Int64Value(v.Timestamp)/1000, 0).String(),
			}
			result = append(result, r)
			log.Print(r)
		}

		if aws.StringValue(token) == aws.StringValue(resp.NextForwardToken) {
			break
		}

		token = resp.NextForwardToken
	}

	data, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	write(string(data))
}

func search_log_group(service *cloudwatchlogs.CloudWatchLogs, group string, prefix string, start int64, end int64) {
	log.Print("Searching log groups...")

	var token *string = nil

	for {
		input := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(group),
			NextToken:    token,
		}

		if prefix != "" {
			input.LogStreamNamePrefix = aws.String(prefix)
		}

		resp, err := service.DescribeLogStreams(input)

		if err != nil {
			log.Fatal(err)
		}

		for _, v := range resp.LogStreams {
			if aws.Int64Value(v.CreationTime) <= start && end <= aws.Int64Value(v.LastEventTimestamp) {
				log.Print("Found log group: " + aws.StringValue(v.Arn))
				log_events(service, group, *v.LogStreamName, start, end)
			}
		}

		token = resp.NextToken

		if token == nil {
			break
		}
	}
}

func main() {
	layout := "2006-01-02 15:04:05"
	now := time.Now()

	group := flag.String("group", "", "Log group")
	prefix := flag.String("prefix", "", "Log group prefix")
	start := flag.String("start", now.Add(-5*time.Minute).Format(layout), "Start date")
	end := flag.String("end", now.Format(layout), "End date")

	flag.Parse()

	s, err := time.Parse(layout, *start)

	if err != nil {
		log.Fatal(err)
	}

	e, err := time.Parse(layout, *end)

	if err != nil {
		log.Fatal(err)
	}

	sess := session.Must(session.NewSession())
	service := cloudwatchlogs.New(
		sess,
		aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")),
	)
	file = "./dist/result_" + time.Now().Format("2006010230405") + ".log"

	search_log_group(service, *group, *prefix, int64(s.Unix()*1000), int64(e.Unix()*1000))

	log.Print("Generated log file: " + file)
}
