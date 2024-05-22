package Lambda

import (
	"fmt"
	// "github.com/Appkube-awsx/awsx-common/cmdb"
	"log"
	"strconv"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxLambdaInvocationTrendCmd = &cobra.Command{

	Use:   "invocation_trend_panel",
	Short: "Get invocation trend metrics data",
	Long:  `Command to get invocation trend metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running invocation trend panel command")

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
			panel, err := GetInvocationTrendData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetInvocationTrendData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventSource| filter eventSource = "lambda.amazonaws.com"| stats count() as InvocationCount by bin(1h)`, cloudWatchLogs)
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
					if *data.Field == "InvocationCount" {
						invocationCount, err := strconv.Atoi(*data.Value)
						if err != nil {
							log.Println("Failed to convert InvocationCount to integer:", err)
							continue
						}
						log.Printf("Invocation Count: %d\n", invocationCount)

						// You can perform further processing or store the invocation count data as needed
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
	comman_function.InitAwsCmdFlags(AwsxLambdaInvocationTrendCmd)
}
