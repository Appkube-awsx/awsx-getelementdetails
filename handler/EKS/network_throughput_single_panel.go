package EKS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

const (
	PodNetworkRXByte = "pod_network_rx_bytes"
	PodNetworkTXByte = "pod_network_tx_bytes"
)

type NetworKThroughputResult struct {
	Throughput []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Throughput"`
}

var AwsxEKSNetworkThroughputSingleCmd = &cobra.Command{
	Use:   "network_throughput_single_panel",
	Short: "get Network throughput single graph metrics data",
	Long:  `command to get Network throughput single graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkThroughputSinglePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Network throughput data: ", err)
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

func GetNetworkThroughputSinglePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, string, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing time: %v", err)
	}
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return nil, "", fmt.Errorf("error getting instance ID: %v", err)
	}
	// Fetch network in raw data
	networkInRawData, err := GetMetricData(clientAuth, instanceId, elementType, startTime, endTime, PodNetworkRXByte, cloudWatchClient)
	if err != nil {
		log.Println("Error fetching network in raw data: ", err)
		return nil, "", err
	}

	// Fetch network out raw data
	networkOutRawData, err := GetmetricData(clientAuth, instanceId, elementType, startTime, endTime, PodNetworkTXByte, cloudWatchClient)
	if err != nil {
		log.Println("Error fetching network out raw data: ", err)
		return nil, "", err
	}

	// Calculate network throughput
	result := calculateNetworKThroughput(networkInRawData, networkOutRawData)

	// Marshal result to JSON string
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return nil, "", err
	}

	return networkInRawData, string(jsonString), nil
}

// Function to fetch CloudWatch metric data
func GetmetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"), // Using Sum as an example, change as needed
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

// Function to calculate network throughput
func calculateNetworKThroughput(networkInRawData, networkOutRawData *cloudwatch.GetMetricDataOutput) NetworKThroughputResult {
	var result NetworKThroughputResult

	result.Throughput = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(networkInRawData.MetricDataResults[0].Timestamps))

	for i, timestamp := range networkInRawData.MetricDataResults[0].Timestamps {
		// Calculate network throughput (difference between network in and out)
		throughput := *networkInRawData.MetricDataResults[0].Values[i] - *networkOutRawData.MetricDataResults[0].Values[i]
		result.Throughput[i].Timestamp = *timestamp
		result.Throughput[i].Value = throughput
	}

	return result
}

func init() {
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNetworkThroughputSingleCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
