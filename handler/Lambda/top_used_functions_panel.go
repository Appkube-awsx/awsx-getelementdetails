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

var AwsxTopUsedFunctionsLogCmd = &cobra.Command{

	Use:   "top_used_functions_panel",
	Short: "Get top used functions data",
	Long:  `Command to get top used functions data`,

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
			panel, err := GetTopUsedFunctionsLogData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetTopUsedFunctionsLogData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	// elementId, _ := cmd.PersistentFlags().GetString("elementId")
	// cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
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
	| filter eventName=="GetFunction20150331v2" and  requestParameters.functionName != ""
	| stats count(*) as InvocationCount by requestParameters.functionName ,eventName, eventTime
	| sort by InvocationCount desc
	| limit 10
	`, cloudWatchLogs)
	if err != nil{
		return nil, nil
	}
	processedResults := ProcessQuerysResultss(results)

	return processedResults, nil

}

func ProcessQuerysResultss(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
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

					} else if *data.Field == "InvocationCount" {

						log.Printf("InvocationCount: %s\n", *data)

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
	AwsxTopUsedFunctionsLogCmd.PersistentFlags().String("logGroupName", "", "log group name")
    AwsxTopUsedFunctionsLogCmd.PersistentFlags().String("functionName", "", "Lambda function name")
	AwsxTopUsedFunctionsLogCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxTopUsedFunctionsLogCmd.PersistentFlags().String("endTime", "", "end time")
}
