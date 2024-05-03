package EC2

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

var AwsxEc2DiskReadCmd = &cobra.Command{
	Use:   "disk_read_panel",
	Short: "get disk read metrics data",
	Long:  `command to get disk read metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskReadPanel(cmd, clientAuth, nil)
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

func GetDiskReadPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "DiskReadBytes", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk read data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Disk_Reads"] = rawData
	jsonOutput := make(map[string]float64)
	if len(rawData.MetricDataResults) > 0 && len(rawData.MetricDataResults[0].Values) > 0 {
		jsonOutput["Disk_Reads"] = *rawData.MetricDataResults[0].Values[0]
	}
	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func init() {
	AwsxEc2DiskReadCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2DiskReadCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2DiskReadCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2DiskReadCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2DiskReadCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2DiskReadCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2DiskReadCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2DiskReadCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2DiskReadCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2DiskReadCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2DiskReadCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2DiskReadCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2DiskReadCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2DiskReadCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2DiskReadCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2DiskReadCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
