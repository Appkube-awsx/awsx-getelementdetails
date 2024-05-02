package NLB

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxNLBErrorLogCmd = &cobra.Command{

	Use:   "error_log_panel",
	Short: "Get error log logs data",
	Long:  `Command to get error log logs data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running error log panel command")

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
			panel, err := GetNLBErrorLogData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetNLBErrorLogData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = commanFunction.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := commanFunction.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventType, errorMessage| filter eventSource = 'elasticloadbalancing.amazonaws.com'| filter ispresent(errorMessage)| display @timestamp, eventType, errorMessage`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQueryResults(results)

	return processedResults, nil

}

func ProcessQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventTime" {

						log.Printf("eventTime: %s\n", *data)

					} else if *data.Field == "eventType" {

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
	AwsxNLBErrorLogCmd.PersistentFlags().String("logGroupName", "", "log group name")

	AwsxNLBErrorLogCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBErrorLogCmd.PersistentFlags().String("endTime", "", "end time")
}
