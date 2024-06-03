package EKS

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
)

type NetworkResults struct {
	InboundTraffic  float64 `json:"InboundTraffic"`
	OutboundTraffic float64 `json:"OutboundTraffic"`
	DataTransferred float64 `json:"DataTransferred"`
}

const (
	bytesToMegabytes = 1024 * 1024
)

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

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	inboundTraffic, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_network_rx_bytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}
	inboundTrafficMegabytes := *inboundTraffic.MetricDataResults[0].Values[0] / bytesToMegabytes
	cloudwatchMetricData["InboundTraffic"] = createMetricDataOutput(inboundTrafficMegabytes)

	// Get Outbound Traffic
	outboundTraffic, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "pod_network_tx_bytes", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	outboundTrafficMegabytes := *outboundTraffic.MetricDataResults[0].Values[0] / bytesToMegabytes
	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

	dataTransferred := inboundTrafficMegabytes + outboundTrafficMegabytes
	cloudwatchMetricData["DataTransferred"] = createMetricDataOutput(dataTransferred)

	jsonOutput := NetworkResults{
		InboundTraffic:  inboundTrafficMegabytes,
		OutboundTraffic: outboundTrafficMegabytes,
		DataTransferred: dataTransferred,
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
	comman_function.InitAwsCmdFlags(AwsxEKSNetworkUtilizationCmd)
}
