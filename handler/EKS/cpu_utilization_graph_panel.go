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

// type CPUUtilizationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU Utilization"`
// }

var AwsxEKSCpuUtilizationGraphCmd = &cobra.Command{
	Use:   "cpu_utilization_graph_panel",
	Short: "get cpu utilization graph metrics data",
	Long:  `command to get cpu utilization graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUUtilizationData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization graph data : ", err)
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

func GetCPUUtilizationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "pod_cpu_utilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU Utilization"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processCPUUtilizationRawData(result *cloudwatch.GetMetricDataOutput) CPUUtilizationResult {
// 	var rawData CPUUtilizationResult
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
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSCpuUtilizationGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
