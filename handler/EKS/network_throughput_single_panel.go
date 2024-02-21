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
			return nil,"",err
		}
		instanceId = cmdbData.InstanceId

	}

	startTime, endTime := ParseTime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Fetch network in raw data
	networkInRawData, err := GetmetricData(clientAuth, instanceId, elementType, startTime, endTime, PodNetworkRXByte, cloudWatchClient)
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

// Function to parse time strings and return time pointers
func ParseTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		} else {
			startTime = &parsedStartTime
		}
	} else {
		// If startTimeStr is empty, default to the last five minutes
		now := time.Now()
		startTime = &now
		minusFiveMinutes := now.Add(-5 * time.Minute)
		startTime = &minusFiveMinutes
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
		} else {
			endTime = &parsedEndTime
		}
	} else {
		// If endTimeStr is empty, default to the current time
		now := time.Now()
		endTime = &now
	}

	return startTime, endTime
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