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

const (
	PodNetworkRXByte = "pod_network_rx_bytes"
	PodNetworkTXByte = "pod_network_tx_bytes"
)

type NetworKThroughputResult struct {
	Throughput []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Throughput"`
}

func GetNetworkThroughputSinglePanel(cmd *cobra.Command, clientAuth *model.Auth) (*cloudwatch.GetMetricDataOutput, string, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime := ParseTime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Fetch network in raw data
	networkInRawData, err := GetmetricData(clientAuth, clusterName, namespace, startTime, endTime, PodNetworkRXByte)
	if err != nil {
		log.Println("Error fetching network in raw data: ", err)
		return nil, "", err
	}

	// Fetch network out raw data
	networkOutRawData, err := GetmetricData(clientAuth, clusterName, namespace, startTime, endTime, PodNetworkTXByte)
	if err != nil {
		log.Println("Error fetching network out raw data: ", err)
		return nil, "", err
	}

	// Calculate network throughput
	result := calculateNetworKThroughput(networkInRawData, networkOutRawData)

	// Marshal result to JSON string
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return nil, "", err
	}

	return networkInRawData, string(jsonString), nil
}

// Function to parse time strings and return time pointers
func ParseTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
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

// Function to fetch CloudWatch metric data
func GetmetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String(metricName),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"), // Using Sum as an example, change as needed
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

// Function to calculate network throughput
func calculateNetworKThroughput(networkInRawData, networkOutRawData *cloudwatch.GetMetricDataOutput) NetworKThroughputResult {
	var result NetworKThroughputResult

	result.Throughput = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(networkInRawData.MetricDataResults[0].Timestamps))

	for i, timestamp := range networkInRawData.MetricDataResults[0].Timestamps {
		// Calculate network throughput (difference between network in and out)
		throughput := *networkInRawData.MetricDataResults[0].Values[i] - *networkOutRawData.MetricDataResults[0].Values[i]
		result.Throughput[i].Timestamp = *timestamp
		result.Throughput[i].Value = throughput
	}

	return result
}
