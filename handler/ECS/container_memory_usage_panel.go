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

type ContainerMemoryUsageResult struct {
	TimeSeries []struct {
		Timestamp   time.Time
		MemoryUsage float64
	}
}

var AwsxECSContainerMemoryUsageCmd = &cobra.Command{
	Use:   "container_memory_usage_panel",
	Short: "get container memory usage metrics data",
	Long:  `command to get container memory usage metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetContainerMemoryUsageData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting container memory usage data : ", err)
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

func GetContainerMemoryUsageData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := GetContainerMemoryUsageMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RawData"] = rawData

	result := processContainerMemoryUsageRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetContainerMemoryUsageMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {

	elmType := "ECS/ContainerInsights"
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

func processContainerMemoryUsageRawData(result *cloudwatch.GetMetricDataOutput) ContainerMemoryUsageResult {
	var containerMemoryUsageResult ContainerMemoryUsageResult

	for i := range result.MetricDataResults[0].Timestamps {
		timestamp := *result.MetricDataResults[0].Timestamps[i]
		memoryUsage := *result.MetricDataResults[0].Values[i]
		containerMemoryUsageResult.TimeSeries = append(containerMemoryUsageResult.TimeSeries, struct {
			Timestamp   time.Time
			MemoryUsage float64
		}{Timestamp: timestamp, MemoryUsage: memoryUsage})
	}

	return containerMemoryUsageResult
}

func init() {
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("query", "", "query")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSContainerMemoryUsageCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
