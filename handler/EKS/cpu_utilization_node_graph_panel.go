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

// type CPU_UtilizationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU Utilization"`
// }

var AwsxEKSCpuUtilizationNodeGraphCmd = &cobra.Command{
	Use:   "cpu_utilization_node_graph_panel",
	Short: "get cpu utilization node graph metrics data",
	Long:  `command to get cpu utilization node graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUUtilizationNodeData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization node graph data : ", err)
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

func GetCPUUtilizationNodeData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "node_cpu_utilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU Utilization"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processCPU_UtilizationRawData(result *cloudwatch.GetMetricDataOutput) CPU_UtilizationResult {
// 	var rawData CPU_UtilizationResult
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
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSCpuUtilizationNodeGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
