package RDS

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

var AwsxRecentEventLogsCmd = &cobra.Command{
	Use:   "recent_event_logs",
	Short: "Get recent event logs",
	Long:  `Command to retrieve recent event logs`,
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
			jsonResp, rawLogs, err := GetRecentEventLogsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Printf("Error retrieving recent event logs: %v\n", err)
				return
			}
			fmt.Println("JSON Response:")
			fmt.Println(jsonResp)
			fmt.Println("\nRaw Logs:")
			fmt.Println(rawLogs)
		}
	},
}

type RecentEventLogEntry struct {
	Timestamp       string // Change type to string
	EventName       string // No changes here
	SourceIPAddress string // No changes here
	EventSource     string // No changes here
	UserAgent       string // No changes here
}

func GetRecentEventLogsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (string, string, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

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
            return "","", err
        }
        logGroupName = cmdbData.LogGroup
	}

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", "", err
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
				return "", "", err
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	results, err := filterCloudWatchLogRDS(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return "", "", err
	}

	processedResults := processQueryResult(results)
	jsonResp, err := json.Marshal(processedResults)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", "", err
	}

	// Concatenate raw logs
	rawLogs := ""
	for _, log := range processedResults {
		rawLogs += fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", log.Timestamp, log.EventName, log.SourceIPAddress, log.EventSource, log.UserAgent)
	}

	return string(jsonResp), rawLogs, nil
}

func filterCloudWatchLogRDS(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, eventName, sourceIPAddress, eventSource, userAgent
| filter eventSource = 'rds.amazonaws.com' 
| limit 1000`),
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

func processQueryResult(results []*cloudwatchlogs.GetQueryResultsOutput) []RecentEventLogEntry {
	eventLogs := make([]RecentEventLogEntry, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, res := range result.Results {
				entry := RecentEventLogEntry{}

				for _, field := range res {
					if *field.Field == "@timestamp" {
						entry.Timestamp = *field.Value
					}
					if *field.Field == "eventName" {
						entry.EventName = *field.Value
					}
					if *field.Field == "sourceIPAddress" {
						entry.SourceIPAddress = *field.Value
					}
					if *field.Field == "eventSource" {
						entry.EventSource = *field.Value
					}
					if *field.Field == "userAgent" {
						entry.UserAgent = *field.Value
					}
				}

				eventLogs = append(eventLogs, entry)
			}
		}
	}

	return eventLogs
}

func init() {
	AwsxRecentEventLogsCmd.PersistentFlags().String("elementId", "", "Element ID")
	AwsxRecentEventLogsCmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API URL")
	AwsxRecentEventLogsCmd.PersistentFlags().String("logGroupName", "", "Log Group Name")
	AwsxRecentEventLogsCmd.PersistentFlags().String("startTime", "", "Start Time")
	AwsxRecentEventLogsCmd.PersistentFlags().String("endTime", "", "End Time")
}
