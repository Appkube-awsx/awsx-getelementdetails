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
			createdEvents, err := GetECSResourceCreatedEvents(cmd, clientAuth)
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

func GetECSResourceCreatedEvents(cmd *cobra.Command, clientAuth *model.Auth) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime, err := parseTimerange(startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	createdEvents, err := FilterCreatedEvents(clientAuth, startTime, endTime, logGroupName)
	if err != nil {
		return nil, err
	}

	return createdEvents, nil
}

func FilterCreatedEvents(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	cloudWatchLogs := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)

	queryString := `fields @timestamp, eventName
	| filter eventSource = "ecs.amazonaws.com" and (eventName = "CreateCluster" or eventName = "RegisterContainerInstance" or eventName = "CreateService" or eventName = "RegisterTaskDefinition" or eventName = "CreateTask" or eventName = "RunTask")
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
	var queryResults []*cloudwatchlogs.GetQueryResultsOutput

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
				queryResults = append(queryResults)
				fmt.Println(r)
			}
		}

		break
	}
	return queryResults, nil
}

func init() {
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
