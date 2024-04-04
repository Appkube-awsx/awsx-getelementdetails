package RDS

import (
    "fmt"
    "log"
    "time"

    "github.com/Appkube-awsx/awsx-common/authenticate"
    "github.com/Appkube-awsx/awsx-common/awsclient"
    "github.com/Appkube-awsx/awsx-common/config"
    "github.com/Appkube-awsx/awsx-common/model"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
    "github.com/spf13/cobra"
)

// Define a struct to hold the log entry data
type ErrorAnalysisEntry struct {
    EventTime    time.Time
    EventType    string
    EventSource  string
    ErrorCode    string
    ErrorMessage string
}

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
                log.Printf("Error retrieving error analysis panel data: %v\n", err)
                return
            }
            // Print the data
            for _, entry := range data {
                fmt.Printf("%+v\n", entry)
            }
        }
    },
}

// Function to retrieve error analysis panel data
func GetErrorAnalysisData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]ErrorAnalysisEntry, error) {
    // Retrieve necessary parameters from command flags
    elementId, _ := cmd.PersistentFlags().GetString("elementId")
    cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
    logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
    startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
    endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

    // Fetching cloud-element data from cmdb if necessary
    if elementId != "" {
        log.Println("Getting cloud-element data from cmdb")
        apiUrl := cmdbApiUrl
        if cmdbApiUrl == "" {
            log.Println("Using default cmdb url")
            apiUrl = config.CmdbUrl
        }
        log.Println("CMDB url:", apiUrl)
    }

    // Parsing start and end times for querying logs
    var startTime, endTime time.Time
    if startTimeStr != "" {
        parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
        if err != nil {
            return nil, fmt.Errorf("Error parsing start time: %v", err)
        }
        startTime = parsedStartTime
    } else {
        startTime = time.Now().Add(-5 * time.Minute)
    }
    if endTimeStr != "" {
        parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
        if err != nil {
            return nil, fmt.Errorf("Error parsing end time: %v", err)
        }
        endTime = parsedEndTime
    } else {
        endTime = time.Now()
    }

    // Fetching query results from CloudWatch Logs
    results, err := filterCloudWatchLogs(clientAuth, &startTime, &endTime, logGroupName, cloudWatchLogs)
    if err != nil {
        return nil, err
    }

    // Process query results into ErrorAnalysisEntry structs
    data := processQueryResultss(results)
    return data, nil
}

// Function to filter CloudWatch Logs and retrieve query results
func filterCloudWatchLogs(clientAuth *model.Auth, startTime, endTime *time.Time, logGroupName string, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
    params := &cloudwatchlogs.StartQueryInput{
        LogGroupName: aws.String(logGroupName),
        StartTime:    aws.Int64(startTime.Unix() * 1000),
        EndTime:      aws.Int64(endTime.Unix() * 1000),
        QueryString: aws.String(`fields @timestamp, eventType, eventSource, errorCode, errorMessage
            | filter eventSource = 'rds.amazonaws.com' 
            | filter ispresent(responseElements) or ispresent(errorCode)
            | stats count(errorMessage) as errorCode by eventTime,errorMessage,eventName`),
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
            time.Sleep(1 * time.Second)
            continue
        }

        break
    }

    return queryResults, nil
}

// Function to process query results into ErrorAnalysisEntry structs
func processQueryResultss(results []*cloudwatchlogs.GetQueryResultsOutput) []ErrorAnalysisEntry {
    var data []ErrorAnalysisEntry

    for _, result := range results {
        if *result.Status == "Complete" {
            for _, res := range result.Results {
                var entry ErrorAnalysisEntry
                for _, field := range res {
                    switch *field.Field {
                    case "@timestamp":
                        t, err := time.Parse("2006-01-02 15:04:05.000", *field.Value)
                        if err != nil {
                            log.Printf("Error parsing timestamp: %v", err)
                        }
                        entry.EventTime = t
                    case "eventType":
                        entry.EventType = *field.Value
                    case "eventSource":
                        entry.EventSource = *field.Value
                    case "errorCode":
                        entry.ErrorCode = *field.Value
                    case "errorMessage":
                        entry.ErrorMessage = *field.Value
                    }
                }
                data = append(data, entry)
            }
        } else {
            log.Println("Query status is not complete.")
        }
    }

    return data
}

func init() {
    AwsxRDSErrorAnalysisCmd.PersistentFlags().String("elementId", "", "Element ID")
    AwsxRDSErrorAnalysisCmd.PersistentFlags().String("cmdbApiUrl", "", "CMDB API URL")
    AwsxRDSErrorAnalysisCmd.PersistentFlags().String("logGroupName", "", "Log Group Name")
    AwsxRDSErrorAnalysisCmd.PersistentFlags().String("startTime", "", "Start Time")
    AwsxRDSErrorAnalysisCmd.PersistentFlags().String("endTime", "", "End Time")
}
