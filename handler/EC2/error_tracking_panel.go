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

var AwsxEc2ErrorTrackingPanelCmd = &cobra.Command{

	Use:   "error_tracking_panel",
	Short: "Get error tracking panel metrics data",
	Long:  `Command to get error tracking panel metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running error tracking panel command")

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
			results, err := GetErrorTrackingPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error in getting instance error rate panel: ", err)
				return
			}
			// processedResults := ProcessQueryResults(results)
			fmt.Println(results)
		}
	},
}

func GetErrorTrackingPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	events, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource=="ec2.amazonaws.com"| filter  eventName=="RunInstances"  and errorCode!=""| filter ispresent(responseElements) or ispresent(errorCode)| stats count(*) as errorCount by eventTime,eventName,errorCode,errorMessage`, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return nil, err
	}
	processedResults := ProcessQueryResultz(events)

	return processedResults, nil
}
func ProcessQueryResultz(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventName" {

						log.Printf("eventName: %s\n", *data)

					} else if *data.Field == "eventTime" {

						log.Printf("eventTime: %s\n", *data)

					} else if *data.Field == "errorCode" {

						log.Printf("errorCode: %s\n", *data)

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
	comman_function.InitAwsCmdFlags(AwsxEc2ErrorTrackingPanelCmd)
}