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

// Define a struct to hold the log entry data
// type ErrorAnalysisEntry struct {
// 	EventTime    time.Time
// 	EventType    string
// 	EventSource  string
// 	ErrorCode    string
// 	ErrorMessage string
// }

var AwsxRDSErrorAnalysisCmd = &cobra.Command{
	Use:   "error_analysis_panel",
	Short: "Get error analysis panel for RDS instances",
	Long:  `Command to retrieve error analysis panel for RDS instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running error analysis panel command")

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
			data, err := GetErrorAnalysisData(cmd, clientAuth, nil)
			if err != nil {

				return
			}
			fmt.Println(data)
		}
	},
}

// Function to retrieve error analysis panel data
func GetErrorAnalysisData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	// Retrieve necessary parameters from command flags

	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)

	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, eventType, eventSource, errorCode, errorMessage| filter eventSource = 'rds.amazonaws.com' | filter ispresent(responseElements) or ispresent(errorCode)| stats count(errorMessage) as errorCode by eventTime,errorMessage,eventName`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := comman_function.ProcessQueryResult(results)

	return processedResults, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxRDSErrorAnalysisCmd)
}
