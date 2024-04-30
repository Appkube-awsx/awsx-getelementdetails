package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

//type CpuUsageNice struct {
//	CPU_Nice []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"CPU_Nice"`
//}

var AwsxEc2CpuUsageNiceCmd = &cobra.Command{
	Use:   "cpu_usage_nice_utilization_panel",
	Short: "get cpu usage nice utilization metrics data",
	Long:  `command to get cpu usage nice utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUUsageNicePanel(cmd, clientAuth, nil)
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

func GetCPUUsageNicePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "CWAgent", "cpu_usage_nice", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu usage nice data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_Nice"] = rawData

	//result := processedRawData(rawData)
	//
	//jsonString, err := json.Marshal(result)
	//if err != nil {
	//	log.Println("Error in marshalling json in string: ", err)
	//	return "", nil, err
	//}

	return "", cloudwatchMetricData, nil
}

//	func processedRawData(result *cloudwatch.GetMetricDataOutput) CpuUsageNice {
//		var rawData CpuUsageNice
//		rawData.CPU_Nice = make([]struct {
//			Timestamp time.Time
//			Value     float64
//		}, len(result.MetricDataResults[0].Timestamps))
//
//		for i, timestamp := range result.MetricDataResults[0].Timestamps {
//			rawData.CPU_Nice[i].Timestamp = *timestamp
//			rawData.CPU_Nice[i].Value = *result.MetricDataResults[0].Values[i]
//		}
//
//		return rawData
//	}
func init() {
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
