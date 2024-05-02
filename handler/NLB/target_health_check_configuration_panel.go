package NLB

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
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
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = commanFunction.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := commanFunction.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp| filter eventSource = "elasticloadbalancing.amazonaws.com"| filter eventName= "CreateTargetGroup"| display responseElements.targetGroups.0.healthCheckProtocol,responseElements.targetGroups.0.healthCheckPort,responseElements.targetGroups.0.healthCheckPath,responseElements.targetGroups.0.healthCheckTimeoutSeconds,responseElements.targetGroups.0.healthCheckIntervalSeconds,responseElements.targetGroups.0.unhealthyThresholdCount,responseElements.targetGroups.0.healthyThresholdCount`, cloudWatchLogs)
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

	AwsxNLBTargetHealthCheckCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBTargetHealthCheckCmd.PersistentFlags().String("endTime", "", "end time")
}
