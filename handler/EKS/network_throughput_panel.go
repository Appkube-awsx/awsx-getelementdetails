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
	PodNetworkRXBytes = "pod_network_rx_bytes"
	PodNetworkTXBytes = "pod_network_tx_bytes"
)

type NetworkThroughputResult struct {
	NetworkIn  []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NetworkIn"`
	NetworkOut []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NetworkOut"`
}

func GetNetworkThroughputPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime := parseTime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	networkInRawData, err := GetMetricData(clientAuth, clusterName, namespace, startTime, endTime, PodNetworkRXBytes)
	if err != nil {
		log.Println("Error fetching network in raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkIn"] = networkInRawData

	networkOutRawData, err := GetMetricData(clientAuth, clusterName, namespace, startTime, endTime, PodNetworkTXBytes)
	if err != nil {
		log.Println("Error fetching network out raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkOut"] = networkOutRawData

	result, _ := calculateNetworkThroughput(networkInRawData, networkOutRawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func parseTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		} else {
			startTime = &parsedStartTime
		}
	} else {
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
		now := time.Now()
		endTime = &now
	}

	return startTime, endTime
}

func GetMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
					Stat:   aws.String("Sum"), 
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

func calculateNetworkThroughput(networkInRawData, networkOutRawData *cloudwatch.GetMetricDataOutput) (NetworkThroughputResult, string) {
	var result NetworkThroughputResult

	result.NetworkIn = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(networkInRawData.MetricDataResults[0].Timestamps))
	for i, timestamp := range networkInRawData.MetricDataResults[0].Timestamps {
		result.NetworkIn[i].Timestamp = *timestamp
		result.NetworkIn[i].Value = *networkInRawData.MetricDataResults[0].Values[i]
	}

	result.NetworkOut = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(networkOutRawData.MetricDataResults[0].Timestamps))
	for i, timestamp := range networkOutRawData.MetricDataResults[0].Timestamps {
		result.NetworkOut[i].Timestamp = *timestamp
		result.NetworkOut[i].Value = *networkOutRawData.MetricDataResults[0].Values[i]
	}

	return result, ""
}
