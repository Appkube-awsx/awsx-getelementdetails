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

var AwsxActiveTaskPanelCmd = &cobra.Command{
	Use:   "active_task_panel",
	Short: "Get ECS active task events",
	Long:  `Command to retrieve ECS active task events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ECS active task panel command")

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
			panel, err := GetECSActiveTaskEvents(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)

		}
	},
}

func GetECSActiveTaskEvents(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource = "ecs.amazonaws.com" and @message like /task/ and not(@message like /ERROR|Exception|Failed/)| stats count() as ActiveTaskCount by @timestamp| sort @timestamp desc`, cloudWatchLogs)
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
	comman_function.InitAwsCmdFlags(AwsxActiveTaskPanelCmd)
}
