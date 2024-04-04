package ECS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type AvailableMemoryOverTimeResult struct {
	TimeSeries []struct {
		Timestamp       time.Time
		AvailableMemory float64
	}
}

type AllocateResult struct {
	RawData []AvailableMemoryOverTimeResult `json:"RawData"`
}

var AwsxECSAvailableMemoryOverTimeCmd = &cobra.Command{
	Use:   "available_memory_overtime_panel",
	Short: "get available memory over time metrics data",
	Long:  `command to get available memory over time metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")
		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err := GetAvailableMemoryOverTimeData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting available memory over time data : ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetAvailableMemoryOverTimeData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

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

	// Fetch raw data
	rawData, err := GetAvailableMemoryOverTimeMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	// Process the raw data if needed
	result := processAvailableMemoryOverTimeRawData(rawData)

	// Collect all timestamps and values separately
	timestamps := make([]time.Time, len(result.TimeSeries))
	values := make([]float64, len(result.TimeSeries))

	// Populate the slices with actual data
	for i, data := range result.TimeSeries {
		// Assigning values directly to slices without taking their addresses
		timestamps[i] = data.Timestamp
		values[i] = data.AvailableMemory
	}

	// Initialize the MetricDataResults slice
	metricDataResults := make([]*cloudwatch.MetricDataResult, len(result.TimeSeries))

	// Populate the MetricDataResults with actual data
	for i := range result.TimeSeries {
		metricDataResults[i] = &cloudwatch.MetricDataResult{
			Timestamps: []*time.Time{&timestamps[i]},
			Values:     []*float64{&values[i]},
		}
	}

	// Assign the processed data to cloudwatchMetricData
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"RawData": {
			MetricDataResults: metricDataResults,
		},
	}

	// Initialize AllocateResult and populate RawData field
	allocateResult := AllocateResult{
		RawData: []AvailableMemoryOverTimeResult{result}, // Convert single result to slice and assign
	}

	// Convert the result to JSON
	jsonString, err := json.Marshal(allocateResult)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetAvailableMemoryOverTimeMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {

	elmType := "ECS/ContainerInsights"
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("memory_reserved"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("MemoryReserved"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("memory_utilized"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("MemoryUtilized"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processAvailableMemoryOverTimeRawData(result *cloudwatch.GetMetricDataOutput) AvailableMemoryOverTimeResult {
	var availableMemoryOverTimeResult AvailableMemoryOverTimeResult

	// Iterate over all timestamps and values
	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		if i >= len(result.MetricDataResults[1].Values) {
			break
		}
		if i >= len(result.MetricDataResults[0].Values) {
			break
		}

		memoryUtilized := *result.MetricDataResults[1].Values[i]
		memoryReservation := *result.MetricDataResults[0].Values[i]
		availableMemory := calculateAvailableMemory(memoryUtilized, memoryReservation)

		availableMemoryOverTimeResult.TimeSeries = append(availableMemoryOverTimeResult.TimeSeries, struct {
			Timestamp       time.Time
			AvailableMemory float64
		}{
			Timestamp:       *timestamp,
			AvailableMemory: availableMemory,
		})
	}

	return availableMemoryOverTimeResult
}

func calculateAvailableMemory(memoryUtilized, memoryReservation float64) float64 {
	return memoryReservation - memoryUtilized
}

func init() {
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("query", "", "query")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSAvailableMemoryOverTimeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
