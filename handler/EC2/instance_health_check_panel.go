package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxEc2InstanceHealthCheckCmd = &cobra.Command{

	Use:   "ec2_instance_health_check_panel",
	Short: "Get ec2 instance health check data",
	Long:  `Command to get ec2 instance health check data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ec2 instance health check panel command")

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
			panel, err := GetEc2InstanceHealthCheckData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetEc2InstanceHealthCheckData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp| filter eventSource=="ec2.amazonaws.com"| filter eventName=="RunInstances"| fields responseElements.instancesSet.items.0.instanceId as instanceId, requestParameters.instanceType as instanceType, responseElements.instancesSet.items.0.launchTime as launchTime, responseElements.instancesSet.items.0.placement.availabilityZone as availabilityZone, responseElements.instancesSet.items.0.instanceState.name as instanceStatus| sort @timestamp desc`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQuerysResultzss(results)

	return processedResults, nil

}

func ProcessQuerysResultzss(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventTime" {

						log.Printf("eventTime: %s\n", *data)

					} else if *data.Field == "instanceId" {

						log.Printf("instanceId: %s\n", *data)

					} else if *data.Field == "instanceType" {

						log.Printf("instanceType: %s\n", *data)

					} else if *data.Field == "launchTime" {

						log.Printf("launchTime: %s\n", *data)

					} else if *data.Field == "availabilityZone" {

						log.Printf("availabilityZone: %s\n", *data)

					} else if *data.Field == "instanceStatus" {

						log.Printf("instanceStatus: %s\n", *data)
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
	comman_function.InitAwsCmdFlags(AwsxEc2InstanceHealthCheckCmd)
}
