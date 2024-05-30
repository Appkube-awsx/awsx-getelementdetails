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

var AwsxEc2InactiveInstanceCmd = &cobra.Command{

	Use:   "inactive_instances_panel",
	Short: "Get inactive instances count metrics data",
	Long:  `Command to get inactive instance count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running instance count panel command")

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
			panel, err := GetInactiveInstancesCountPanel(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetInactiveInstancesCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource=="ec2.amazonaws.com"| filter eventName=="StopInstances"| stats count(*) as InstanceCount`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}

	processedResults := comman_function.ProcessQueryResult(results)

	return processedResults, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2InactiveInstanceCmd)
}
