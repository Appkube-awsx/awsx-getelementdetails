package EKS

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
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

var AwsxEKSNodeEventLogsCmd = &cobra.Command{
	Use:   "node_event_logs_panel",
	Short: "get node event logs data",
	Long:  `command to get node event logs data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")
		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err := GetNodeEventLogsSinglePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Node event logs data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNodeEventLogsSinglePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", "", fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", "", fmt.Errorf("error getting instance ID: %v", err)
	}

	// Fetch node event logs
	nodeEventLogs, err := GetNodeEventLogs(clientAuth, instanceId, startTime, endTime, cloudWatchClient)
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

func GetNodeEventLogs(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) ([]NodeEventLog, error) {
	logGroupName := "/aws/containerinsights/" + instanceId + "/host"
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

func init() {
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeEventLogsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
