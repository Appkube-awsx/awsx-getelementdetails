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

var AwsxEC2ListOfInstancesFailureCmd = &cobra.Command{

	Use:   "list_of_instances_failure_panel",
	Short: "Get list of instances failure logs data",
	Long:  `Command to get list of instances failure logs data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running list of instances failure panel command")

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
			panel, err := GetListOfInstancesFailureData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetListOfInstancesFailureData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource=="ec2.amazonaws.com"| filter  eventName=="RunInstances" and failureCode!=""| filter ispresent(responseElements) or ispresent(failureCode)| stats count() as failureCode by eventName,responseElements.instancesSet.items.0.instanceId,responseElements.instancesSet.items.0.instanceType,responseElements.instancesSet.items.0.placement.availabilityZone,errorMessage`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQueryResults(results)

	return processedResults, nil

}

func ProcessQueryResultss(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventName" {

						log.Printf("eventName: %s\n", *data)

					} else if *data.Field == "responseElements.instancesSet.items.0.instanceId" {

						log.Printf("responseElements.instancesSet.items.0.instanceId: %s\n", *data)

					} else if *data.Field == "responseElements.instancesSet.items.0.instanceType" {

						log.Printf("responseElements.instancesSet.items.0.instanceType: %s\n", *data)

					} else if *data.Field == "responseElements.instancesSet.items.0.placement.availabilityZone" {

						log.Printf("responseElements.instancesSet.items.0.placement.availabilityZone: %s\n", *data)
					} else if *data.Field == "errorMessage" {

						log.Printf("errorMessage: %s\n", *data)
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
	comman_function.InitAwsCmdFlags(AwsxEC2ListOfInstancesFailureCmd)
}
