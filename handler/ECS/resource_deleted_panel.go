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

var AwsxResourceDeletedPanelCmd = &cobra.Command{
	Use:   "resource_deleted_panel",
	Short: "Get ECS resource deletion events",
	Long:  `Command to retrieve ECS resource deletion events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ECS resource deleted panel command")

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
			deletedEvents, err := GetECSResourceDeletedEvents(cmd, clientAuth, nil)
			if err != nil {
				log.Fatalf("Error retrieving ECS resource deletion events: %v", err)
				return
			}
			for _, event := range deletedEvents {
				fmt.Println(event)
			}
		}
	},
}

func GetECSResourceDeletedEvents(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	deletedEvents, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventName| filter eventSource = "ecs.amazonaws.com" and (eventName = "DeleteCluster" or eventName = "DeregisterContainerInstance" or eventName = "DeleteService" or eventName = "DeleteTaskSet" or eventName = "DeregisterTaskDefinition" or eventName = "StopTask")| stats count(*) as EventCount by eventName`, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	return deletedEvents, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxResourceDeletedPanelCmd)
}
