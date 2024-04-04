package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	// "net/http"
	// "strconv"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"

	"github.com/aws/aws-sdk-go/aws"
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
			jsonResp, rawLogs, err := GetRdsErrorLogsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Printf("Error retrieving RDS error logs: %v\n", err)
				return
			}
			fmt.Println("JSON Response:")
			fmt.Println(jsonResp)
			fmt.Println("\nRaw Logs:")
			fmt.Println(rawLogs)
		}
	},
}

type RdsErrorLogEntry struct {
	Timestamp   string // Change type to string
	ErrorType   string // Change field name to ErrorType
	ErrorCode   int    // Store HTTP status code only
	Description string // No changes here
}

func GetRdsErrorLogsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (string, string, error) {
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
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", "", err
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
				return "", "", err
			}
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	results, err := filterCloudWatchLogsRDS(clientAuth, startTime, endTime, logGroupName, cloudWatchLogs)
	if err != nil {
		return "", "", err
	}

	processedResults := processQueryResults(results)
	jsonResp, err := json.Marshal(processedResults)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", "", err
	}

	// Concatenate raw logs
	rawLogs := ""
	for _, log := range processedResults {
		rawLogs += fmt.Sprintf("%s\t%s\t%s\t%d\n", log.Timestamp, log.ErrorType, log.Description, log.ErrorCode)
	}

	return string(jsonResp), rawLogs, nil
}

func filterCloudWatchLogsRDS(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		StartTime:    aws.Int64(startTime.Unix() * 1000),
		EndTime:      aws.Int64(endTime.Unix() * 1000),
		QueryString: aws.String(`fields @timestamp, @message, errorCode, eventType, errorMessage
| filter eventSource = 'rds.amazonaws.com' 
| filter ispresent(responseElements) or ispresent(errorCode)
| limit 1000`),
	}

	if cloudWatchLogs == nil {
		cloudWatchLogs = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	}

	queryResult, err := cloudWatchLogs.StartQuery(params)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %v", err)
	}

	queryId := queryResult.QueryId
	var queryResults []*cloudwatchlogs.GetQueryResultsOutput
	for {
		queryStatusInput := &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		}

		queryResult, err := cloudWatchLogs.GetQueryResults(queryStatusInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %v", err)
		}

		queryResults = append(queryResults, queryResult)

		if *queryResult.Status != "Complete" {
			log.Println("Query status is not complete.")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	return queryResults, nil
}

func processQueryResults(results []*cloudwatchlogs.GetQueryResultsOutput) []RdsErrorLogEntry {
	errorLogs := make([]RdsErrorLogEntry, 0)

	// Map error codes to HTTP status codes
	errorCodeMap := map[string]int{
		"DBInstanceNotFoundFault":              404,
		"AccessDenied":                         403,
		"InvalidParameterCombinationException":  400,
		"InvalidParameterValueException":       400,
		"InternalFailure":                      500, // Mapping for InternalFailure
		// Add more mappings as needed
	}

	for _, result := range results {
		if *result.Status == "Complete" {
			for _, res := range result.Results {
				entry := RdsErrorLogEntry{}

				for _, field := range res {
					if *field.Field == "@timestamp" {
						t, err := time.Parse("2006-01-02 15:04:05.000", *field.Value) // Adjust timestamp format
						if err != nil {
							log.Printf("Error parsing timestamp: %v", err)
						}
						entry.Timestamp = t.String()
					}
					if *field.Field == "errorCode" {
						entry.ErrorType = *field.Value
					}
					if *field.Field == "errorCode" {
						if code, ok := errorCodeMap[*field.Value]; ok {
							entry.ErrorCode = code
						} else {
							entry.ErrorCode = 500 // Default to Internal Server Error if code not found
						}
					}
					if *field.Field == "errorMessage" {
						entry.Description = *field.Value
					}
					// Add more field mappings as needed
				}

				errorLogs = append(errorLogs, entry)
			}
		}
	}

	return errorLogs
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
