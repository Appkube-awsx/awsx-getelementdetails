package ECS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NetworkRxInBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"RawData"`
// }

var AwsxECSNetworkRxInBytesCmd = &cobra.Command{
	Use:   "network_rxinbytes_panel",
	Short: "get network received inbytes metrics data",
	Long:  `command to get network received inbytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSNetworkRxInBytesPanel(cmd, clientAuth, nil)
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

func GetECSNetworkRxInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "ECS/ContainerInsights", "NetworkRxBytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting net received bytes data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Network_bytes_received"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processECSNetworkRxInbytesRawdata(result *cloudwatch.GetMetricDataOutput) NetworkRxInBytes {
// 	var rawData NetworkRxInBytes
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
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSNetworkRxInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
