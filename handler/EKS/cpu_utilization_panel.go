package EKS

import (
	"encoding/json"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"

	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type Result struct {
// 	CurrentUsage float64 `json:"CurrentUsage"`
// 	AverageUsage float64 `json:"AverageUsage"`
// 	MaxUsage     float64 `json:"MaxUsage"`
// }

var AwsxEKSCpuUtilizationCmd = &cobra.Command{
	Use:   "cpu_utilization_panel",
	Short: "get cpu utilization metrics data",
	Long:  `command to get cpu utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetEKScpuUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetEKScpuUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	//if queryName == "cpu_utilization_panel" {
	currentUsage, err := metricData.GetMetricClusterData(clientAuth, instanceId, "AWS/"+elementType, "node_cpu_utilization", startTime, endTime, "SampleCount", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CurrentUsage"] = currentUsage
	// Get average usage
	averageUsage, err := metricData.GetMetricClusterData(clientAuth, instanceId, "AWS/"+elementType, "node_cpu_utilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["AverageUsage"] = averageUsage
	// Get max usage
	maxUsage, err := metricData.GetMetricClusterData(clientAuth, instanceId, "AWS/"+elementType, "node_cpu_utilization", startTime, endTime, "Maximum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting maximum: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["MaxUsage"] = maxUsage
	jsonOutput := map[string]float64{
		"CurrentUsage": *currentUsage.MetricDataResults[0].Values[0],
		"AverageUsage": *averageUsage.MetricDataResults[0].Values[0],
		"MaxUsage":     *maxUsage.MetricDataResults[0].Values[0],
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil

}

func init() {
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSCpuUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
