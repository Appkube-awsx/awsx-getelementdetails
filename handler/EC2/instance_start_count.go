package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	"log"
)

var AwsxEc2InstanceStartCmd = &cobra.Command{

	Use: "instance_start_count_panel",

	Short: "get instance start count metrics data",

	Long: `command to get instance start count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("running from child command")

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
			panel, err := GetInstanceStartCountPanel(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)

		}

	},
}

func GetInstanceStartCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = commanFunction.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := commanFunction.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource=="ec2.amazonaws.com"| filter eventName=="StartInstances"| stats count(*) as InstanceCount by bin(1mo)| sort @timestamp desc`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := commanFunction.ProcessQueryResult(results)

	return processedResults, nil
}

func init() {
	AwsxEc2InstanceStartCmd.PersistentFlags().String("logGroupName", "", "log group name")
	AwsxEc2InstanceStartCmd.PersistentFlags().String("filterPattern", "", "filter pattern")
	AwsxEc2InstanceStartCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceStartCmd.PersistentFlags().String("endTime", "", "end time")
}
