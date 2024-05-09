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
			// for _, result := range results {
			fmt.Println(results)
			// }
		}
	},
}

func GetDowntimeIncidentsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventType, errorMessage| filter eventSource = 'apigateway.amazonaws.com'| sort @timestamp desc`, cloudWatchLogs)
	if err != nil {
		return nil, err
	}
	processedResult := processQueryResult(results)

	return processedResult, nil
}

func processQueryResult(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResult := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "errorMessage" {

						log.Printf("errorMessage: %s\n", *data)

					}
					if *data.Field == "@timestamp" {

						log.Printf("timestamp: %s\n", *data)

					}
				}
			}
			processedResult = append(processedResult, result)
		} else {
			log.Println("Query status is not complete.")
		}
	}

	return processedResult
}

func init() {
	AwsxApiDowntimeIncidentsCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxApiDowntimeIncidentsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiDowntimeIncidentsCmd.PersistentFlags().String("endTime", "", "end time")
}
