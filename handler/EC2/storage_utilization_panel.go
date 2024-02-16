package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type volumeUsage struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type volumeMetrics struct {
	RootVolumeUsage volumeUsage `json:"RootVolumeUsage"`
	EBSVolume1Usage volumeUsage `json:"EBSVolume1Usage"`
	EBSVolume2Usage volumeUsage `json:"EBSVolume2Usage"`
}

func GetVolumeMetricsPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	RootVolumeId, _ := cmd.PersistentFlags().GetString("RootVolumeId")
	EBSVolume1Id, _ := cmd.PersistentFlags().GetString("EBSVolume1Id")
	EBSVolume2Id, _ := cmd.PersistentFlags().GetString("EBSVolume2Id")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, err
		}
		instanceId = cmdbData.InstanceId

	}

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

	// Get metrics for root volume
	rootVolumeMetrics, err := GetMetrics(clientAuth, instanceId, RootVolumeId, elementType, startTime, endTime, "VolumeStalledIOCheck")
	if err != nil {
		log.Println("Error in getting metrics for root volume: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolume"] = rootVolumeMetrics

	// Get metrics for EBS1 volume
	ebsVolume1Metrics, err := GetMetrics(clientAuth, instanceId, EBSVolume1Id, elementType, startTime, endTime, "VolumeStalledIOCheck")
	if err != nil {
		log.Println("Error in getting metrics for EBS1 volume: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBSVolume1"] = ebsVolume1Metrics

	// Get metrics for EBS2 volume
	ebsVolume2Metrics, err := GetMetrics(clientAuth, instanceId, EBSVolume2Id, elementType, startTime, endTime, "VolumeStalledIOCheck")
	if err != nil {
		log.Println("Error in getting metrics for EBS2 volume: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBSVolume2"] = ebsVolume2Metrics

	// JSON output for volume metrics
	var volumeMetricsOutput volumeMetrics

	if len(rootVolumeMetrics.MetricDataResults) >= 3 &&
		len(rootVolumeMetrics.MetricDataResults[0].Values) >= 1 &&
		len(rootVolumeMetrics.MetricDataResults[1].Values) >= 1 &&
		len(rootVolumeMetrics.MetricDataResults[2].Values) >= 1 {
		volumeMetricsOutput = volumeMetrics{
			RootVolumeUsage: volumeUsage{
				Value: *rootVolumeMetrics.MetricDataResults[0].Values[0],
				Unit:  "GB",
			},
			EBSVolume1Usage: volumeUsage{
				Value: *rootVolumeMetrics.MetricDataResults[1].Values[0],
				Unit:  "GB",
			},
			EBSVolume2Usage: volumeUsage{
				Value: *rootVolumeMetrics.MetricDataResults[2].Values[0],
				Unit:  "GB",
			},
		}
	} else {
		log.Println("Error: Not enough data in MetricDataResults.")
	}

	jsonString, err := json.Marshal(volumeMetricsOutput)
	if err != nil {
		log.Println("Error in marshalling volume metrics json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetMetrics(clientAuth *model.Auth, instanceID, volumeID, elementType string, startTime, endTime *time.Time, metrics ...string) (*cloudwatch.GetMetricDataOutput, error) {
	var metricDataQueries []*cloudwatch.MetricDataQuery

	for i, metricName := range metrics {
		query := &cloudwatch.MetricDataQuery{
			Id: aws.String(fmt.Sprintf("m%d", i+1)),
			MetricStat: &cloudwatch.MetricStat{
				Metric: &cloudwatch.Metric{
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  aws.String("InstanceId"),
							Value: aws.String(instanceID),
						},
						{
							Name:  aws.String("VolumeId"),
							Value: aws.String(volumeID),
						},
					},
					MetricName: aws.String(metricName),
					Namespace:  aws.String("AWS/" + elementType),
				},
				Period: aws.Int64(300),
				Stat:   aws.String("SampleCount"), // You can customize this if needed
			},
		}
		metricDataQueries = append(metricDataQueries, query)
	}

	input := &cloudwatch.GetMetricDataInput{
		EndTime:           endTime,
		StartTime:         startTime,
		MetricDataQueries: metricDataQueries,
	}

	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		log.Printf("Error in GetMetricData: %v", err)
		return nil, err
	}

	return result, nil
}
