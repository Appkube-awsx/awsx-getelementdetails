package EKS

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

type NetworkResultMB struct {
	InboundTraffic  float64 `json:"InboundTraffic"`
	OutboundTraffic float64 `json:"OutboundTraffic"`
	DataTransferred float64 `json:"DataTransferred"`
}

// Function to convert bytes to megabytes
func bytesToMegabytes(bytes float64) float64 {
	return bytes / (1024 * 1024)
}

var AwsxEKSNetworkUtilizationCmd = &cobra.Command{
	Use:   "network_utilization_panel",
	Short: "get network_utilization metrics data",
	Long:  `command to get network_utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Network utilization data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNetworkUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	inboundTraffic, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_network_rx_bytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}
	// Convert inbound traffic to megabytes
	var inboundTrafficMegabytes float64
	if len(inboundTraffic.MetricDataResults) > 0 && len(inboundTraffic.MetricDataResults[0].Values) > 0 {
		inboundTrafficMegabytes = bytesToMegabytes(*inboundTraffic.MetricDataResults[0].Values[0])
	} else {
		log.Println("No data available for inbound traffic")
	}
	cloudwatchMetricData["InboundTraffic"] = createMetricDataOutput(inboundTrafficMegabytes)

	// Get Outbound Traffic
	outboundTraffic, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_network_tx_bytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	// Convert outbound traffic to megabytes
	var outboundTrafficMegabytes float64
	if len(outboundTraffic.MetricDataResults) > 0 && len(outboundTraffic.MetricDataResults[0].Values) > 0 {
		outboundTrafficMegabytes = bytesToMegabytes(*outboundTraffic.MetricDataResults[0].Values[0])
	} else {
		log.Println("No data available for outbound traffic")
	}
	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

	var dataTransferred float64
	if len(inboundTraffic.MetricDataResults) > 0 && len(inboundTraffic.MetricDataResults[0].Values) > 0 &&
		len(outboundTraffic.MetricDataResults) > 0 && len(outboundTraffic.MetricDataResults[0].Values) > 0 {
		dataTransferred = *inboundTraffic.MetricDataResults[0].Values[0] + *outboundTraffic.MetricDataResults[0].Values[0]
	} else {
		log.Println("Not enough data available to calculate data transferred")
		dataTransferred = 0
	}
	cloudwatchMetricData["DataTransferred"] = createMetricDataOutput(dataTransferred)

	// Convert values to MB
	dataTransferredMB := bytesToMegabytes(dataTransferred)

	jsonOutput := NetworkResultMB{
		InboundTraffic:  inboundTrafficMegabytes,
		OutboundTraffic: outboundTrafficMegabytes,
		DataTransferred: dataTransferredMB,
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func createMetricDataOutput(value float64) *cloudwatch.GetMetricDataOutput {
	return &cloudwatch.GetMetricDataOutput{
		MetricDataResults: []*cloudwatch.MetricDataResult{
			{
				Values: []*float64{&value},
			},
		},
	}
}

func init() {
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNetworkUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
