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

// type NodeStabilityResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"NodeStabilityindex"`
// }

var AwsxEKSNodeStabilityCmd = &cobra.Command{
	Use:   "node_stability_panel",
	Short: "get node stability metrics data",
	Long:  `command to get node stability metrics data`,

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

func GetNodeStabilityData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_number_of_running_containers", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NodeStabilityindex"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processNodeStabilityRawData(result *cloudwatch.GetMetricDataOutput) NodeStabilityResult {
// 	var rawData NodeStabilityResult
// 	rawData.RawData = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.RawData[i].Timestamp = *timestamp
// 		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i] / 60
// 	}

// 	return rawData
// }

func init() {
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeStabilityCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
