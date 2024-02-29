package EC2

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/spf13/cobra"
)

var AwsxEc2StartCountCmd = &cobra.Command{
	Use:   "instance_start_count",
	Short: "Get instance start count metric data",
	Long:  `Command to get instance start count metric data`,

	Run: func(cmd *cobra.Command, args []string) {
		var clientAuth *model.Auth
		var jsonString bool
		var err error

		jsonString, clientAuth, err = authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		fmt.Println(jsonString)
		startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
		endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

		var startTime, endTime time.Time

		if startTimeStr != "" {
			startTime, err = time.Parse(time.RFC3339, startTimeStr)
			if err != nil {
				log.Printf("Error parsing start time: %v", err)
				return
			}
		} else {
			startTime = time.Now().AddDate(0, -12, 0) // 12 months ago
		}

		if endTimeStr != "" {
			endTime, err = time.Parse(time.RFC3339, endTimeStr)
			if err != nil {
				log.Printf("Error parsing end time: %v", err)
				return
			}
		} else {
			endTime = time.Now() // Current time
		}

		err = GetInstanceStartCountMetricData(clientAuth, startTime, endTime)
		if err != nil {
			log.Println("Error getting instance start count metric data: ", err)
			return
		}
	},
}

func GetInstanceStartCountMetricData(clientAuth *model.Auth, startTime, endTime time.Time) error {
	cwLogs := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)

	// Start the query
	queryInput := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String("CloudTrail/DefaultLogGroup"), // Specify the log group name
		StartTime:    aws.Int64(startTime.Unix()),
		EndTime:      aws.Int64(endTime.Unix()),

		QueryString: aws.String(`fields @timestamp, @message
            | filter eventSource=="ec2.amazonaws.com"
            | filter eventName=="StartInstances"
            | stats count(*) as InstanceCount by bin(1mo)
            | sort @timestamp desc`), // Your CloudWatch Logs Insights query
	}

	queryResult, err := cwLogs.StartQuery(queryInput)
	if err != nil {
		return fmt.Errorf("failed to start query: %v", err)
	}

	queryId := queryResult.QueryId
	queryStatus := ""
	var queryResults *cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop
	for queryStatus != "Complete" {
		// Check query status
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResults, err = cwLogs.GetQueryResults(queryStatusInput) // Assign value to queryResults
		if err != nil {
			return fmt.Errorf("failed to get query results: %v", err)
		}

		queryStatus = aws.StringValue(queryResults.Status)
		time.Sleep(1 * time.Second) // Wait for a second before checking status again
	}

	// Query is complete, now process results
	for _, resultRow := range queryResults.Results {
		for _, resultField := range resultRow {
			fmt.Println(*resultField.Value)
		}
	}

	return nil
}

func init() {
	AwsxEc2StartCountCmd.PersistentFlags().String("startTime", "", "start time (RFC3339 format)")
	AwsxEc2StartCountCmd.PersistentFlags().String("endTime", "", "end time (RFC3339 format)")
}

// package EC2

// import (
//     "fmt"
//     "github.com/Appkube-awsx/awsx-common/awsclient"
//     "github.com/aws/aws-sdk-go/aws"
//     "github.com/aws/aws-sdk-go/aws/awserr"
//     "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
//     "time"

//     "github.com/Appkube-awsx/awsx-common/model"
// )

// func GetInstanceStartCountMetricData(clientAuth *model.Auth) error {
//     cwLogs := awsclient.GetClient(*clientAuth, awsclient.CLOUD_WATCH_LOGS).(*cloudwatchlogs.CloudWatchLogs)

//     // Start the query
//     queryInput := &cloudwatchlogs.StartQueryInput{
//         LogGroupName: aws.String("CloudTrail/DefaultLogGroup"), // Specify the log group name
//         StartTime:    aws.Int64(time.Now().Add(-12 * 30 * 24 * time.Hour).Unix()), // Example start time (12 months ago)
//         EndTime:      aws.Int64(time.Now().Unix()),                   // Example end time (current time)
//         QueryString:  aws.String(`fields @timestamp, @message
//             | filter eventSource=="ec2.amazonaws.com"
//             | filter eventName=="StartInstances"
//             | stats count(*) as InstanceCount by bin(1mo)
//             | sort @timestamp desc`), // Your CloudWatch Logs Insights query
//     }

//     queryResult, err := cwLogs.StartQuery(queryInput)
//     if err != nil {
//         return fmt.Errorf("failed to start query: %v", err)
//     }

//     queryId := queryResult.QueryId
//     queryStatus := ""
//     for queryStatus != "Complete" {
//         // Check query status
//         queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
//             QueryId: queryId,
//         }

//         queryResults, err := cwLogs.GetQueryResults(queryStatusInput)
//         if err != nil {
//             return fmt.Errorf("failed to get query results: %v", err)
//         }

//         queryStatus = aws.StringValue(queryResults.Status)
//         time.Sleep(1 * time.Second) // Wait for a second before checking status again
//     }

//     // Query is complete, now process results
//     for _, resultRow := range queryResults.Results {
//         for _, resultField := range resultRow {
//             fmt.Println(*resultField.Value)
//         }
//     }

//     return nil
// }
