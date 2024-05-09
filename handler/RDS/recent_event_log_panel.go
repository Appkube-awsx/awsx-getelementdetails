package RDS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	 "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxRecentEventLogsCmd = &cobra.Command{
	Use:   "recent_event_logs",
	Short: "Get recent event logs",
	Long:  `Command to retrieve recent event logs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")

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
			data, err := GetRecentEventLogsPanel(cmd, clientAuth, nil)
			if err != nil {
				
				return
			}
			fmt.Println(data)

		}
	},
}

// type RecentEventLogEntry struct {
// 	Timestamp       string // Change type to string
// 	EventName       string // No changes here
// 	SourceIPAddress string // No changes here
// 	EventSource     string // No changes here
// 	UserAgent       string // No changes here
// }

func GetRecentEventLogsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)

	

	

	
	if err != nil {
		return nil, fmt.Errorf("Error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName,`fields @timestamp, eventName, sourceIPAddress, eventSource, userAgent| filter eventSource = 'rds.amazonaws.com' | limit 1000`, cloudWatchLogs)


	
	if err != nil {
		return nil, nil
	}

	processedResults := comman_function.ProcessQueryResult(results)

	return processedResults, nil

}



func init() {
	comman_function.InitAwsCmdFlags(AwsxRecentEventLogsCmd)
}
