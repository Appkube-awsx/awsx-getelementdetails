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

type StorageUtilizationResult struct {
	RootVolumeUsage float64 `json:"rootVolumeUsage"`
	EBSVolume1Usage float64 `json:"ebsVolume1Usage"`
	EBSVolume2Usage float64 `json:"ebsVolume2Usage"`
}

func GetStorageUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
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
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Root Volume Usage
	rootVolumeUsage, err := GetStorageMetricData(clientAuth, clusterName, namespace, startTime, endTime, "node_filesystem_utilization")
	if err != nil {
		log.Println("Error in getting root volume usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolumeUsage"] = rootVolumeUsage

	// Get EBS Volume 1 Usage
	ebsVolume1Usage, err := GetStorageMetricData(clientAuth, clusterName, namespace, startTime, endTime, "node_filesystem_utilization")
	if err != nil {
		log.Println("Error in getting EBS volume 1 usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBSVolume1Usage"] = ebsVolume1Usage

	// Get EBS Volume 2 Usage
	ebsVolume2Usage, err := GetStorageMetricData(clientAuth, clusterName, namespace, startTime, endTime, "node_filesystem_utilization")
	if err != nil {
		log.Println("Error in getting EBS volume 2 usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBSVolume2Usage"] = ebsVolume2Usage

	// Create JSON output
	jsonOutput := StorageUtilizationResult{
		RootVolumeUsage: *rootVolumeUsage.MetricDataResults[0].Values[0],
		EBSVolume1Usage: *ebsVolume1Usage.MetricDataResults[0].Values[0],
		EBSVolume2Usage: *ebsVolume2Usage.MetricDataResults[0].Values[0],
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetStorageMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
							// Add dimensions for specific EBS volumes if needed
						},
						MetricName: aws.String(metricName),
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
		return nil, err
	}

	return result, nil
}
