package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	"log"
)

var AwsxEc2InstanceTerminatedCountCmd = &cobra.Command{

	Use:   "instance_terminated_count_panel",
	Short: "Get instance terminated count metrics data",
	Long:  `Command to get instance terminated count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running instance terminated count panel command")

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
			panel, err := GetInstanceTerminatedCountPanel(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetInstanceTerminatedCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message,eventTime, eventName, userIdentity.invokedBy, responseElements.instancesSet.items.0.currentState.name
	| filter eventName = "TerminateInstances"
	|display eventTime, eventName, userIdentity.invokedBy,responseElements.instancesSet.items.0.currentState.name`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := comman_function.ProcessQueryResult(results)

	return processedResults, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2InstanceTerminatedCountCmd)
}
