package RDS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type DiskQueueDepth struct {
// 	Timestamp time.Time
// 	Value     float64
// }

var AwsxRDSDiskQueueDepthCmd = &cobra.Command{
	Use:   "disk_queue_depth_panel",
	Short: "get disk queue depth metrics data",
	Long:  `command to get disk queue depth metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSDiskQueueDepthPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk queue depth data: ", err)
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

func GetRDSDiskQueueDepthPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/RDS", "DiskQueueDepth", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk queue depth data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DiskQueueDepth"] = rawData
	return "", cloudwatchMetricData, nil

}

// func processedRawDiskQueueDepthData(result *cloudwatch.GetMetricDataOutput) []DiskQueueDepth {
// 	var processedData []DiskQueueDepth

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		value := *result.MetricDataResults[0].Values[i]
// 		processedData = append(processedData, DiskQueueDepth{
// 			Timestamp: *timestamp,
// 			Value:     value,
// 		})
// 	}

// 	return processedData
// }

func init() {
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
