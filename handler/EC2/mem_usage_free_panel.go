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

// type MemUsageFree struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Mem_Free"`
// }

var AwsxEc2MemoryUsageFreeCmd = &cobra.Command{
	Use:   "memory_usage_free_utilization_panel",
	Short: "get cpu memory usage free utilization metrics data",
	Long:  `command to get cpu usage free utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetMemUsageFreePanel(cmd, clientAuth, nil)
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

func GetMemUsageFreePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := metricData.GetMetricData(clientAuth, instanceId, "CWAgent", "mem_free", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting memeory usage free data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Mem_Free"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processRawDatas(result *cloudwatch.GetMetricDataOutput) MemUsageFree {
// 	var rawData MemUsageFree
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
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2MemoryUsageFreeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
