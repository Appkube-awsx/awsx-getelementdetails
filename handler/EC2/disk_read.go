package EC2

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type DiskReadPanelData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

func GetDiskReadPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID, _ := cmd.PersistentFlags().GetString("instanceID")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
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
			log.Printf("Error parsing end time: %v", err)
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := GetDiskReadPanelMetricData(clientAuth, instanceID, namespace, startTime, endTime)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RawData"] = rawData

	result := processDiskReadPanelRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetDiskReadPanelMetricData(clientAuth *model.Auth, instanceID string, namespace string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, namespace, startTime, endTime)
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceID),
							},
						},
						MetricName: aws.String("DiskReadBytes"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
			},
		},
	}
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		log.Println("Error fetching metric data:", err)
		return nil, err
	}
	log.Println("Metric data result:", result)

	return result, nil
}

func processDiskReadPanelRawData(result *cloudwatch.GetMetricDataOutput) DiskReadPanelData {
	var rawData DiskReadPanelData

	// Initialize an empty slice to store the raw data
	rawData.RawData = []struct {
		Timestamp time.Time
		Value     float64
	}{}

	// Iterate over each metric data result
	for _, metricDataResult := range result.MetricDataResults {
		// Iterate over each timestamp and value pair in the current metric data result
		for i, timestamp := range metricDataResult.Timestamps {
			// Append the timestamp and value to the rawData slice
			rawData.RawData = append(rawData.RawData, struct {
				Timestamp time.Time
				Value     float64
			}{
				Timestamp: *timestamp,
				Value:     *metricDataResult.Values[i],
			})
		}
	}

	return rawData
}
