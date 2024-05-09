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
	
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource = "lambda.amazonaws.com"| filter @message like /Dead|Deadletter queue/| stats count(*) as DeadLetterErrorCount by bin(1h)`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQuerysResult(results)

	return processedResults, nil
}

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
	comman_function.InitAwsCmdFlags(AwsxDeadLetterErrorsTrendsCmd)
}
