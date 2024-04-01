package ECS

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
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
			updatedEvents, err := GetECSResourceUpdatedEvents(cmd, clientAuth)
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

func GetECSResourceUpdatedEvents(cmd *cobra.Command, clientAuth *model.Auth) ([]*cloudwatchlogs.ResultField, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime, err := parseTimerange(startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	updatedEvents, err := FilterUpdatedEvents(clientAuth, startTime, endTime, logGroupName)
	if err != nil {
		return nil, err
	}

	return updatedEvents, nil
}

func parseTimerange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing start time: %v", err)
		}
		startTime = &parsedStartTime
	}

	// Parse end time if provided
	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing end time: %v", err)
		}
		endTime = &parsedEndTime
	}

	return startTime, endTime, nil
}

func FilterUpdatedEvents(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string) ([]*cloudwatchlogs.ResultField, error) {
	cloudWatchLogs := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)

	queryString := `fields @timestamp, eventName
	| filter eventSource = "ecs.amazonaws.com" and (eventName = "UpdateCluster" or eventName = "UpdateContainerInstance" or eventName = "UpdateService" or eventName = "UpdateTaskSet" or eventName = "RegisterTaskDefinition")
	| stats count(*) as EventCount by eventName`

	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString:  aws.String(queryString),
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	queryId := queryResult.QueryId
	var queryResults []*cloudwatchlogs.ResultField

	for {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		result, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		if *result.Status != "Complete" {
			time.Sleep(5 * time.Second)
			continue
		}

		// Flatten and append each element individually
		for _, res := range result.Results {
			for _, r := range res {
				queryResults = append(queryResults, r)
			}
		}

		break
	}
	return queryResults, nil
}

func init() {
	AwsxResourceUpdatedPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxResourceUpdatedPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceUpdatedPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
