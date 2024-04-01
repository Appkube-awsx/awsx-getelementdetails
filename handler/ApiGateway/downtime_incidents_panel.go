package ApiGateway

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

var AwsxApiDowntimeIncidentsCmd = &cobra.Command{
	Use:   "downtime_incidents",
	Short: "Get downtime incidents data",
	Long:  `Command to get downtime incidents data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running downtime incidents command")

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
			results, err := GetDowntimeIncidentsData(cmd, clientAuth, nil)
			if err != nil {
				log.Printf("Error getting downtime incidents data: %v\n", err)
				return
			}
			// Print the results
			for _, result := range results {
				fmt.Println(result)
			}
		}
	},
}

func GetDowntimeIncidentsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]string, error) {
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
	startTime, endTime := parseStartEndTime(cmd)

	results, err := FilterDowntimeIncidentsLogs(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func parseStartEndTime(cmd *cobra.Command) (*time.Time, *time.Time) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		} else {
			startTime = &parsedStartTime
		}
	}

	// Parse end time if provided
	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
		} else {
			endTime = &parsedEndTime
		}
	}

	return startTime, endTime
}

func FilterDowntimeIncidentsLogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]string, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`
            fields @timestamp, eventType, errorMessage
            | filter eventSource = 'apigateway.amazonaws.com'
            | sort @timestamp desc
        `),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	// Wait for query to complete
	for {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryResult.QueryId,
		}
		queryResults, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		// If query is complete, return results
		if *queryResults.Status == "Complete" {
			// Process the query results to extract lines below and above
			results := processQueryResult(queryResults.Results)
			return results, nil
		}

		// If query is still running, wait before checking again
		time.Sleep(5 * time.Second)
	}
}

func processQueryResult(results [][]*cloudwatchlogs.ResultField) []string {
	var output []string
	for _, event := range results {
		var parsedResult string
		var errorMessage string
		for _, field := range event {
			if *field.Field == "errorMessage" {
				errorMessage = *field.Value
			}
			if *field.Field == "@timestamp" || *field.Field == "eventType" {
				parsedResult += fmt.Sprintf("%s: %s\n", *field.Field, *field.Value)
			}
		}
		if errorMessage != "" {
			parsedResult += fmt.Sprintf("errorMessage: %s\n", errorMessage)
			output = append(output, parsedResult)
		}
	}
	return output
}

func init() {
	AwsxApiDowntimeIncidentsCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxApiDowntimeIncidentsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiDowntimeIncidentsCmd.PersistentFlags().String("endTime", "", "end time")
}
