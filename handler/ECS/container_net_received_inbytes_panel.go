package ECS

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

// type ContainerNetRxInBytes struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"RawData"`
// }

var AwsxECSContainerNetRxInBytesCmd = &cobra.Command{
	Use:   "container_net_rxinbytes_panel",
	Short: "get container net received inbytes metrics data",
	Long:  `command to get container net received inbytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSContainerNetRxInBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting container net received inbytes metrics data: ", err)
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

func GetECSContainerNetRxInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := metricData.GetMetricClusterData(clientAuth, instanceId, "ECS/ContainerInsights", "NetworkTxBytes", startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting net received bytes data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Container_net_received_inbytes"] = rawData
	return "", cloudwatchMetricData, nil
}

// func processECSContainerNetRxInbytesRawdata(result *cloudwatch.GetMetricDataOutput) ContainerNetRxInBytes {
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
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
