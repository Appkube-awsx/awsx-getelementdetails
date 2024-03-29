package Lambda

import (
	"fmt"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"

	"github.com/Appkube-awsx/awsx-common/model"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/spf13/cobra"

	"log"

	"time"
)

var AwsxLambdaFunctionCmd = &cobra.Command{

	Use: "lambda_function_panel",

	Short: "get function metrics data",

	Long: `command to get function metrics data`,

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

			GetFunctionPanel(cmd, clientAuth, nil)

		}

	},
}

func GetFunctionPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) {

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

	events, err := filterCloudWatchLogs(clientAuth, startTime, endTime, logGroupName, filterPattern, cloudWatchLogs)

	if err != nil {

		log.Println("Error in getting sample count: ", err)

		// handle error

	}

	for _, event := range events {

		fmt.Println(event)

	}

}

func filterCloudWatchLogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, filterPattern string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.ResultField, error) {

	// Construct input parameters

	params := &cloudwatchlogs.StartQueryInput{

		LogGroupName: aws.String(logGroupName),

		StartTime: aws.Int64(startTime.Unix() * 1000), // Convert to milliseconds

		EndTime: aws.Int64(endTime.Unix() * 1000), // Convert to milliseconds

		QueryString: aws.String(`fields @timestamp, @message
	   | filter eventSource=="lambda.amazonaws.com"
	   | filter eventName=="GetPolicy20150331"
	   | stats count(*) as functionCount by bin(1mo)

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
	 AwsxLambdaFunctionCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("elementId", "", "element id")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("vaultToken", "", "vault token")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("accountId", "", "aws account number")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("zone", "", "aws region")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("externalId", "", "aws external id")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("elementType", "", "element type")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("instanceId", "", "instance id")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("clusterName", "", "cluster name")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("query", "", "query")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("startTime", "", "start time")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("endTime", "", "endcl time")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	 AwsxLambdaFunctionCmd.PersistentFlags().String("logGroupName", "", "log group name")
}

