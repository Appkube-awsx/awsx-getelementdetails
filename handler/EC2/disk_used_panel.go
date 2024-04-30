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

//type DiskUsedPanelData struct {
//	RawData []struct {
//		Timestamp time.Time
//		Value     float64
//	} `json:"Disk_Used"`
//}

var AwsxEc2DiskUsedCmd = &cobra.Command{
	Use:   "disk_used_panel",
	Short: "get disk used metrics data",
	Long:  `command to get disk used metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskUsedPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk read  utilization: ", err)
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

func GetDiskUsedPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "CWAgent", "disk_used", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk used data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Disk_Used"] = rawData

	//result := processDiskUsedPanelRawData(rawData)
	//
	//jsonString, err := json.Marshal(result)
	//if err != nil {
	//	log.Println("Error in marshalling json in string: ", err)
	//	return "", nil, err
	//}

	return "", cloudwatchMetricData, nil
}

//
//
//func processDiskUsedPanelRawData(result *cloudwatch.GetMetricDataOutput) DiskUsedPanelData {
//	var rawData DiskUsedPanelData
//
//	// Initialize an empty slice to store the raw data
//	rawData.RawData = []struct {
//		Timestamp time.Time
//		Value     float64
//	}{}
//
//	// Iterate over each metric data result
//	for _, metricDataResult := range result.MetricDataResults {
//		// Iterate over each timestamp and value pair in the current metric data result
//		for i, timestamp := range metricDataResult.Timestamps {
//			// Append the timestamp and value to the rawData slice
//			rawData.RawData = append(rawData.RawData, struct {
//				Timestamp time.Time
//				Value     float64
//			}{
//				Timestamp: *timestamp,
//				Value:     *metricDataResult.Values[i],
//			})
//		}
//	}
//
//	return rawData
//}

func init() {
	AwsxEc2DiskUsedCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
