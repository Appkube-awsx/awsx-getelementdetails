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

var AwsxLambdaErrorAndWarningCmd = &cobra.Command{
	Use:   "error_and_warning_events_panel",
	Short: "Get error and warning events metrics data",
	Long:  `Command to get error and warning events metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running from child command")

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
			panel, err := GetLambdaErrorAndWarningData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetLambdaErrorAndWarningData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp| filter eventSource = 'lambda.amazonaws.com' and (errorCode != '')| stats count(*) as TotalWarnings, count(errorCode) as TotalErrors by bin(1month)| sort @timestamp asc`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := processQuerysResults(results)

	return processedResults, nil
}

func processQuerysResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "TotalErrors" {

						log.Printf("TotalErrors: %s\n", *data)

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
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("startTime", "", "Start time in RFC3339 format, e.g., 2024-02-20T00:00:00Z")
	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("endTime", "", "End time in RFC3339 format, e.g., 2024-03-26T23:59:59Z")

	AwsxLambdaErrorAndWarningCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
