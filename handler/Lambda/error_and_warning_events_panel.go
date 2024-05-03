package Lambda

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxLambdaErrorAndWarningCmd = &cobra.Command{
	Use:   "error_and_warning_events_panel",
	Short: "Get error and warning events metrics data",
	Long:  `Command to get error and warning events metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running from child command")

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
			panel, err := GetLambdaErrorAndWarningData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetLambdaErrorAndWarningData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return nil, err
		}
		logGroupName = cmdbData.LogGroup

	}
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

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)
	results, err := GetLambdaErrorAndWarningMetricData(clientAuth, logGroupName, startTime, endTime, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := processQuerysResults(results)

	return processedResults, nil
}

func GetLambdaErrorAndWarningMetricData(clientAuth *model.Auth, logGroupName string, startTime, endTime *time.Time, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp
		| filter eventSource = 'lambda.amazonaws.com' and (errorCode != '')
		| stats count(*) as TotalWarnings, count(errorCode) as TotalErrors by bin(1month)
		| sort @timestamp asc`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	// Start the query
	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	// Get the query ID
	queryId := queryResult.QueryId
	var queryResults []*cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop
	// Wait for the query to complete
	for {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryResults = append(queryResults, queryResult)

		// Check if query is complete
		if *queryResult.Status != "Complete" {
			time.Sleep(5 * time.Second) // wait before querying again
			continue
		}

		break // exit loop if query is complete
	}

	return queryResults, nil
}

func processQuerysResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "TotalErrors" {

						log.Printf("TotalErrors: %s\n", *data)

					}
				}
			}
			processedResults = append(processedResults, result)

		} else {
			log.Println("Query status is not complete.")
		}
	}

	return processedResults
}

func init() {
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("startTime", "", "Start time in RFC3339 format, e.g., 2024-02-20T00:00:00Z")
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("endTime", "", "End time in RFC3339 format, e.g., 2024-03-26T23:59:59Z")

	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
