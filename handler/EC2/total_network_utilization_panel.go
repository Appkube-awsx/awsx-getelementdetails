package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEc2NetworkUtilizationAcrossAllInstanceCmd = &cobra.Command{
	Use:   "total_network_utilization_panel",
	Short: "get network utilization metrics data for all instances",
	Long:  `command to get network utilization metrics data for all instances`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkUtilizationAcrossAllInstancesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network utilization: ", err)
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

func GetNetworkUtilizationAcrossAllInstancesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	inboundTraffic, err := GetTotalNetworkUtilizationData(clientAuth, elementType, startTime, endTime, "Average", "NetworkIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}
	if len(inboundTraffic.MetricDataResults) == 0 || len(inboundTraffic.MetricDataResults[0].Values) == 0 {
		return "null", nil, nil
	}

	// Get Outbound Traffic
	outboundTraffic, err := GetTotalNetworkUtilizationData(clientAuth, elementType, startTime, endTime, "Average", "NetworkOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	if len(outboundTraffic.MetricDataResults) == 0 || len(outboundTraffic.MetricDataResults[0].Values) == 0 {
		return "null", nil, nil
	}

	// Convert traffic from bytes to megabytes and collect time series data
	inboundSeries := convertToTimeSeries(inboundTraffic.MetricDataResults[0])
	outboundSeries := convertToTimeSeries(outboundTraffic.MetricDataResults[0])

	cloudwatchMetricData["InboundTraffic"] = inboundTraffic
	cloudwatchMetricData["OutboundTraffic"] = outboundTraffic

	jsonOutput := map[string]interface{}{
		"InboundTraffic":  inboundSeries,
		"OutboundTraffic": outboundSeries,
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetTotalNetworkUtilizationData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, statistic string, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data in namespace %s from %v to %v", elementType, startTime, endTime)
	elmType := "AWS/EC2"
	if elementType == "EC2" {
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

func convertToTimeSeries(metricDataResult *cloudwatch.MetricDataResult) []map[string]interface{} {
	timeSeries := []map[string]interface{}{}
	for i, timestamp := range metricDataResult.Timestamps {
		value := *metricDataResult.Values[i] / bytesToMegabytes
		timeSeries = append(timeSeries, map[string]interface{}{
			"timestamp": timestamp,
			"value":     value,
		})
	}
	return timeSeries
}


func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2NetworkUtilizationAcrossAllInstanceCmd)
}
