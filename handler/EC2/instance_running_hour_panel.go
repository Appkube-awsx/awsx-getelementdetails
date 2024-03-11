package EC2

import (
	"encoding/json"
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

var AwsxEc2InstanceRunningHourCmd = &cobra.Command{

	Use:   "instance_running_hour_panel",
	Short: "get instance running hour metrics data",
	Long:  `command to get instance running hour metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")

		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)

		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			rawEvents, jsonData, err := GetInstanceRunningHoursPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error retrieving instance running hours data:", err)
				return
			}

			// Now you have access to rawEvents and jsonData
			// You can process or use them as needed
			fmt.Println("Raw Events:")
			fmt.Println(rawEvents)
			fmt.Println("JSON Data:", string(jsonData))
		}
	},
}

func GetInstanceRunningHoursPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (string, []byte, error) {
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
				log.Printf("Error parsing start time: %v", err)
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

	events, err := filterCloudwatchLogs(clientAuth, startTime, endTime, logGroupName, filterPattern, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}

	// Convert events to a string
	var rawEventsStr string
	for _, event := range events {
		rawEventsStr += fmt.Sprintf("%+v\n", event)
	}

	// Convert events to JSON
	jsonData, err := json.Marshal(events)
	if err != nil {
		log.Println("Error marshaling JSON: ", err)
		return "", nil, err
	}

	// Return raw data and JSON data
	return rawEventsStr, jsonData, nil
}

func filterCloudwatchLogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, filterPattern string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.ResultField, error) {
	// Construct input parameters
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, @timestamp 
            | filter eventSource=="ec2.amazonaws.com"
            | filter eventName=="RunInstances"
            | stats sum(responseElements.instancesSet.items.0.launchTime/3600) as totalDurationHours by bin(1month)`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	queryId := queryResult.QueryId
	queryStatus := ""
	var queryResults *cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop
	for queryStatus != "Complete" {
		// Check query status
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResults, err = cloudWatchLogs.GetQueryResults(queryStatusInput) // Assign value to queryResults
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryStatus = aws.StringValue(queryResults.Status)
		time.Sleep(1 * time.Second) // Wait for a second before checking status again
	}

	// Query is complete, now process results
	var results []*cloudwatchlogs.ResultField
	for _, resultRow := range queryResults.Results {
		results = append(results, resultRow...)
	}

	return results, nil
}

func init() {
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
