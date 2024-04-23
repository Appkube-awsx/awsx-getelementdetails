package ECS

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEcsMemoryUtilizationCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization: ", err)
				fmt.Println("null")
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetMemoryUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	elementType, _ := cmd.PersistentFlags().GetString("elementType")

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

	currentUsage, err := metricData.GetMetricClusterData(clientAuth, instanceId, "AWS/"+elementType, "MemoryUtilization", startTime, endTime, "SampleCount", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CurrentUsage"] = currentUsage
	} else {
		log.Println("No data found for current usage")
	}

	// cloudwatchMetricData["CurrentUsage"] = &cloudwatch.GetMetricDataOutput{
	// 	MetricDataResults: []*cloudwatch.MetricDataResult{{Values: []*float64{aws.Float64(0)}}},
	// }

	// Get average utilization
	averageUsage, err := metricData.GetMetricClusterData(clientAuth, instanceId, "AWS/"+elementType, "MemoryUtilization", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = averageUsage
	} else {
		log.Println("No data found for average usage")
	}
	// cloudwatchMetricData["CurrentUsage"] = &cloudwatch.GetMetricDataOutput{
	// 	MetricDataResults: []*cloudwatch.MetricDataResult{{Values: []*float64{aws.Float64(0)}}},
	// }

	maxUsage, err := metricData.GetMetricClusterData(clientAuth, instanceId, "AWS/"+elementType, "MemoryUtilization", startTime, endTime, "Maximum", cloudWatchClient)
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
	// cloudwatchMetricData["CurrentUsage"] = &cloudwatch.GetMetricDataOutput{
	// 	MetricDataResults: []*cloudwatch.MetricDataResult{{Values: []*float64{aws.Float64(0)}}},
	// }

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
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEcsMemoryUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
