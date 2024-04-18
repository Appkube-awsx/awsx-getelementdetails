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

// type Networkoutbound struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"NetworkOutbound"`
// }

var AwsxEc2NetworkOutboundCmd = &cobra.Command{
	Use:   "network_out_bound_panel",
	Short: "get network out bound metrics data",
	Long:  `command to get network out bound metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkOutBoundPanel(cmd, clientAuth, nil)
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

func GetNetworkOutBoundPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := metricData.GetMetricData(clientAuth, instanceId, "AWS/EC2", "NetworkOut", startTime, endTime, "Sum", cloudWatchClient)

	if err != nil {
		log.Println("Error in network outbounds data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["NetworkOutbound"] = rawData
	return "", cloudwatchMetricData, nil
}

// func processtheRawdata(result *cloudwatch.GetMetricDataOutput) Networkoutbound {
// 	var rawData Networkoutbound
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
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkOutboundCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
