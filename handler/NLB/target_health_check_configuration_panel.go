package NLB

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxNLBTargetHealthCheckCmd = &cobra.Command{

	Use:   "target_health_check_configuration_panel",
	Short: "Get target health check configuration logs data",
	Long:  `Command to get target health check configuration logs data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running target health check configuration panel command")

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
			panel, err := GetNLBTargetHealthCheckData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetNLBTargetHealthCheckData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
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
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	results, err := FilterCloudWatchLogss(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQueryResult(results)

	return processedResults, nil

}

func FilterCloudWatchLogss(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp
		| filter eventSource = "elasticloadbalancing.amazonaws.com"
		| filter eventName= "CreateTargetGroup" 
		| display responseElements.targetGroups.0.healthCheckProtocol,responseElements.targetGroups.0.healthCheckPort,responseElements.targetGroups.0.healthCheckPath,responseElements.targetGroups.0.healthCheckTimeoutSeconds,responseElements.targetGroups.0.healthCheckIntervalSeconds,responseElements.targetGroups.0.unhealthyThresholdCount,responseElements.targetGroups.0.healthyThresholdCount`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)

	}
	queryId := queryResult.QueryId
	var queryResults []*cloudwatchlogs.GetQueryResultsOutput

	for {
		// Check query status
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryResults = append(queryResults, queryResult)

		if *queryResult.Status != "Complete" {
			time.Sleep(5 * time.Second) // wait before querying again
			continue
		}

		break // exit loop if query is complete
	}
	return queryResults, nil
}

func ProcessQueryResult(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "protocol" {

						log.Printf("protocol: %s\n", *data)

					} else if *data.Field == "port" {

						log.Printf("port: %s\n", *data)

					} else if *data.Field == "path" {

						log.Printf("path: %s\n", *data)

					} else if *data.Field == "timeout" {

						log.Printf("timeout: %s\n", *data)

					} else if *data.Field == "interval" {

						log.Printf("interval: %s\n", *data)

					} else if *data.Field == "unhealthyThreshold" {

						log.Printf("unhealthyThreshold: %s\n", *data)

					} else if *data.Field == "healthyThreshold" {

						log.Printf("healthyThreshold: %s\n", *data)
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
	AwsxNLBTargetHealthCheckCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxNLBTargetHealthCheckCmd.PersistentFlags().String("clusterName", "", "ECS cluster name")
	AwsxNLBTargetHealthCheckCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBTargetHealthCheckCmd.PersistentFlags().String("endTime", "", "end time")
}
