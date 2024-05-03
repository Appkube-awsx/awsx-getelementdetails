package ECS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxFailedServicePanelCmd = &cobra.Command{
	Use:   "failed_services_panel",
	Short: "Get ECS failed services events",
	Long:  `Command to retrieve ECS failed services events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ECS failed services command")

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
			panel, err := GetECSFailedServiceEvents(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)

		}
	},
}

func GetECSFailedServiceEvents(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource = "ecs.amazonaws.com" and @message like /ERROR|Exception|Failed/ and @message like /service/| stats count() as FailedServiceCount by @timestamp| sort @timestamp desc`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := processQuerysResult(results)

	return processedResults, nil
}

func processQuerysResult(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "failed" {
						log.Printf("failed: %s\n", *data)
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
	AwsxFailedServicePanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxFailedServicePanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxFailedServicePanelCmd.PersistentFlags().String("endTime", "", "end time")
}
