package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var file string
var events []map[string]string

func write(data string) {
	log.Print("Write results...")

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer fp.Close()
	fmt.Fprintln(fp, data)
}

func searchLogEvents(service *cloudwatchlogs.CloudWatchLogs, group string, stream string, start int64, end int64, reg *regexp.Regexp) {
	log.Print("Searching log events...")

	var token *string = nil

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
			size := len(reg.String())

			if size == 0 || size > 0 && reg.Match([]byte(*v.Message)) {
				event := map[string]string{
					"LogStream":     stream,
					"Message":       aws.StringValue(v.Message),
					"Timestamp":     time.Unix(aws.Int64Value(v.Timestamp)/1000, 0).String(),
					"IngestionTime": time.Unix(aws.Int64Value(v.IngestionTime)/1000, 0).String(),
				}
				events = append(events, event)
				log.Print(event)
			}
		}

		if aws.StringValue(token) == aws.StringValue(resp.NextForwardToken) {
			break
		}

		token = resp.NextForwardToken
	}
}

func searchLogGroup(service *cloudwatchlogs.CloudWatchLogs, group string, prefix string, start int64, end int64, reg *regexp.Regexp) {
	log.Print("Searching log groups...")

	var token *string

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
			if aws.Int64Value(v.FirstEventTimestamp) <= start && aws.Int64Value(v.LastEventTimestamp) <= end {
				log.Print("Found log group: " + aws.StringValue(v.Arn))
				searchLogEvents(service, group, *v.LogStreamName, start, end, reg)
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

	group := flag.String("group", "", "Log group name")
	prefix := flag.String("prefix", "", "Prefix name when searching log groups")
	start := flag.String("start", now.Add(-10*time.Minute).Format(layout), "Log stream event search start date and time (UTC)")
	end := flag.String("end", now.Format(layout), "Log stream event search end date and time (UTC)")
	pattern := flag.String("pattern", "", "Regular expressions for events to be extracted")

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

	searchLogGroup(service, *group, *prefix, int64(s.Unix()*1000), int64(e.Unix()*1000), regexp.MustCompile(*pattern))

	if len(events) > 0 {
		data, err := json.MarshalIndent(events, "", "  ")

		if err != nil {
			log.Fatal(err)
		}

		write(string(data))
		log.Print("Generated log file: " + file)
	} else {
		log.Print("Event was not found.")
	}
}
