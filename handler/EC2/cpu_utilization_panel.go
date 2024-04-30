package EC2

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

var AwsxEc2CpuUtilizationCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetCpuUtilizationPanel(cmd, clientAuth, nil)
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

func GetCpuUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	currentUsage, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "SampleCount", cloudWatchClient)
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
	averageUsage, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "Average", cloudWatchClient)
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
	maxUsage, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "Maximum", cloudWatchClient)
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
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
