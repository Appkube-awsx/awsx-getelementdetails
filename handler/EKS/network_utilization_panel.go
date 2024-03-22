package EKS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
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
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, err
		}
		instanceId = cmdbData.InstanceId

	}

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	inboundTraffic, err := GetNetworkMetricData(clientAuth, instanceId, elementType, startTime, endTime, "pod_network_rx_bytes", cloudWatchClient)
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
	outboundTraffic, err := GetNetworkMetricData(clientAuth, instanceId, elementType, startTime, endTime, "pod_network_tx_bytes", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

	// Calculate Data Transferred (sum of inbound and outbound)
	dataTransferred := *inboundTraffic.MetricDataResults[0].Values[0] + *outboundTraffic.MetricDataResults[0].Values[0]
	cloudwatchMetricData["DataTransferred"] = createMetricDataOutput(dataTransferred)

	// Convert values to MB
	outboundTrafficMB := bytesToMegabytes(*outboundTraffic.MetricDataResults[0].Values[0])
	dataTransferredMB := bytesToMegabytes(dataTransferred)

	jsonOutput := NetworkResultMB{
		InboundTraffic:  inboundTrafficMegabytes,
		OutboundTraffic: outboundTrafficMB,
		DataTransferred: dataTransferredMB,
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNetworkMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	elmType := "ContainerInsights"
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Sum"),
				},
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func extractMetricValue(result *cloudwatch.GetMetricDataOutput, index int) float64 {
	if len(result.MetricDataResults) > 0 && len(result.MetricDataResults[0].Values) > index {
		return *result.MetricDataResults[0].Values[index]
	}
	return 0
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
