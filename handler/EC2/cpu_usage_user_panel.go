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

//type CpuUsageUser struct {
//	RawData []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"CPU_User"`
//}

var AwsxEc2CpuUsageUserCmd = &cobra.Command{
	Use:   "cpu_usage_user_utilization_panel",
	Short: "get cpu usage user utilization metrics data",
	Long:  `command to get cpu usage user utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUUsageUserPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu usage user panel utilization: ", err)
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

func GetCPUUsageUserPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "CWAgent", "cpu_usage_user", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_User"] = rawData

	//result := processRawData(rawData)
	//
	//jsonString, err := json.Marshal(result)
	//if err != nil {
	//	log.Println("Error in marshalling json in string: ", err)
	//	return "", nil, err
	//}

	return "", cloudwatchMetricData, nil
}

//
//func processRawData(result *cloudwatch.GetMetricDataOutput) CpuUsageUser {
//	var rawData CpuUsageUser
//	rawData.RawData = make([]struct {
//		Timestamp time.Time
//		Value     float64
//	}, len(result.MetricDataResults[0].Timestamps))
//
//	for i, timestamp := range result.MetricDataResults[0].Timestamps {
//		rawData.RawData[i].Timestamp = *timestamp
//		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
//	}
//
//	return rawData
//}

func init() {
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUsageUserCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
