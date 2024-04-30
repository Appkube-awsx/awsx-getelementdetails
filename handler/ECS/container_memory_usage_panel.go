package ECS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type ContainerMemoryUsageResult struct {
// 	TimeSeries []struct {
// 		Timestamp   time.Time
// 		MemoryUsage float64
// 	} `json:"RawData"`
// }

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
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := commanFunction.ParseTimes(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}
	instanceId, err = commanFunction.GetCmdbData(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ECS/ContainerInsights", "MemoryUtilized", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting container memory usage  data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Container_memory_usage"] = rawData
	return "", cloudwatchMetricData, nil

}

// func processContainerMemoryUsageRawData(result *cloudwatch.GetMetricDataOutput) ContainerMemoryUsageResult {
// 	var containerMemoryUsageResult ContainerMemoryUsageResult

// 	for i := range result.MetricDataResults[0].Timestamps {
// 		timestamp := *result.MetricDataResults[0].Timestamps[i]
// 		memoryUsage := *result.MetricDataResults[0].Values[i]
// 		containerMemoryUsageResult.TimeSeries = append(containerMemoryUsageResult.TimeSeries, struct {
// 			Timestamp   time.Time
// 			MemoryUsage float64
// 		}{Timestamp: timestamp, MemoryUsage: memoryUsage})
// 	}

// 	return containerMemoryUsageResult
// }

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
