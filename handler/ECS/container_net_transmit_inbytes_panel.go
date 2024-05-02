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

// type ContainerNetTxInBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"RawData"`
// }

var AwsxECSContainerNetTxInBytesCmd = &cobra.Command{
	Use:   "container_net_txinbytes_panel",
	Short: "get container net transmit inbytes metrics data",
	Long:  `command to get container net transmit inbytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSContainerNetTxInBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting container net transmit inbytes metrics data: ", err)
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

func GetECSContainerNetTxInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ECS/ContainerInsights", "NetworkRxBytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)

	if err != nil {
		log.Println("Error in getting net transmitted bytes data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Container_net_transmit_inbytes"] = rawData
	return "", cloudwatchMetricData, nil
}

// func processECSContainerNetTxInbytesRawdata(result *cloudwatch.GetMetricDataOutput) ContainerNetRxInBytes {
// 	var rawData ContainerNetRxInBytes
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
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
