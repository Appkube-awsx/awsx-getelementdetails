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

// type MemoryGraphUtilizationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Memory utilization"`
// }

var AwsxEc2MemoryUtilizationGraphCmd = &cobra.Command{
	Use:   "memory_utilization_graph_panel",
	Short: "get memory utilization graph metrics data",
	Long:  `command to get memory utilization graph metrics data`,
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
			jsonResp, cloudwatchMetricResp, err := GetMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization graph: ", err)
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

func GetMemoryUtilizationGraphPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	// Get average utilization
	rawData, err := metricData.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "mem_used_percent", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting rawdata: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory utilization"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processMemoryUtilizationGraphRawData(result *cloudwatch.GetMetricDataOutput) MemoryGraphUtilizationResult {
// 	var rawData MemoryGraphUtilizationResult
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
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2MemoryUtilizationGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
