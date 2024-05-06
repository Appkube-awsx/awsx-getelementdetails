package ApiGateway

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxApiSuccessEventCmd = &cobra.Command{

	Use:   "successful_event_panel",
	Short: "Get successful event metrics data",
	Long:  `Command to get successful event metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running successful event panel command")

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
			panel, err := GetSuccessEventData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetSuccessEventData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventType, errorMessage| filter eventSource = 'apigateway.amazonaws.com' | filter !ispresent(errorMessage) | display @timestamp, eventType`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := processQueryResults(results)

	return processedResults, nil

}

func processQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventType" {

						log.Printf("eventType: %s\n", *data)

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
	AwsxApiSuccessEventCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxApiSuccessEventCmd.PersistentFlags().String("functionName", "", "Lambda function name")
	AwsxApiSuccessEventCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiSuccessEventCmd.PersistentFlags().String("endTime", "", "end time")
}
