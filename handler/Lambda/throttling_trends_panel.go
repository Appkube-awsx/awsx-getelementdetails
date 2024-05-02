package Lambda

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxLambdaThrottlingTrendsCmd = &cobra.Command{

	Use:   "throttling_trends_panel",
	Short: "Get throttling trends metrics data",
	Long:  `Command to get throttling trends metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running throttling trends panel command")

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
			panel, err := GetThrottlingTrendsData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetThrottlingTrendsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	// elementId, _ := cmd.PersistentFlags().GetString("elementId")
	// cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = commanFunction.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	// if elementId != "" {
	// 	log.Println("getting cloud-element data from cmdb")
	// 	apiUrl := cmdbApiUrl
	// 	if cmdbApiUrl == "" {
	// 		log.Println("using default cmdb url")
	// 		apiUrl = config.CmdbUrl
	// 	}
	// 	log.Println("cmdb url: " + apiUrl)

	// }

	// startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	// endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	// var startTime, endTime time.Time

	// // Parse start time if provided
	// if startTimeStr != "" {
	// 	parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
	// 	if err != nil {
	// 		log.Printf("Error parsing start time: %v", err)
	// 		err := cmd.Help()
	// 		if err != nil {
	// 			// handle error
	// 		}
	// 		return nil, err
	// 	}
	// 	startTime = parsedStartTime

	// } else {
	// 	defaultStartTime := time.Now().Add(-5 * time.Minute)
	// 	startTime = defaultStartTime
	// }

	// if endTimeStr != "" {
	// 	parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
	// 	if err != nil {
	// 		log.Printf("Error parsing end time: %v", err)
	// 		err := cmd.Help()
	// 		if err != nil {
	// 			// handle error
	// 		}
	// 		return nil, err
	// 	}
	// 	endTime = parsedEndTime
	// } else {
	// 	defaultEndTime := time.Now()
	// 	endTime = defaultEndTime
	// }

	results, err := commanFunction.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, InvocationCount, errorCount| filter eventSource = "lambda.amazonaws.com"| stats count() as InvocationCount, count(errorCode) as errorCount by bin(1m)`, cloudWatchLogs)
	if err != nil {
		return nil, err
	}
	processedResults := ProcessQueryResults(results)

	return processedResults, nil

}

// func filterCloudWatchLogsss(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
// 	params := &cloudwatchlogs.StartQueryInput{
// 		LogGroupName: aws.String(logGroupName),
// 		StartTime:    aws.Int64(startTime.Unix() * 1000),
// 		EndTime:      aws.Int64(endTime.Unix() * 1000),
// 		QueryString: aws.String(`fields @timestamp, InvocationCount, errorCount
// 		| filter eventSource = "lambda.amazonaws.com"
// 		| stats count() as InvocationCount, count(errorCode) as errorCount by bin(1m)`),
// 	}

// 	if cloudWatchLogs == nil {
// 		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
// 	}

// 	queryResult, err := cloudWatchLogs.StartQuery(params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to start query: %v", err)

// 	}
// 	queryId := queryResult.QueryId
// 	var queryResults []*cloudwatchlogs.GetQueryResultsOutput

// 	for {
// 		// Check query status
// 		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
// 			QueryId: queryId,
// 		}

// 		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get query results: %v", err)
// 		}

// 		queryResults = append(queryResults, queryResult)

// 		if *queryResult.Status != "Complete" {
// 			time.Sleep(5 * time.Second) // wait before querying again
// 			continue
// 		}

// 		break // exit loop if query is complete
// 	}
// 	return queryResults, nil
// }

func ProcessQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "InvocationCount" {
						invocationCount, err := strconv.Atoi(*data.Value)
						if err != nil {
							log.Println("Failed to convert InvocationCount to integer:", err)
							continue
						}
						log.Printf("Invocation Count: %d\n", invocationCount)

						// You can perform further processing or store the invocation count data as needed
					} else if *data.Field == "errorCount" {
						errorCount, err := strconv.Atoi(*data.Value)
						if err != nil {
							log.Println("Failed to convert errorCount to integer:", err)
							continue
						}
						log.Printf("Error Count: %d\n", errorCount)
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
	AwsxLambdaThrottlingTrendsCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxLambdaThrottlingTrendsCmd.PersistentFlags().String("functionName", "", "Lambda function name")
	AwsxLambdaThrottlingTrendsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaThrottlingTrendsCmd.PersistentFlags().String("endTime", "", "end time")
}
