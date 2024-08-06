package S3

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxMaximumErrorsMessageCmd = &cobra.Command{

	Use:   "maximum_errors_message_panel",
	Short: "Get maximum errors message data",
	Long:  `Command to get maximum errors message data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running maximum errors message panel command")

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
			panel, err := GetMaximumErrorsMessageData(cmd, clientAuth, nil)
			if err != nil {
				return
			}
			fmt.Println(panel)
		}
	},
}

func GetMaximumErrorsMessageData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = comman_function.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := comman_function.GetLogsData(clientAuth, startTime, endTime, logGroupName, `fields @timestamp, @message| filter eventSource == "s3.amazonaws.com"| filter eventName == "GetBucketAcl"| filter ispresent(errorCode) and errorCode != ""| filter errorCode in ["AccessDenied", "NoSuchBucket", "NoSuchKey", "InvalidBucketName", "AllAccessDisabled", "InvalidObjectState", "RequestTimeTooSkewed"]| stats count(errorCode) as ErrorCodeCount,count(errorMessage) as ErrorCount by requestParameters.bucketName as BucketName,errorCode| sort ErrorCodeCount desc| limit 10`, cloudWatchLogs)
	if err != nil {
		return nil, nil
	}
	processedResults := ProcessQueryResult(results)

	return processedResults, nil

}

func ProcessQueryResult(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "requestParameters.bucketName" {

						log.Printf("requestParameters.bucketName: %s\n", *data)

					} else if *data.Field == "errorCode" {

						log.Printf("errorCode: %s\n", *data)

					} else if *data.Field == "ErrorCodeCount" {

						log.Printf("ErrorCodeCount: %s\n", *data)

					} else if *data.Field == "ErrorCount" {

						log.Printf("ErrorCount: %s\n", *data)
					
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
	comman_function.InitAwsCmdFlags(AwsxMaximumErrorsMessageCmd)
}
