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

type ServiceAvailabilityResult struct {
	Availability float64 `json:"availability"`
}

type TimeseriesDataPoint struct {
	Timestamp     time.Time `json:"timestamp"`
	Availability  float64   `json:"availability"`
}

func GetServiceAvailabilityData(cmd *cobra.Command, clientAuth *model.Auth) (string, []TimeseriesDataPoint, error) {
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

	rawData, err := GetServiceAvailabilityMetricData(clientAuth, clusterName, namespace, startTime, endTime)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	var timeSeriesData []TimeseriesDataPoint
	for i, timestamp := range rawData.MetricDataResults[0].Timestamps {
		availability := ProcessServiceAvailabilityRawData(rawData, i)
		dataPoint := TimeseriesDataPoint{
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

func GetServiceAvailabilityMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("pod_status_running"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("SampleCount"),
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
						MetricName: aws.String("pod_status_pending"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("SampleCount"),
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
						MetricName: aws.String("pod_status_ready"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("SampleCount"),
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

func ProcessServiceAvailabilityRawData(result *cloudwatch.GetMetricDataOutput, index int) float64 {
	// Calculate service availability based on the metrics
	totalRunning := float64(0)
	totalPending := float64(0)
	totalReady := float64(0)

	for _, result := range result.MetricDataResults {
		if *result.Id == "m1" {
			for _, value := range result.Values {
				totalRunning += *value
			}
		} else if *result.Id == "m2" {
			for _, value := range result.Values {
				totalPending += *value
			}
		} else if *result.Id == "m3" {
			for _, value := range result.Values {
				totalReady += *value
			}
		}
	}

	// Calculate service availability
	totalPods := totalRunning + totalPending + totalReady
	if totalPods > 0 {
		return (totalReady / totalPods) * 100
	} else {
		return 0
	}
}
