package EKS

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEKSMemoryUtilizationCmd = &cobra.Command{
	Use:   "memory_utilization_panel",
	Short: "get memory utilization metrics data",
	Long:  `command to get memory utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GeteksMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization: ", err)
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

func GeteksMemoryUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	currentUsage, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "node_memory_utilization", startTime, endTime, "SampleCount", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CurrentUsage"] = currentUsage
	} else {
		log.Println("No data found for current usage")
	}
	// Get average utilization
	averageUsage, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "node_memory_utilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = averageUsage
	} else {
		log.Println("No data found for average usage")
	}
	maxUsage, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "node_memory_utilization", startTime, endTime, "Maximum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting maximum: ", err)
		return "", nil, err
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MaxUsage"] = maxUsage
	} else {
		log.Println("")
		return "null", nil, nil
	}
	jsonOutput := make(map[string]float64)
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["CurrentUsage"] = *currentUsage.MetricDataResults[0].Values[0]
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["AverageUsage"] = *averageUsage.MetricDataResults[0].Values[0]
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["MaxUsage"] = *maxUsage.MetricDataResults[0].Values[0]
	}
	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil

}

func init() {
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
