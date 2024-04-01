package Lambda

import (
	"encoding/json"
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
		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			return
		}
		if !authFlag {
			log.Println("Authentication failed")
			return
		}

		responseType, _ := cmd.PersistentFlags().GetString("responseType")
		if responseType != "json" && responseType != "frame" {
			log.Println("Invalid response type. Valid options are 'json' or 'frame'.")
			return
		}

		jsonResp, cloudwatchMetricResp, err := GetLambdaErrorAndWarningData(cmd, clientAuth, nil)
		if err != nil {
			log.Println("Error getting Lambda error and warning data: ", err)
			return
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaErrorAndWarningData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (string, map[string]*cloudwatchlogs.GetQueryResultsOutput, error) {
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
			return "",nil, err
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
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := make(map[string]*cloudwatchlogs.GetQueryResultsOutput)

	// Fetch raw data
	rawData, err := GetLambdaErrorAndWarningMetricData(clientAuth,logGroupName, startTime, endTime, cloudWatchLogs)
	if err != nil {
		return "", nil, err
	}

	// Filter out unwanted fields and keep only "bin(1h)" and "TotalErrors"
	filteredResults := filterResults(rawData.Results)

	// Create a new GetQueryResultsOutput with filtered results
	filteredOutput := &cloudwatchlogs.GetQueryResultsOutput{
		Results:   filteredResults,
		Statistics: rawData.Statistics,
		Status:    rawData.Status,
	}

	cloudwatchMetricData["ErrorAndWarningEvents"] = filteredOutput

	// Generate JSON response
	jsonString, err := json.Marshal(filteredOutput)
	if err != nil {
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, err
}

func filterResults(results [][]*cloudwatchlogs.ResultField) [][]*cloudwatchlogs.ResultField {
	var filteredResults [][]*cloudwatchlogs.ResultField
	for _, row := range results {
		var hasTotalErrors bool
		for _, field := range row {
			if *field.Field == "TotalErrors" {
				hasTotalErrors = true
				break
			}
		}
		if hasTotalErrors {
			var filteredRow []*cloudwatchlogs.ResultField
			for _, field := range row {
				if *field.Field == "bin(1h)" || *field.Field == "TotalErrors" {
					filteredRow = append(filteredRow, field)
				}
			}
			filteredResults = append(filteredResults, filteredRow)
		}
	}
	return filteredResults
}


func GetLambdaErrorAndWarningMetricData(clientAuth *model.Auth, logGroupName string, startTime, endTime *time.Time, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, errorCode
| filter eventSource = 'lambda.amazonaws.com' and (errorCode != '')
| stats count(*) as TotalEvents, count(errorCode) as TotalErrors by errorCode, bin(1h)
| sort @timestamp asc`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	// Start the query
	startQueryOutput, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	// Get the query ID
	queryId := startQueryOutput.QueryId

	// Wait for the query to complete
	for {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResults, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		// Check if query is complete
		if *queryResults.Status == cloudwatchlogs.QueryStatusComplete {
			return queryResults, nil
		}

		// If query is not complete, wait for some time before checking again
		time.Sleep(5 * time.Second)
	}
}



func init() {
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("startTime", "", "Start time in RFC3339 format, e.g., 2024-02-20T00:00:00Z")
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("endTime", "", "End time in RFC3339 format, e.g., 2024-03-26T23:59:59Z")
	
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
