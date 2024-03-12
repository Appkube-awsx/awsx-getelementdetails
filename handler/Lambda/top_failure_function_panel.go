package Lambda

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

// FunctionDetails struct to hold details of a function
type FunctionDetails struct {
	FunctionName string
	Timestamp    string
	FailureCount int64
}

var AwsxLambdaFunctionFailureCmd = &cobra.Command{

	Use: "lambda_failure_panel",

	Short: "get lambda failure metrics data",

	Long: `command to get lambda failure metrics data`,

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

			GetTotalFailureFunctionsPanel(cmd, clientAuth, nil)

		}

	},
}

func GetTotalFailureFunctionsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) {
	logGroupName := "CloudTrail/DefaultLogGroup"

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

	// Get total failure count
	totalFailureCount, err := getTotalFailureCount(clientAuth, startTime, endTime, logGroupName, filterPattern, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting total failure count: ", err)
		// handle error
	}

	fmt.Printf("Total Failure Count for All Functions: %d\n", totalFailureCount)

	// Get top failure functions
	topFunctions, err := getTopFailureFunctions(clientAuth, startTime, endTime, logGroupName, filterPattern, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting top failure functions: ", err)
		// handle error
	}

	// Display top functions and individual function details
	fmt.Println("Top Failure Functions:")
	for _, function := range topFunctions {
		fmt.Printf("Function Name: %s, Time: %s, Failure Count: %d\n", function.FunctionName, function.Timestamp, function.FailureCount)
	}
}

func getTotalFailureCount(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, filterPattern string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (int64, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000), // Convert to milliseconds
		EndTime:      aws.Int64(endTime.Unix() * 1000),   // Convert to milliseconds
		QueryString: aws.String(`fields @timestamp, @message
		| filter eventSource=="lambda.amazonaws.com"
		| filter @message like /ERROR|Exception|Failed/
		| stats count(*) as FailureCount`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return 0, fmt.Errorf("failed to start query: %v", err)
	}

	queryId := queryResult.QueryId
	queryStatus := ""
	var queryResults *cloudwatchlogs.GetQueryResultsOutput

	for queryStatus != "Complete" {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}
		queryResults, err = cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return 0, fmt.Errorf("failed to get query results: %v", err)
		}
		queryStatus = aws.StringValue(queryResults.Status)
		time.Sleep(1 * time.Second)
	}

	// Extract total failure count
	var totalFailureCount int64
	for _, resultRow := range queryResults.Results {
		for _, resultField := range resultRow {
			if aws.StringValue(resultField.Field) == "FailureCount" {
				// Correctly convert string to int64
				totalFailureCount, err = strconv.ParseInt(aws.StringValue(resultField.Value), 10, 64)
				if err != nil {
					return 0, fmt.Errorf("failed to parse FailureCount: %v", err)
				}
			}
		}
	}

	return totalFailureCount, nil
}

func getTopFailureFunctions(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, filterPattern string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*FunctionDetails, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000), // Convert to milliseconds
		EndTime:      aws.Int64(endTime.Unix() * 1000),   // Convert to milliseconds
		QueryString: aws.String(`fields @timestamp, @message
		| filter eventSource=="lambda.amazonaws.com"
		| filter @message like /ERROR|Exception|Failed/
		| stats count(*) as FailureCount by requestParameters.functionName
		| sort -FailureCount
		| limit 10`),
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
	var queryResults *cloudwatchlogs.GetQueryResultsOutput

	for queryStatus != "Complete" {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}
		queryResults, err = cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}
		queryStatus = aws.StringValue(queryResults.Status)
		time.Sleep(1 * time.Second)
	}

	// Extract failure functions details
	var functionDetailsList []*FunctionDetails
	for _, resultRow := range queryResults.Results {
		functionDetails := &FunctionDetails{}
		for _, resultField := range resultRow {
			switch aws.StringValue(resultField.Field) {
			case "requestParameters.functionName":
				functionDetails.FunctionName = aws.StringValue(resultField.Value)
			case "@timestamp":
				// Parse timestamp with a custom layout
				timestamp, err := time.Parse("2006-01-02 15:04:05.999", aws.StringValue(resultField.Value))
				if err != nil {
					return nil, fmt.Errorf("failed to parse timestamp: %v", err)
				}
				functionDetails.Timestamp = timestamp.Format(time.RFC3339)
			case "FailureCount":
				// Correctly convert string to int64
				functionDetails.FailureCount, err = strconv.ParseInt(aws.StringValue(resultField.Value), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse FailureCount: %v", err)
				}
			}
		}
		functionDetailsList = append(functionDetailsList, functionDetails)
	}

	return functionDetailsList, nil
}
func init() {
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaFunctionFailureCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
