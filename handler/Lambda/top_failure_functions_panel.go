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

var AwsxTopFailureFunctionsLogCmd = &cobra.Command{

	Use:   "top_failure_functions_panel",
	Short: "Get top failure functions data",
	Long:  `Command to get top failure functions data`,

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
			panel, err := GetTopFailureFunctionsLogData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetTopFailureFunctionsLogData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
    
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	
	results, err :=  comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName,`fields @timestamp, @message
	| filter eventSource=="lambda.amazonaws.com"
	| filter @message like /ERROR|Exception|Failed/
	| stats count(*) as FailureCount by requestParameters.functionName,eventTime,eventName
	| limit 10
	`, cloudWatchLogs)
	if err != nil{
		return nil, nil
	}
	processedResults := ProcessQuerysResults(results)

	return processedResults, nil

}

func ProcessQuerysResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "requestParameters.functionName" {

						log.Printf("requestParameters.functionName: %s\n", *data)

					} else if *data.Field == "eventName" {

						log.Printf("eventName: %s\n", *data)

					} else if *data.Field == "eventTime" {

						log.Printf("eventTime: %s\n", *data)

					} else if *data.Field == "FailureCount" {

						log.Printf("FailureCount: %s\n", *data)

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
	AwsxTopFailureFunctionsLogCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxTopFailureFunctionsLogCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxTopFailureFunctionsLogCmd.PersistentFlags().String("endTime", "", "end time")
}
