package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type CPULimitsResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU limits"`
// }

var AwsxEKSCpuLimitsCmd = &cobra.Command{
	Use:   "cpu_limits_panel",
	Short: "get cpu limits metrics data",
	Long:  `command to get cpu limits metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPULimitsData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu limits data : ", err)
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

func GetCPULimitsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "pod_cpu_limit", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU limits"] = rawData

	// Debug prints
	// log.Printf("RawData Result: %+v", rawData)

	return "", cloudwatchMetricData, nil
}

// func processCPULimitsRawData(result *cloudwatch.GetMetricDataOutput) CPULimitsResult {
// 	var rawData CPULimitsResult
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
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSCpuLimitsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
