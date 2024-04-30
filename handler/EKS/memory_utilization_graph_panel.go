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

// type MemoryUtilizationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Memory utilization"`
// }

var AwsxEKSMemoryUtilizationGraphCmd = &cobra.Command{
	Use:   "memory_utilization_graph_panel",
	Short: "get memory_utilization graph metrics data",
	Long:  `command to get memory_utilization graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetMemoryUtilizationGraphData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory_utilization: ", err)
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

func GetMemoryUtilizationGraphData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := metricData.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "pod_memory_utilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory utilization"] = rawData

	// result := processMemoryUtilizationGraphRawData(rawData)

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling json in string: ", err)
	// 	return "", nil, err
	// }

	return "", cloudwatchMetricData, nil
}

// func processMemoryUtilizationGraphRawData(result *cloudwatch.GetMetricDataOutput) MemoryUtilizationResult {
// 	var rawData MemoryUtilizationResult
// 	rawData.RawData = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.RawData[i].Timestamp = *timestamp
// 		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
// 	}

// 	return rawData
// }

func init() {
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSMemoryUtilizationGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
