package ECS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxResourceCreatedPanelCmd = &cobra.Command{
	Use:   "resource_created_panel",
	Short: "Get ECS resource creation events",
	Long:  `Command to retrieve ECS resource creation events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ECS resource created panel command")

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
			createdEvents, err := GetECSResourceCreatedEvents(cmd, clientAuth, nil)
			if err != nil {
				log.Fatalf("Error retrieving ECS resource creation events: %v", err)
				return
			}
			for _, event := range createdEvents {
				fmt.Println(event)
			}
		}
	},
}

func GetECSResourceCreatedEvents(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = commanFunction.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	createdEvents, err := commanFunction.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventName| filter eventSource = "ecs.amazonaws.com" and (eventName = "CreateCluster" or eventName = "RegisterContainerInstance" or eventName = "CreateService" or eventName = "RegisterTaskDefinition" or eventName = "CreateTask" or eventName = "RunTask")| stats count(*) as EventCount by eventName`, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	return createdEvents, nil
}

func init() {
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
