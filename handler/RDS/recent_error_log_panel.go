package RDS

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxRdsErrorLogsCmd = &cobra.Command{
	Use:   "rds_error_logs",
	Short: "Get recent error logs for RDS instances",
	Long:  `Command to retrieve recent error logs for RDS instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")

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
			_, err := GetRdsErrorLogsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Printf("Error retrieving RDS error logs: %v\n", err)
				return
			}
			// fmt.Println(panel) // Not printing panel directly to demonstrate the custom processing
		}
	},
}

type RdsErrorLogEntry struct {
	Timestamp   time.Time
	ErrorType   string
	ErrorCode   string
	Description string
	Resolution  string
}

func GetRdsErrorLogsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]RdsErrorLogEntry, error) {
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
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return nil, err
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
				return nil, err
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	results, err := filterCloudWatchLogsRDS(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	processedResults := processQueryResults(results)
	return processedResults, nil
}

func filterCloudWatchLogsRDS(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`filter @logStream = 'postgresql.0' and @message like /ERROR/
| fields @timestamp, @message
| limit 20`),
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
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryResults = append(queryResults, queryResult)

		if *queryResult.Status != "Complete" {
			log.Println("Query status is not complete.")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	return queryResults, nil
}

func processQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []RdsErrorLogEntry {
	errorLogs := make([]RdsErrorLogEntry, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, res := range result.Results {
				entry := RdsErrorLogEntry{}

				for _, field := range res {
					if *field.Field == "@timestamp" {
						t, err := time.Parse("2006-01-02 15:04:05.000", *field.Value)
						if err != nil {
							log.Printf("Error parsing timestamp: %v", 							err)
							continue
						}
						entry.Timestamp = t
					} else if *field.Field == "@message" {
						// Example parsing logic for message field to extract error type, error code, description, and resolution
						// Assuming message format is: ERROR: type: <type>, code: <code>, description: <description>, resolution: <resolution>
						message := *field.Value
						var err error
						_, err = fmt.Sscanf(message, "ERROR: type: %s, code: %s, description: %s, resolution: %s", &entry.ErrorType, &entry.ErrorCode, &entry.Description, &entry.Resolution)
						if err != nil {
							log.Printf("Error parsing message: %v", err)
						}
					}
				}

				errorLogs = append(errorLogs, entry)
			}
		} else {
			log.Println("Query status is not complete.")
		}
	}

	return errorLogs
}

func init() {
	AwsxRdsErrorLogsCmd.PersistentFlags().String("elementId", "", "Element ID")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API URL")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("logGroupName", "", "Log Group Name")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("startTime", "", "Start Time")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("endTime", "", "End Time")
}

