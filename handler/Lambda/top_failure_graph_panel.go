package Lambda

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxLambdaTopFailureCmd = &cobra.Command{

	Use:   "top_failure_count_panel",
	Short: "Get top failure count metrics data",
	Long:  `Command to get top failure count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running top failure count panel command")

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
			panel, err := GetLambdaTopFailurePanel(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetLambdaTopFailurePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
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

	results, err := filterCloudWatchsLogss(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := processQueryResultss(results)

	return processedResults, nil
}

func filterCloudWatchsLogss(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	// Construct input parameters
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, @message
		| filter eventSource=="lambda.amazonaws.com"
		| filter @message like /ERROR|Exception|Failed/
		| stats count(*) as FailureCount by bin(1m)`),
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
func processQueryResultss(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "TopFailure" {
						topFailureCount, err := strconv.Atoi(*data.Value)
						if err != nil {
							log.Println("Failed to convert TopFailureCount to integer:", err)
							continue
						}
						log.Printf("Top Failure Count: %d\n", topFailureCount)

						// You can perform further processing or store the top failure count data as needed
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
	AwsxLambdaTopFailureCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxLambdaTopFailureCmd.PersistentFlags().String("filterPattern", "", "filter pattern")
	AwsxLambdaTopFailureCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaTopFailureCmd.PersistentFlags().String("endTime", "", "end time")
}
