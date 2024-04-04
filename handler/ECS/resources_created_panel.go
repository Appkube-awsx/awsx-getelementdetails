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
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
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

	}

	startTime, endTime, err := parseTimerange(startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	createdEvents, err := FilterCreatedEvents(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return nil, err
	}

	return createdEvents, nil
}

func FilterCreatedEvents(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}
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

		queryResults = append(queryResults, result)

		break
	}

	return queryResults, nil
}

func init() {
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxResourceCreatedPanelCmd.PersistentFlags().String("endTime", "", "end time")
}
