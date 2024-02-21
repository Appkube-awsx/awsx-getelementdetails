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
	PodNetworkRXBytes = "pod_network_rx_bytes"
	PodNetworkTXBytes = "pod_network_tx_bytes"
)

type NetworkThroughputResult struct {
	NetworkIn  []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NetworkIn"`
	NetworkOut []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NetworkOut"`
}

var AwsxEKSNetworkThroughputCmd = &cobra.Command{
	Use:   "network_throughput_panel",
	Short: "get Network throughput graph metrics data",
	Long:  `command to get Network throughput graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkThroughputPanel(cmd, clientAuth, nil)
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

func GetNetworkThroughputPanel(cmd *cobra.Command, clientAuth *model.Auth,cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	
	startTime, endTime := parseTime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	networkInRawData, err := GetMetricData(clientAuth, instanceId, elementType, startTime, endTime, PodNetworkRXBytes,cloudWatchClient)
	if err != nil {
		log.Println("Error fetching network in raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkIn"] = networkInRawData

	networkOutRawData, err := GetMetricData(clientAuth, instanceId, elementType, startTime, endTime, PodNetworkTXBytes,cloudWatchClient)
	if err != nil {
		log.Println("Error fetching network out raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkOut"] = networkOutRawData

	result, _ := calculateNetworkThroughput(networkInRawData, networkOutRawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func parseTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		} else {
			startTime = &parsedStartTime
		}
	} else {
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
		now := time.Now()
		endTime = &now
	}

	return startTime, endTime
}

func GetMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, metricName string,cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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

func calculateNetworkThroughput(networkInRawData, networkOutRawData *cloudwatch.GetMetricDataOutput) (NetworkThroughputResult, string) {
	var result NetworkThroughputResult

	result.NetworkIn = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(networkInRawData.MetricDataResults[0].Timestamps))
	for i, timestamp := range networkInRawData.MetricDataResults[0].Timestamps {
		result.NetworkIn[i].Timestamp = *timestamp
		result.NetworkIn[i].Value = *networkInRawData.MetricDataResults[0].Values[i]
	}

	result.NetworkOut = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(networkOutRawData.MetricDataResults[0].Timestamps))
	for i, timestamp := range networkOutRawData.MetricDataResults[0].Timestamps {
		result.NetworkOut[i].Timestamp = *timestamp
		result.NetworkOut[i].Value = *networkOutRawData.MetricDataResults[0].Values[i]
	}

	return result, ""
}

func init() {
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNetworkThroughputCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}