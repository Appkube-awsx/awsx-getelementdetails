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

var AwsxResourceUpdatedPanelCmd = &cobra.Command{
	Use:   "resource_updated_panel",
	Short: "Get ECS resource update events",
	Long:  `Command to retrieve ECS resource update events`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ECS resource updated panel command")

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
			updatedEvents, err := GetECSResourceUpdatedEvents(cmd, clientAuth, nil)
			if err != nil {
				log.Fatalf("Error retrieving ECS resource update events: %v", err)
				return
			}
			for _, event := range updatedEvents {
				fmt.Println(event)
			}
		}
	},
}

func GetECSResourceUpdatedEvents(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	updatedEvents, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventName| filter eventSource = "ecs.amazonaws.com" and (eventName = "UpdateCluster" or eventName = "UpdateContainerInstance" or eventName = "UpdateService" or eventName = "UpdateTaskSet" or eventName = "RegisterTaskDefinition")| stats count(*) as EventCount by eventName`, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	return updatedEvents, nil
}

// func parseTimerange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
// 	var startTime, endTime *time.Time

// 	// Parse start time if provided
// 	if startTimeStr != "" {
// 		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("error parsing start time: %v", err)
// 		}
// 		startTime = &parsedStartTime
// 	}

// 	// Parse end time if provided
// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("error parsing end time: %v", err)
// 		}
// 		endTime = &parsedEndTime
// 	}

// 	return startTime, endTime, nil
// }

func init() {
	AwsxResourceUpdatedPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxResourceUpdatedPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceUpdatedPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
