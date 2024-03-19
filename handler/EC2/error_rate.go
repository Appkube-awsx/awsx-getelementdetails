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

var AwsxEc2InstanceErrorRateCmd = &cobra.Command{

	Use:   "instance_error_rate_panel",
	Short: "get instance error rate metrics data",
	Long:  `command to get instance error rate metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")

		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)

		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}

		if authFlag {
			cloudwatchMetricData := GetInstanceErrorRatePanel(cmd, clientAuth, nil)
			fmt.Println(cloudwatchMetricData)
		}
	},
}

func GetInstanceErrorRatePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) map[string][]*cloudwatchlogs.ResultField {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	// filterPattern, _ := cmd.PersistentFlags().GetString("filterPattern")
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
	cloudwatchMetricData := make(map[string][]*cloudwatchlogs.ResultField)

	events, err := filterCloudWatchlogs(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		// handle error
	}
	for _, event := range events {
		fmt.Println(event)
	}
	cloudwatchMetricData["Instance_Error_Rate"] = events
	return cloudwatchMetricData
}

func filterCloudWatchlogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.ResultField, error) {
	// Construct input parameters
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, @message
            | filter eventSource=="ec2.amazonaws.com"
            | filter eventName=="RunInstances" and errorCode!=""
            | stats count(*) as ErrorCount by bin(1d)
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
		for _, resultField := range resultRow {
			results = append(results, resultField)
		}
	}

	return results, nil
}

func init() {
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxEc2InstanceErrorRateCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
