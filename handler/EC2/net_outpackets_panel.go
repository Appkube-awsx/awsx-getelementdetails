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

// type NetworkOutPackets struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Net_Outpackets"`
// }

var AwsxEc2NetworkOutPacketsCmd = &cobra.Command{
	Use:   "network_outpackets_utilization_panel",
	Short: "get network outpackts utilization metrics data",
	Long:  `command to get network outpackets utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkOutPacketsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network outpackets metric data: ", err)
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

func GetNetworkOutPacketsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	
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
	rawData, err := metricData.GetMetricData(clientAuth, instanceId, "AWS/EC2", "NetworkPacketsOut", startTime, endTime, "Sum", cloudWatchClient)
	
	if err != nil {
		log.Println("Error in getting network outpackets data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Net_Outpackets"] = rawData
	return "", cloudwatchMetricData, nil
}


// func processOutPacketsRawData(result *cloudwatch.GetMetricDataOutput) NetworkOutPackets {
// 	var rawData NetworkOutPackets
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
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkOutPacketsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
