package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"

	// "github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NetworkResult struct {
	InboundTraffic  float64 `json:"Inbound Traffic"`
	OutboundTraffic float64 `json:"Outbound Traffic"`
	DataTransferred float64 `json:"DataTransferred"`
}

const (
	bytesToMegabytes = 1024 * 1024
)

var AwsxRDSNetworkUtilizationCmd = &cobra.Command{
	Use:   "network_utilization_panel",
	Short: "get network utilization metrics data",
	Long:  `command to get network utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSNetworkUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network utilization: ", err)
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

func GetRDSNetworkUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	// instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		// cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		// if err != nil {
		// 	return "", nil, err
		// }
		// instanceId = cmdbData.InstanceId

	}

	startTime, endTime, err := comman_function.ParseTimes(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}


	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	inboundTraffic, err := GetRDSNetworkUtilizationMetricData(clientAuth, elementType, startTime, endTime, "Average", "NetworkReceiveThroughput", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}

	// Check if any metric data is returned for inbound traffic
	if len(inboundTraffic.MetricDataResults) == 0 || len(inboundTraffic.MetricDataResults[0].Values) == 0 {
		log.Println("")
		return "null", nil, nil
	}

	// Convert inbound traffic from bytes to megabytes
	inboundTrafficMegabytes := *inboundTraffic.MetricDataResults[0].Values[0] / bytesToMegabytes
	cloudwatchMetricData["Network RX"] = createMetricDataOutput(inboundTrafficMegabytes)

	// Get Outbound Traffic
	outboundTraffic, err := GetRDSNetworkUtilizationMetricData(clientAuth, elementType, startTime, endTime, "Average", "NetworkTransmitThroughput", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}

	// Check if any metric data is returned for outbound traffic
	if len(outboundTraffic.MetricDataResults) == 0 || len(outboundTraffic.MetricDataResults[0].Values) == 0 {
		log.Println("")
		return "null", nil, nil
	}

	// Convert outbound traffic from bytes to megabytes
	outboundTrafficMegabytes := *outboundTraffic.MetricDataResults[0].Values[0] / bytesToMegabytes
	cloudwatchMetricData["Network TX"] = createMetricDataOutput(outboundTrafficMegabytes)

	// Calculate Data Transferred (sum of inbound and outbound) and convert to megabytes
	dataTransferred := inboundTrafficMegabytes + outboundTrafficMegabytes
	cloudwatchMetricData["DataTransferred"] = createMetricDataOutput(dataTransferred)

	jsonOutput := NetworkResult{
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

func GetRDSNetworkUtilizationMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, statistic string, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", elementType, startTime, endTime)
	elmType := "AWS/RDS"
	if elementType == "RDS" {
		elmType = "AWS/" + elementType
	}
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							// {
							// 	Name:  aws.String("InstanceId"),
							// 	Value: aws.String(instanceID),
							// },
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
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
	comman_function.InitAwsCmdFlags(AwsxRDSNetworkUtilizationCmd)
}
