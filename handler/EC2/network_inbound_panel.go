package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NetworkInbound struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"NetworkInbound"`
// }

var AwsxEc2NetworkInboundCmd = &cobra.Command{
	Use:   "network_in_bound_panel",
	Short: "get network in bound metrics data",
	Long:  `command to get network in bound metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkInBoundPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network in bytes metrics data: ", err)
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

func GetNetworkInBoundPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/EC2", "NetworkIn", startTime, endTime, "Sum", "InstanceId", cloudWatchClient)

	if err != nil {
		log.Println("Error in getting network inbounds data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["NetworkInbound"] = rawData
	return "", cloudwatchMetricData, nil
}

// func processTheRawdata(result *cloudwatch.GetMetricDataOutput) NetworkInbound {
// 	var rawData NetworkInbound
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
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkInboundCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
