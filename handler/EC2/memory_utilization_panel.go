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

func GetMemoryUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID, _ := cmd.PersistentFlags().GetString("instanceID")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

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

	currentUsage, err := GetMemoryUtilizationMetricData(clientAuth, instanceID, namespace, startTime, endTime, "SampleCount")
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CurrentUsage"] = currentUsage
	} else {
		log.Println("No data found for current usage")
	}

	// Get average utilization
	averageUsage, err := GetMemoryUtilizationMetricData(clientAuth, instanceID, namespace, startTime, endTime, "Average")
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = averageUsage
	} else {
		log.Println("No data found for average usage")
	}

	maxUsage, err := GetMemoryUtilizationMetricData(clientAuth, instanceID, namespace, startTime, endTime, "Maximum")
	if err != nil {
		log.Println("Error in getting maximum: ", err)
		return "", nil, err
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MaxUsage"] = maxUsage
	} else {
		log.Println("No data found for maximum usage")
	}

	jsonOutput := make(map[string]float64)

	if currentUsage != nil && len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["CurrentUsage"] = *currentUsage.MetricDataResults[0].Values[0]
	}
	if averageUsage != nil && len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["AverageUsage"] = *averageUsage.MetricDataResults[0].Values[0]
	}
	if maxUsage != nil && len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["MaxUsage"] = *maxUsage.MetricDataResults[0].Values[0]
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetMemoryUtilizationMetricData(clientAuth *model.Auth, instanceID, namespace string, startTime, endTime *time.Time, statistic string) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("mem_used_percent"),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
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
