package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

//type CpuUsageIdle struct {
//	CPU_Idle []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"CPU_Idle"`
//}

var AwsxEc2CpuUsageIdleCmd = &cobra.Command{
	Use:   "cpu_usage_Idle_utilization_panel",
	Short: "get cpu usage idle utilization metrics data",
	Long:  `command to get cpu usage idle utilization metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command..")
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
			jsonResp, cloudwatchMetricResp, err := GetCPUUsageIdlePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu usage idle utilization: ", err)
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

func GetCPUUsageIdlePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawData, err := metricData.GetMetricData(clientAuth, instanceId, "CWAgent", "cpu_usage_idle", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu usage idle data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_Idle"] = rawData

	return "", cloudwatchMetricData, nil
}

//	func processTheRawData(result *cloudwatch.GetMetricDataOutput) CpuUsageIdle {
//		var rawData CpuUsageIdle
//		rawData.CPU_Idle = make([]struct {
//			Timestamp time.Time
//			Value     float64
//		}, len(result.MetricDataResults[0].Timestamps))
//
//		for i, timestamp := range result.MetricDataResults[0].Timestamps {
//			rawData.CPU_Idle[i].Timestamp = *timestamp
//			rawData.CPU_Idle[i].Value = *result.MetricDataResults[0].Values[i]
//		}
//
//		return rawData
//	}
func init() {
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
