package ECS

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
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
			deletedEvents, err := GetECSResourceDeletedEvents(cmd, clientAuth)
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

func GetECSResourceDeletedEvents(cmd *cobra.Command, clientAuth *model.Auth) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return nil, err
		}
		logGroupName = cmdbData.LogGroup
		fmt.Println("HELLPPP", logGroupName)

	}
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime, err := parseTimeRange(startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	deletedEvents, err := FilterDeletedEvents(clientAuth, startTime, endTime, logGroupName)
	if err != nil {
		return nil, err
	}

	return deletedEvents, nil
}

func parseTimeRange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
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

func FilterDeletedEvents(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	cloudWatchLogs := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)

	queryString := `fields @timestamp, eventName
	| filter eventSource = "ecs.amazonaws.com" and (eventName = "DeleteCluster" or eventName = "DeregisterContainerInstance" or eventName = "DeleteService" or eventName = "DeleteTaskSet" or eventName = "DeregisterTaskDefinition" or eventName = "StopTask")
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
				queryResults = append(queryResults, result)
				fmt.Println(r)

			}
		}
		break
	}
	return queryResults, nil
}

func init() {
	AwsxResourceDeletedPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxResourceDeletedPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceDeletedPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
