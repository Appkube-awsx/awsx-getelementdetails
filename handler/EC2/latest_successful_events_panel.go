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

var AwsxEc2LatestSuccessfulEventsCmd = &cobra.Command{

	Use:   "latest_sucessful_events_panel",
	Short: "Get latest successful events metrics data",
	Long:  `Command to get latest successful events metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running latest successful events count panel command")

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
			panel, err := GetLatestSucessfulEventsCountPanel(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetLatestSucessfulEventsCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message, user, userIdentity.sessionContext.sessionIssuer.userName,userIdentity.sessionContext.sessionIssuer.type
	| filter eventSource = "ec2.amazonaws.com" and ispresent(errorCode)
	| display eventTime, eventName, sourceIPAddress,userIdentity.sessionContext.sessionIssuer.userName, userIdentity.sessionContext.sessionIssuer.type`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := comman_function.ProcessQueryResult(results)

	return processedResults, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2LatestSuccessfulEventsCmd)
}
