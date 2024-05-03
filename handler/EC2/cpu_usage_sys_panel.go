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

//type CpuUsageSys struct {
//	CPU_Sys []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"CPU_Sys"`
//}

var AwsxEc2CpuSysTimeCmd = &cobra.Command{
	Use:   "cpu_sys_time_utilization_panel",
	Short: "get cpu sys time utilization metrics data",
	Long:  `command to get cpu sys time utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUUsageSysPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu sys time utilization: ", err)
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

func GetCPUUsageSysPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, elementType, "cpu_usage_system", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu usage system data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_Sys"] = rawData
	//
	//result := processingRawData(rawData)
	//
	//jsonString, err := json.Marshal(result)
	//if err != nil {
	//	log.Println("Error in marshalling json in string: ", err)
	//	return "", nil, err
	//}

	return "", cloudwatchMetricData, nil
}

//	func processingRawData(result *cloudwatch.GetMetricDataOutput) CpuUsageSys {
//		var rawData CpuUsageSys
//		rawData.CPU_Sys = make([]struct {
//			Timestamp time.Time
//			Value     float64
//		}, len(result.MetricDataResults[0].Timestamps))
//
//		for i, timestamp := range result.MetricDataResults[0].Timestamps {
//			rawData.CPU_Sys[i].Timestamp = *timestamp
//			rawData.CPU_Sys[i].Value = *result.MetricDataResults[0].Values[i]
//		}
//
//		return rawData
//	}
func init() {
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuSysTimeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
