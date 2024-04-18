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

// type NetworkOutBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Net_Outbytes"`
// }

var AwsxEc2NetworkOutBytesCmd = &cobra.Command{
	Use:   "network_outbytes_utilization_panel",
	Short: "get network outbytes utilization metrics data",
	Long:  `command to get network out bytes utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkOutBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network outbytes metrics data: ", err)
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

func GetNetworkOutBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	
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
		log.Println("Error in getting network outbytes data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Net_Outbytes"] = rawData
	return "", cloudwatchMetricData, nil
}



// func processOutbytesRawdata(result *cloudwatch.GetMetricDataOutput) NetworkOutBytes {
// 	var rawData NetworkOutBytes
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
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkOutBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
