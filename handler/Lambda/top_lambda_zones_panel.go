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

var AwsxLambdaTopLambdaZonesCmd = &cobra.Command{

	Use:   "top_lambda_zones_panel",
	Short: "Get top 5 Lambda zones, event sources, and function names",
	Long:  `Command to get top 5 Lambda zones along with their event sources and function names.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running top lambda zones panel command")

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
			panel, err := GetTopLambdaZonesData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetTopLambdaZonesData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields eventSource, requestParameters.functionName|filter eventSource = "lambda.amazonaws.com"| stats count(*) as EventCount by eventSource, requestParameters.functionName, awsRegion| sort EventCount desc| limit 5`, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	processedResults := ProcesssQueryResults(results)

	return processedResults, nil

}

func ProcesssQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				// Check if there are enough fields in resultField
				if len(resultField) >= 3 {
					region := *resultField[1].Value
					eventSource := *resultField[2].Value
					if len(resultField) >= 4 {
						functionName := *resultField[3].Value
						log.Printf("Region: %s, Event Source: %s, Function Name: %s\n", region, eventSource, functionName)
					} else {
						log.Println("Not enough fields in resultField to extract function name.")
					}
				} else {
					log.Println("Not enough fields in resultField.")
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
	comman_function.InitAwsCmdFlags(AwsxLambdaTopLambdaZonesCmd)
}
