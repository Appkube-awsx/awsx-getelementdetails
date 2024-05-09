package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
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

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "CWAgent", "cpu_usage_idle", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
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
	comman_function.InitAwsCmdFlags(AwsxEc2CpuUsageIdleCmd)
}
