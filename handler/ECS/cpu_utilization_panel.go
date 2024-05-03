package ECS

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEcsCpuUtilizationCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetECScpuUtilizationPanel(cmd, clientAuth, nil)
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

func GetECScpuUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	currentUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "SampleCount", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}

	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CurrentUsage"] = currentUsage
	} else {
		log.Println("No data available for current Usage")
	}

	// Get average usage
	averageUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}

	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = averageUsage
	} else {
		log.Println("No data available for average Usage")
	}

	// Get max usage
	maxUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "Maximum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting maximum: ", err)
		return "", nil, err
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MaxUsage"] = maxUsage
	} else {
		log.Println("No data available for maximum Usage")
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
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEcsCpuUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
