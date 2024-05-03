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

var AwsxECSRegistrationEventsCmd = &cobra.Command{

	Use:   "registration_events_panel",
	Short: "Get registration events logs data",
	Long:  `Command to get registration events logs data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running registration events panel command")

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
			panel, err := GetRegistrationEventsData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetRegistrationEventsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message, @logStream, @log| filter eventSource = "ecs.amazonaws.com"| filter eventName = "RegisterContainerInstance"| display eventTime,awsRegion,requestParameters.cluster,requestParameters.totalResources.0.name,responseElements.containerInstance.ec2InstanceId| sort @timestamp desc| limit 10`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQueryResult(results)

	return processedResults, nil

}

func ProcessQueryResult(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventTime" {

						log.Printf("eventTime: %s\n", *data)

					} else if *data.Field == "region" {

						log.Printf("awsRegion: %s\n", *data)

					} else if *data.Field == "clusterName" {

						log.Printf("clusterName: %s\n", *data)

					} else if *data.Field == "resource" {

						log.Printf("resource: %s\n", *data)
					} else if *data.Field == "instanceId" {

						log.Printf("instanceId: %s\n", *data)
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
	AwsxECSRegistrationEventsCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxECSRegistrationEventsCmd.PersistentFlags().String("clusterName", "", "ECS cluster name")
	AwsxECSRegistrationEventsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSRegistrationEventsCmd.PersistentFlags().String("endTime", "", "end time")
}
