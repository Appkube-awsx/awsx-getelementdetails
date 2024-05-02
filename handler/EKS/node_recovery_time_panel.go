package EKS

import (
	"fmt"
	"log"

	//"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEKSNodeRecoveryPanelCmd = &cobra.Command{
	Use:   "node_recovery_time_panel",
	Short: "get node recovery time metrics data",
	Long:  `command to get node recovery time metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNodeRecoveryTime(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting node recovery time data: ", err)
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

func GetNodeRecoveryTime(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_status_condition_ready", startTime, endTime, "Maximum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_User"] = rawData

	return "", cloudwatchMetricData, nil

}

// func ProcessNodeReadyData(result *cloudwatch.GetMetricDataOutput) []NodeRecoveryData {
// 	var recoveryTimeSeries []NodeRecoveryData

// 	for i := 1; i < len(result.MetricDataResults[0].Timestamps); i++ {
// 		currentTimestamp := *result.MetricDataResults[0].Timestamps[i]
// 		previousTimestamp := *result.MetricDataResults[0].Timestamps[i-1]

// 		if *result.MetricDataResults[0].Values[i-1] == 0 && *result.MetricDataResults[0].Values[i] == 1 {
// 			recoveryTime := currentTimestamp.Sub(previousTimestamp)

// 			recoveryTimeSeries = append(recoveryTimeSeries, NodeRecoveryData{
// 				Timestamp:    currentTimestamp,
// 				RecoveryTime: recoveryTime,
// 			})
// 		}
// 	}

// 	return recoveryTimeSeries
// }

func init() {
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeDowntimeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
