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

var AwsxRdsErrorLogsCmd = &cobra.Command{
	Use:   "rds_error_logs",
	Short: "Get recent error logs for RDS instances",
	Long:  `Command to retrieve recent error logs for RDS instances`,
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
			data, err := GetRdsErrorLogsPanel(cmd, clientAuth, nil)
			if err != nil {

				return
			}
			fmt.Println(data)
		}
	},
}

// type RdsErrorLogEntry struct {
// 	Timestamp   string // Change type to string
// 	ErrorType   string // Change field name to ErrorType
// 	ErrorCode   int    // Store HTTP status code only
// 	Description string // No changes here
// }

func GetRdsErrorLogsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {

	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)

	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message, errorCode, eventType, errorMessage| filter eventSource = 'rds.amazonaws.com' | filter ispresent(responseElements) or ispresent(errorCode)| limit 1000`, cloudWatchLogs)

	if err != nil {
		return nil, nil
	}
	processedResults := comman_function.ProcessQueryResult(results)

	return processedResults, nil

}

// func convertToHTTPStatusCode(errorCode int) int {
//     // Map error codes to corresponding HTTP status codes
//     switch errorCode {
//     case 400:
//         return http.StatusBadRequest
//     case 401:
//         return http.StatusUnauthorized
//     case 403:
//         return http.StatusForbidden
//     case 404:
//         return http.StatusNotFound
//     case 500:
//         return http.StatusInternalServerError
//     // Add more mappings as needed
//     default:
//         return errorCode // Return the error code as is if no mapping is found
//     }
// }

func init() {
	AwsxRdsErrorLogsCmd.PersistentFlags().String("elementId", "", "Element ID")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API URL")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("logGroupName", "", "Log Group Name")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("startTime", "", "Start Time")
	AwsxRdsErrorLogsCmd.PersistentFlags().String("endTime", "", "End Time")
}
