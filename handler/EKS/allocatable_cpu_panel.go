package EKS

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

type TimeSeriesData struct {
	Timestamp     time.Time
	AllocatableCPU float64
}

type AllocateResult struct {
	RawData []TimeSeriesData `json:"RawData"`
}

func GetAllocatableCPUData(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
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

	// Debug prints
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := GetAllocatableCPUMetricData(clientAuth, clusterName, namespace, startTime, endTime)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RawData"] = rawData

	// Process the raw data if needed
	result := processCPURawData(rawData)

	// Log only the allocatable CPU and its corresponding timestamp
	for _, data := range result.RawData {
		log.Printf("Timestamp: %v, Allocatable CPU: %v", data.Timestamp, data.AllocatableCPU)
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetAllocatableCPUMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
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
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("node_cpu_limit"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60), // Set a common period for both metrics
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("m2"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("node_cpu_reserved_capacity"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60), // Set a common period for both metrics
					Stat:   aws.String("Average"),
				},
			},
		},
	}
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processCPURawData(result *cloudwatch.GetMetricDataOutput) AllocateResult {
	var rawData AllocateResult
	rawData.RawData = make([]TimeSeriesData, len(result.MetricDataResults[0].Timestamps))

	// Assuming the two metrics have the same number of data points
	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		cpuLimit := *result.MetricDataResults[0].Values[i]
		reservedCapacity := *result.MetricDataResults[1].Values[i]

		// Log the values for troubleshooting
		// log.Printf("Timestamp: %v, cpuLimit: %v, reservedCapacity: %v", *timestamp, cpuLimit, reservedCapacity)

		allocatableCPU := cpuLimit - reservedCapacity

		// Only include the calculated allocatable CPU in the result
		rawData.RawData[i].AllocatableCPU = allocatableCPU
	}

	return rawData
}
