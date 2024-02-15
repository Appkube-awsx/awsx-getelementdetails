package EKS

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

// NodeEventLog represents a node event log entry.
type NodeEventLog struct {
	Timestamp       int64  `json:"Timestamp"`
	EventType       string `json:"EventType"`
	SourceComponent string `json:"SourceComponent"`
	EventMessage    string `json:"EventMessage"`
}

func GetNodeEventLogsSinglePanel(cmd *cobra.Command, clientAuth *model.Auth) (string, string, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime := Parsetime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Fetch node event logs
	nodeEventLogs, err := GetNodeEventLogs(clientAuth, clusterName, startTime, endTime)
	if err != nil {
		log.Println("Error fetching node event logs: ", err)
		return "", "", err
	}

	// Marshal node event logs to JSON string
	jsonString, err := json.Marshal(nodeEventLogs)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", "", err
	}

	// Concatenate raw logs
	rawLogs := ""
	for _, event := range nodeEventLogs {
		rawLogs += event.EventMessage + "\n"
	}

	return string(jsonString), rawLogs, nil
}

// Function to fetch node event logs
func GetNodeEventLogs(clientAuth *model.Auth, clusterName string, startTime, endTime *time.Time) ([]NodeEventLog, error) {
	logGroupName := "/aws/containerinsights/" + clusterName + "/host"
	startTimeMillis := startTime.UnixNano() / int64(time.Millisecond) // Convert to milliseconds
	endTimeMillis := endTime.UnixNano() / int64(time.Millisecond)     // Convert to milliseconds
	filterPattern := "\"node\""
	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  &logGroupName,
		StartTime:     aws.Int64(startTimeMillis), // Pass the calculated value directly
		EndTime:       aws.Int64(endTimeMillis),   // Pass the calculated value directly
		FilterPattern: &filterPattern,             // Pass the address of the filter pattern string
	}
	// Get CloudWatchLogs client from awsclient
	cloudWatchLogsClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	result, err := cloudWatchLogsClient.FilterLogEvents(input)
	if err != nil {
		return nil, err
	}

	var nodeEventLogs []NodeEventLog
	for _, event := range result.Events {
		// Parse the ident field from the message
		var logData map[string]interface{}
		err := json.Unmarshal([]byte(*event.Message), &logData)
		if err != nil {
			log.Printf("Error parsing event log message: %v", err)
			continue
		}

		ident, ok := logData["ident"].(string)
		if !ok {
			log.Println("Error parsing ident field from event log message")
			continue
		}

		nodeEventLogs = append(nodeEventLogs, NodeEventLog{
			Timestamp:       *event.Timestamp / 1000, // Convert to milliseconds
			EventType:       "NodeEvent",
			SourceComponent: ident, // Set SourceComponent to ident
			EventMessage:    *event.Message,
		})
	}

	return nodeEventLogs, nil
}

// Function to parse time strings and return time pointers
func Parsetime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		} else {
			startTime = &parsedStartTime
		}
	} else {
		// If startTimeStr is empty, default to the last five minutes
		now := time.Now()
		startTime = &now
		minusFiveMinutes := now.Add(-5 * time.Minute)
		startTime = &minusFiveMinutes
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
		} else {
			endTime = &parsedEndTime
		}
	} else {
		// If endTimeStr is empty, default to the current time
		now := time.Now()
		endTime = &now
	}

	return startTime, endTime
}
