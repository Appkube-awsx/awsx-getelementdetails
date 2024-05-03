package EKS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type cpuResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU requests"`
// }

var AwsxEKSCpuRequestsCmd = &cobra.Command{
	Use:   "cpu_requests_panel",
	Short: "get cpu requests metrics data",
	Long:  `command to get cpu requests metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPURequestData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu requests data : ", err)
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

func GetCPURequestData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_cpu_request", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU requests"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processRawData(result *cloudwatch.GetMetricDataOutput) cpuResult {
// 	var rawData cpuResult
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
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSCpuRequestsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
