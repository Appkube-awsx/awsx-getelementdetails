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

var AwsxApiFailedEventCmd = &cobra.Command{

	Use:   "failed_event_panel",
	Short: "Get failed event metrics data",
	Long:  `Command to get failed event metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running failed event panel command")

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
			panel, err := GetFailedEventData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetFailedEventData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message | filter eventSource = 'apigateway.amazonaws.com' | filter ispresent(errorMessage) | display eventType, errorMessage`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQueryResultss(results)

	return processedResults, nil

}

func ProcessQueryResultss(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventType" {

						log.Printf("eventType: %s\n", *data)

					} else if *data.Field == "errorMessage" {

						log.Printf("errorMessage: %s\n", *data)
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
	AwsxApiFailedEventCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxApiFailedEventCmd.PersistentFlags().String("functionName", "", "Lambda function name")
	AwsxApiFailedEventCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiFailedEventCmd.PersistentFlags().String("endTime", "", "end time")
}
