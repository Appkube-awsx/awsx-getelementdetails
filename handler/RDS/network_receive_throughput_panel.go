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

// type NetworkReceiveThroughput struct {
// 	Timestamp time.Time
// 	Value     float64
// }

var AwsxRDSNetworkReceiveThroughputCmd = &cobra.Command{
	Use:   "network_receive_throughput_panel",
	Short: "get network receive throughput metrics data",
	Long:  `Command to get network receive throughput metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, cloudwatchMetricResp, err := GetRDSNetworkReceiveThroughputPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network receive throughput data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case: print JSON
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetRDSNetworkReceiveThroughputPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricDatabaseData(clientAuth, instanceId, "AWS/RDS", "NetworkReceiveThroughput", startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting network receive throughput data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkReceiveThroughput"] = rawData
	return "", cloudwatchMetricData, nil

}

// func processedRawNetworkReceiveThroughputData(result *cloudwatch.GetMetricDataOutput) []NetworkReceiveThroughput {
// 	var processedData []NetworkReceiveThroughput

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		value := *result.MetricDataResults[0].Values[i]
// 		processedData = append(processedData, NetworkReceiveThroughput{
// 			Timestamp: *timestamp,
// 			Value:     value,
// 		})
// 	}

// 	return processedData
// }

func init() {
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxRDSNetworkReceiveThroughputCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
