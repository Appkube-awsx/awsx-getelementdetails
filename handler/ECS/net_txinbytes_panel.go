package ECS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NetworkTxInBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"RawData"`
// }

var AwsxECSNetworkTxInBytesCmd = &cobra.Command{
	Use:   "network_txinbytes_panel",
	Short: "get network transmitted inbytes metrics data",
	Long:  `command to get network transmitted inbytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSNetworkTxInBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network received inbytes metrics data: ", err)
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

func GetECSNetworkTxInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricClusterData(clientAuth, instanceId, "ECS/ContainerInsights", "NetworkTxBytes", startTime, endTime, "Sum", cloudWatchClient)

	if err != nil {
		log.Println("Error in getting net transmitted bytes data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Network_bytes_transmitted"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processECSNetworkTxInbytesRawdata(result *cloudwatch.GetMetricDataOutput) NetworkTxInBytes {
// 	var rawData NetworkTxInBytes
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
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSNetworkTxInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
