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

type NodeDowntimeDataPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	NodeDowntime float64   `json:"nodeDowntime"`
}

func GetNodeDowntimePanel(cmd *cobra.Command, clientAuth *model.Auth) (string, []NodeDowntimeDataPoint, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
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

	// Get node metrics
	nodeMetrics, err := GetNodeDowntimeMetrics(clientAuth, clusterName, startTime, endTime)
	if err != nil {
		log.Println("Error in getting node metrics: ", err)
		return "", nil, err
	}

	// Calculate node downtime data points
	var downtimeData []NodeDowntimeDataPoint
	for i := 0; i < len(nodeMetrics.MetricDataResults[0].Values); i++ {
		downtime := 0.0
		if *nodeMetrics.MetricDataResults[0].Values[i] <= 0 {
			downtime = 1.0
		}
		dataPoint := NodeDowntimeDataPoint{
			Timestamp:   *nodeMetrics.MetricDataResults[0].Timestamps[i],
			NodeDowntime: downtime,
		}
		downtimeData = append(downtimeData, dataPoint)
	}

	jsonString, err := json.Marshal(downtimeData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), downtimeData, nil
}

func GetNodeDowntimeMetrics(clientAuth *model.Auth, clusterName string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cpu_utilization"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("node_cpu_utilization"),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(60),

					Stat: aws.String("Average"),
				},
			},
			{
				Id: aws.String("memory_utilization"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("node_memory_utilization"),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(60),
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
