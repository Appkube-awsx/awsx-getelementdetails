package Lambda

import (
	"encoding/json"
	// "errors"
	"fmt"
	"log"
	// "strings"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type ConcurrencyData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxLambdaConcurrencyCmd = &cobra.Command{
	Use:   "concurrency_panel",
	Short: "get lambda concurrency metrics data",
	Long:  `Command to get lambda concurrency metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running from child command")
		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			return
		}
		if !authFlag {
			log.Println("Authentication failed")
			return
		}

		responseType, _ := cmd.PersistentFlags().GetString("responseType")
		if responseType != "json" && responseType != "frame" {
			log.Println("Invalid response type. Valid options are 'json' or 'frame'.")
			return
		}

		jsonResp, cloudwatchMetricResp, err := GetLambdaConcurrencyData(cmd, clientAuth, nil)
		if err != nil {
			log.Println("Error getting lambda concurrency data: ", err)
			return
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaConcurrencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, *map[string]*cloudwatch.GetMetricDataOutput, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := GetLambdaConcurrencyMetricData(clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	cloudwatchMetricData["concurrency"] = rawData

	result := processConcurrencyRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		return "", nil, err
	}

	return string(jsonString), &cloudwatchMetricData, nil
}

func GetLambdaConcurrencyMetricData(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("ConcurrentExecutions"),
						Namespace:  aws.String("AWS/Lambda"),
					},
					Period: aws.Int64(300), // 5 minutes
					Stat:   aws.String("Average"),
				},
			},
		},
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processConcurrencyRawData(result *cloudwatch.GetMetricDataOutput) *ConcurrencyData {
	concurrencyData := &ConcurrencyData{
		RawData: make([]struct {
			Timestamp time.Time
			Value     float64
		}, len(result.MetricDataResults[0].Timestamps)),
	}

	for i := range result.MetricDataResults[0].Timestamps {
		concurrencyData.RawData[i].Timestamp = *result.MetricDataResults[0].Timestamps[i]
		concurrencyData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return concurrencyData
}

func init() {
	AwsxLambdaConcurrencyCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaConcurrencyCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaConcurrencyCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
