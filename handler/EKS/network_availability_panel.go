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

type NetworkAvailabilityResult struct {
	Availability float64 `json:"availability"`
}

type TimeSeriesDataPoint struct {
	Timestamp     time.Time `json:"timestamp"`
	Availability  float64   `json:"availability"`
}

func GetNetworkAvailabilityData(cmd *cobra.Command, clientAuth *model.Auth) (string, []TimeSeriesDataPoint, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
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

	rawData, err := GetNetworkAvailabilityMetricData(clientAuth, clusterName, namespace, startTime, endTime)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	var timeSeriesData []TimeSeriesDataPoint
	for i, timestamp := range rawData.MetricDataResults[0].Timestamps {
		availability := ProcessNetworkAvailabilityRawData(rawData, i)
		dataPoint := TimeSeriesDataPoint{
			Timestamp:    *timestamp,
			Availability: availability,
		}
		timeSeriesData = append(timeSeriesData, dataPoint)
	}

	jsonString, err := json.Marshal(timeSeriesData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), timeSeriesData, nil
}

func GetNetworkAvailabilityMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_interface_network_tx_dropped"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
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
						MetricName: aws.String("node_interface_network_rx_dropped"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("m3"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("pod_network_rx_bytes"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("m4"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("pod_network_tx_bytes"),
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

func ProcessNetworkAvailabilityRawData(result *cloudwatch.GetMetricDataOutput, index int) float64 {
	// Calculate network availability based on the metrics
	totalTxDropped := float64(0)
	totalRxDropped := float64(0)
	totalRxBytes := float64(0)
	totalTxBytes := float64(0)

	for _, result := range result.MetricDataResults {
		if *result.Id == "m1" {
			for _, value := range result.Values {
				totalTxDropped += *value
			}
		} else if *result.Id == "m2" {
			for _, value := range result.Values {
				totalRxDropped += *value
			}
		} else if *result.Id == "m3" {
			for _, value := range result.Values {
				totalRxBytes += *value
			}
		} else if *result.Id == "m4" {
			for _, value := range result.Values {
				totalTxBytes += *value
			}
		}
	}

	// Calculate network availability
	if totalTxBytes > 0 && totalTxBytes > totalTxDropped && totalRxBytes > 0 && totalRxBytes > totalRxDropped {
		return 100 * (1 - ((totalTxDropped + totalRxDropped) / (totalTxBytes + totalRxBytes)))
	} else {
		return 0
	}
}
