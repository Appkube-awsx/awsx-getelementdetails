package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type DiskUsagePanelData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

func GetDiskUsagePanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

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

	// Fetch raw data for disk read and write
	readRawData, err := GetDiskReadpanelMetricData(clientAuth, instanceId, elementType, startTime, endTime)
	if err != nil {
		log.Println("Error in getting disk read raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DiskReadRawData"] = readRawData

	writeRawData, err := GetDiskWritepanelMetricData(clientAuth, instanceId, elementType, startTime, endTime)
	if err != nil {
		log.Println("Error in getting disk write raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DiskWriteRawData"] = writeRawData

	// Calculate disk usage based on read and write
	result, err := processDiskUsagePanelRawData(readRawData, writeRawData)
	if err != nil {
		log.Println("Error processing disk usage panel raw data: ", err)
		return "", nil, err
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func processDiskUsagePanelRawData(readResult, writeResult *cloudwatch.GetMetricDataOutput) (DiskUsagePanelData, error) {
	var rawData DiskUsagePanelData
	rawData.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, 0) // Initialize with zero length

	// Check if there is data available for disk read
	if len(readResult.MetricDataResults) == 0 || len(readResult.MetricDataResults[0].Timestamps) == 0 {
		return rawData, fmt.Errorf("no data available for disk read")
	}

	// Check if there is data available for disk write
	if len(writeResult.MetricDataResults) == 0 || len(writeResult.MetricDataResults[0].Timestamps) == 0 {
		return rawData, fmt.Errorf("no data available for disk write")
	}

	// Assuming the length of timestamps for read and write data are the same
	for i, readTimestamp := range readResult.MetricDataResults[0].Timestamps {
		// Check if the index is valid for both read and write data
		if i >= len(writeResult.MetricDataResults[0].Values) {
			return rawData, fmt.Errorf("index out of range for write data")
		}

		writeValue := *writeResult.MetricDataResults[0].Values[i]
		readValue := *readResult.MetricDataResults[0].Values[i]
		// Calculate disk usage as the sum of read and write
		totalDiskUsage := readValue + writeValue
		rawData.RawData = append(rawData.RawData, struct {
			Timestamp time.Time
			Value     float64
		}{
			Timestamp: *readTimestamp,
			Value:     totalDiskUsage,
		})
	}

	return rawData, nil
}

func GetDiskReadpanelMetricData(clientAuth *model.Auth, instanceID string, namespace string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
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

func GetDiskWritepanelMetricData(clientAuth *model.Auth, instanceID string, elementType string, startTime, endTime *time.Time) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("DiskWriteBytes"),
						Namespace:  aws.String("AWS/" + elementType),
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
