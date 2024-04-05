package EC2

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var AwsxEc2ErrorRatePanelCmd = &cobra.Command{

	Use:   "error_rate_panel",
	Short: "Get error rate panel metrics data",
	Long:  `Command to get error rate panel metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running error rate panel command")

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
			results, err := GetInstanceErrorRatePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error in getting instance error rate panel: ", err)
				return
			}
			// processedResults := ProcessQueryResults(results)
			fmt.Println(results)
		}
	},
}

func GetInstanceErrorRatePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return nil, err
		}
		logGroupName = cmdbData.LogGroup
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				// handle error
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	events, err := filterCloudWatchlogs(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return nil, err
	}
	processedResults := ProcessQueryResults(events)

	return processedResults, nil
}

func filterCloudWatchlogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
    // Construct input parameters
    params := &cloudwatchlogs.StartQueryInput{
        LogGroupName: aws.String(logGroupName),
        StartTime:    aws.Int64(startTime.Unix() * 1000),
        EndTime:      aws.Int64(endTime.Unix() * 1000),
        QueryString: aws.String(`fields @timestamp, @message
            | filter eventSource=="ec2.amazonaws.com"
            | filter eventName=="RunInstances" and errorCode!=""
            | stats count(*) as ErrorCount by bin(1d)
            | sort @timestamp desc`),
    }

    if cloudWatchLogs == nil {
        cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
    }

    queryResult, err := cloudWatchLogs.StartQuery(params)
    if err != nil {
        return nil, fmt.Errorf("failed to start query: %v", err)
    }

    queryId := queryResult.QueryId
    // queryStatus := ""
   var queryResults []*cloudwatchlogs.GetQueryResultsOutput

	for {
		// Check query status
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryResults = append(queryResults, queryResult)

		if *queryResult.Status != "Complete" {
			time.Sleep(5 * time.Second) // wait before querying again
			continue
		}

		break // exit loop if query is complete
	}
	return queryResults, nil
}


func ProcessQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, resultField := range result.Results {
				for _, data := range resultField {
					if *data.Field == "eventName" {

						log.Printf("eventName: %s\n", *data)

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
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
}

// package EC2

// import (
//     "fmt"
//     "log"
//     "time"

//     "github.com/Appkube-awsx/awsx-common/authenticate"
//     "github.com/Appkube-awsx/awsx-common/awsclient"
//     "github.com/Appkube-awsx/awsx-common/cmdb"
//     "github.com/Appkube-awsx/awsx-common/config"
//     "github.com/Appkube-awsx/awsx-common/model"
//     "github.com/aws/aws-sdk-go/aws"
//     "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
//     "github.com/spf13/cobra"
// )

// var AwsxEc2ErrorRatePanelCmd = &cobra.Command{

//     Use:   "error_rate_panel",
//     Short: "get error rate panel metrics data",
//     Long:  `command to get error rate panel metrics data`,

//     Run: func(cmd *cobra.Command, args []string) {
//         fmt.Println("running from child command")

//         authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)

//         if err != nil {
//             log.Printf("Error during authentication: %v\n", err)
//             err := cmd.Help()
//             if err != nil {
//                 return
//             }
//             return
//         }

//         if authFlag {
//             _, err := GetInstanceErrorRatePanel(cmd, clientAuth, nil)
//             if err != nil {
//                 return
//             }
//         }
//     },
// }

// func GetInstanceErrorRatePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
//     elementId, _ := cmd.PersistentFlags().GetString("elementId")
//     cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
//     logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
//     if elementId != "" {
//         log.Println("getting cloud-element data from cmdb")
//         apiUrl := cmdbApiUrl
//         if cmdbApiUrl == "" {
//             log.Println("using default cmdb url")
//             apiUrl = config.CmdbUrl
//         }
//         log.Println("cmdb url: " + apiUrl)
//         cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
//         if err != nil {
//             return nil, err
//         }
//         logGroupName = cmdbData.LogGroup

//     }
//     startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
//     endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
//     var startTime, endTime *time.Time

//     // Parse start time if provided
//     if startTimeStr != "" {
//         parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
//         if err != nil {
//             log.Printf("Error parsing start time: %v", err)
//             err := cmd.Help()
//             if err != nil {
//                 // handle error
//             }
//         }
//         startTime = &parsedStartTime
//     } else {
//         defaultStartTime := time.Now().Add(-5 * time.Minute)
//         startTime = &defaultStartTime
//     }

//     if endTimeStr != "" {
//         parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
//         if err != nil {
//             log.Printf("Error parsing end time: %v", err)
//             err := cmd.Help()
//             if err != nil {
//                 // handle error
//             }
//         }
//         endTime = &parsedEndTime
//     } else {
//         defaultEndTime := time.Now()
//         endTime = &defaultEndTime
//     }

//     events, err := filterCloudWatchlogs(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
//     if err != nil {
//         log.Println("Error in getting sample count: ", err)
//         return nil, err
//     }
//     // for _, event := range events {
//     //     fmt.Println(event)
//     // }
//     processedResults := ProcessQueryResults(events)

//     return processedResults, nil
// }

// func filterCloudWatchlogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.ResultField, error) {
//     // Construct input parameters
//     params := &cloudwatchlogs.StartQueryInput{
//         LogGroupName: aws.String(logGroupName),
//         StartTime:    aws.Int64(startTime.Unix() * 1000),
//         EndTime:      aws.Int64(endTime.Unix() * 1000),
//         QueryString: aws.String(`fields @timestamp, @message
//             | filter eventSource=="ec2.amazonaws.com"
//             | filter eventName=="RunInstances" and errorCode!=""
//             | stats count(*) as ErrorCount by bin(1d)
//             | sort @timestamp desc`),
//     }

//     if cloudWatchLogs == nil {
//         cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
//     }

//     queryResult, err := cloudWatchLogs.StartQuery(params)
//     if err != nil {
//         return nil, fmt.Errorf("failed to start query: %v", err)
//     }

//     queryId := queryResult.QueryId
//     queryStatus := ""
//     var queryResults *cloudwatchlogs.GetQueryResultsOutput // Declare queryResults outside the loop
//     for queryStatus != "Complete" {
//         // Check query status
//         queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
//             QueryId: queryId,
//         }

//         queryResults, err = cloudWatchLogs.GetQueryResults(queryStatusInput) // Assign value to queryResults
//         if err != nil {
//             return nil, fmt.Errorf("failed to get query results: %v", err)
//         }

//         queryStatus = aws.StringValue(queryResults.Status)
//         time.Sleep(1 * time.Second) // Wait for a second before checking status again
//     }

//     // Query is complete, now process results
//     var results []*cloudwatchlogs.ResultField
//     for _, resultRow := range queryResults.Results {
//         for _, resultField := range resultRow {
//             results = append(results, resultField)
//         }
//     }

//     return results, nil
// }
// func ProcessQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []*cloudwatchlogs.GetQueryResultsOutput {
// 	processedResults := make([]*cloudwatchlogs.GetQueryResultsOutput, 0)

// 	for _, result := range results {
// 		if *result.Status == "Complete" {
// 			for _, resultField := range result.Results {
// 				for _, data := range resultField {
// 					if *data.Field == "eventName" {

// 						log.Printf("eventName: %s\n", *data)

// 					}
// 				}
// 			}
// 			processedResults = append(processedResults, result)
// 		} else {
// 			log.Println("Query status is not complete.")
// 		}
// 	}

// 	return processedResults
// }

// func init() {
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("elementId", "", "element id")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("vaultToken", "", "vault token")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("accountId", "", "aws account number")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("zone", "", "aws region")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("accessKey", "", "aws access key")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("secretKey", "", "aws secret key")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("externalId", "", "aws external id")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("ServiceName", "", "Service Name")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("elementType", "", "element type")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("instanceId", "", "instance id")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("clusterName", "", "cluster name")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("query", "", "query")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("startTime", "", "start time")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("endTime", "", "endcl time")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
// 	AwsxEc2ErrorRatePanelCmd.PersistentFlags().String("logGroupName", "", "log group name")
// }
