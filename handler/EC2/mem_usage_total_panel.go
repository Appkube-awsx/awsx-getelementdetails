package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type MemUsageTotal struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Mem_Total"`
// }

var AwsxEc2MemoryUsageTotalCmd = &cobra.Command{
	Use:   "memory_usage_panel",
	Short: "get memory usage metrics data",
	Long:  `command to get memory usage metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetMemUsageTotal(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetMemUsageTotal(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	// Fetch raw data
	rawData, err := metricData.GetMetricData(clientAuth, instanceId, "CWAgent", "mem_total", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting memory usage total data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Mem_Total"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processessRawData(result *cloudwatch.GetMetricDataOutput) MemUsageTotal {
// 	var rawData MemUsageTotal
// 	rawData.RawData = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.RawData[i].Timestamp = *timestamp
// 		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
// 	}

// 	return rawData
// }

func init() {
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2MemoryUsageTotalCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
