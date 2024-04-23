package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type TimeSeriesData struct {
// 	Timestamp      time.Time
// 	AllocatableCPU float64
// }

// type AllocateResult struct {
// 	AllocatableCPU []TimeSeriesData `json:"AllocatableCPU"`
// }

var AwsxEKSAllocatableCpuCmd = &cobra.Command{
	Use:   "allocatable_cpu_panel",
	Short: "get allocatable cpu metrics data",
	Long:  `command to get allocatable cpu metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetAllocatableCPUData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting allocatable cpu: ", err)
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

func GetAllocatableCPUData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

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
	rawData, err := metricData.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "node_cpu_limit", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Allocatble_CPU"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processCPURawData(result *cloudwatch.GetMetricDataOutput) AllocateResult {
// 	var rawData AllocateResult
// 	rawData.AllocatableCPU = make([]TimeSeriesData, len(result.MetricDataResults[0].Timestamps))

// 	// Assuming the two metrics have the same number of data points
// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.AllocatableCPU[i].Timestamp = *timestamp
// 		cpuLimit := *result.MetricDataResults[0].Values[i]
// 		reservedCapacity := *result.MetricDataResults[1].Values[i]

// 		// Log the values for troubleshooting
// 		// log.Printf("Timestamp: %v, cpuLimit: %v, reservedCapacity: %v", *timestamp, cpuLimit, reservedCapacity)

// 		allocatableCPU := cpuLimit - reservedCapacity
// 		// log.Println(allocatableCPU)

// 		// Only include the calculated allocatable CPU in the result
// 		rawData.AllocatableCPU[i].AllocatableCPU = allocatableCPU
// 	}
// 	// log.Println(rawData)
// 	return rawData
// }

func init() {
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
