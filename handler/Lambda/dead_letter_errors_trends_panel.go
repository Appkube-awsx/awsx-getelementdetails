package Lambda

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxDeadLetterErrorsTrendsCmd = &cobra.Command{
	Use:   "dead_letter_errors_trends_panel",
	Short: "Get error trend events",
	Long:  `Command to retrieve error trend events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running  panel command")

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
			panel, err := GetLambdaDeadLetterErrorsTrendsEvents(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)

		}
	},
}

func GetLambdaDeadLetterErrorsTrendsEvents(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	// elementId, _ := cmd.PersistentFlags().GetString("elementId")
	// cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	// if elementId != "" {
	// 	log.Println("getting cloud-element data from cmdb")
	// 	apiUrl := cmdbApiUrl
	// 	if cmdbApiUrl == "" {
	// 		log.Println("using default cmdb url")
	// 		apiUrl = config.CmdbUrl
	// 	}
	// 	log.Println("cmdb url: " + apiUrl)
	// 	cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	logGroupName = cmdbData.LogGroup

	// }
	// startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	// endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	// var startTime, endTime *time.Time

	// // Parse start time if provided
	// if startTimeStr != "" {
	// 	parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
	// 	if err != nil {
	// 		log.Printf("Error parsing start time: %v", err)
	// 		err := cmd.Help()
	// 		if err != nil {

	// 		}
	// 	}
	// 	startTime = &parsedStartTime
	// } else {
	// 	defaultStartTime := time.Now().Add(-5 * time.Minute)
	// 	startTime = &defaultStartTime
	// }

	// if endTimeStr != "" {
	// 	parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
	// 	if err != nil {
	// 		log.Printf("Error parsing end time: %v", err)
	// 		err := cmd.Help()
	// 		if err != nil {
	// 			// handle error
	// 		}
	// 	}
	// 	endTime = &parsedEndTime
	// } else {
	// 	defaultEndTime := time.Now()
	// 	endTime = &defaultEndTime
	// }
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource = "lambda.amazonaws.com"| filter @message like /Dead|Deadletter queue/| stats count(*) as DeadLetterErrorCount by bin(1h)`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQuerysResult(results)

	return processedResults, nil
}

// func FilterDeadLetterErrorsTrendsTasks(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
// 	params := &cloudwatchlogs.StartQueryInput{
// 		LogGroupName: aws.String(logGroupName),
// 		StartTime:    aws.Int64(startTime.Unix() * 1000),
// 		EndTime:      aws.Int64(endTime.Unix() * 1000),
// 		QueryString: aws.String(`fields @timestamp, @message
// 		| filter eventSource = "lambda.amazonaws.com"
// 		| filter @message like /Dead|Deadletter queue/
// 		| stats count(*) as DeadLetterErrorCount by bin(1h)`),
// 	}
// 	if cloudWatchLogs == nil {
// 		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
// 	}

// 	queryResult, err := cloudWatchLogs.StartQuery(params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to start query: %v", err)
// 	}

// 	queryId := queryResult.QueryId
// 	var queryResults []*cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop
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

//		return queryResults, nil
//	}
func ProcessQuerysResult(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "DeadLetterErrorCount" {

						log.Printf("DeadLetterErrorCount: %s\n", *data)

						// You can perform further processing or store the instance count data as needed
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
	AwsxDeadLetterErrorsTrendsCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxDeadLetterErrorsTrendsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxDeadLetterErrorsTrendsCmd.PersistentFlags().String("endTime", "", "end time")
}
