package EC2

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxEc2InstanceStopCmd = &cobra.Command{

	Use:   "instance_stop_count_panel",
	Short: "Get instance stop count metrics data",
	Long:  `Command to get instance stop count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running instance stop count panel command")

		var authFlag bool
		var clientAuth *model.Auth
		var err error
		authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)

		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			cloudwatchMetricData := GetInstanceStopCountPanel(cmd, clientAuth, nil)
			fmt.Println(cloudwatchMetricData)
		}
	},
}

func GetInstanceStopCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) map[string][]*cloudwatchlogs.GetQueryResultsOutput {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	filterPattern, _ := cmd.PersistentFlags().GetString("filterPattern")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	cloudwatchMetricData := make(map[string][]*cloudwatchlogs.GetQueryResultsOutput)

	events, err := filterCloudWatchLogss(clientAuth, startTime, endTime, logGroupName, filterPattern, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		// handle error
	}
	for _, event := range events {
		fmt.Println(event)
	}
	cloudwatchMetricData["Instance_Stop_Count"] = events
	return cloudwatchMetricData
}

func filterCloudWatchLogss(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName, filterPattern string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	// Construct input parameters
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, @message
            | filter eventSource=="ec2.amazonaws.com"
            | filter eventName=="StopInstances"
            | stats count(*) as InstanceCount by bin(1mo)
            | sort @timestamp desc`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	queryId := queryResult.QueryId
	var queryResults []*cloudwatchlogs.GetQueryResultsOutput

	for {
		// Check query status
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryResults = append(queryResults, queryResult)

		if *queryResult.Status != "Complete" {
			time.Sleep(5 * time.Second) // wait before querying again
			continue
		}

		break // exit loop if query is complete
	}

	return queryResults, nil
}

func init() {
	AwsxEc2InstanceStopCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxEc2InstanceStopCmd.PersistentFlags().String("filterPattern", "", "filter pattern")
	AwsxEc2InstanceStopCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceStopCmd.PersistentFlags().String("endTime", "", "end time")
}
