package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NetworkTransmitThroughput struct {
	Timestamp time.Time
	Value     float64
}

var AwsxRDSNetworkTransmitThroughputCmd = &cobra.Command{
	Use:   "network_transmit_throughput_panel",
	Short: "get network transmit throughput metrics data",
	Long:  `Command to get network transmit throughput metrics data`,

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
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSNetworkTransmitThroughputPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network transmit throughput data: ", err)
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

func GetRDSNetworkTransmitThroughputPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("Getting cloud-element data from CMDB")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("Using default CMDB URL")
			apiUrl = config.CmdbUrl
		}
		log.Println("CMDB URL: " + apiUrl)
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", "", nil, err
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
			return "", "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data for network transmit throughput metric
	rawData, err := GetNetworkmetricData(clientAuth, elementType, startTime, endTime, "NetworkTransmitThroughput", cloudWatchClient)
	if err != nil {
		log.Println("Error getting network transmit throughput data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["NetworkTransmitThroughput"] = rawData

	// Process raw data
	result := processedRawNetworkTransmitThroughputData(rawData)
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON for network transmit throughput data: ", err)
		return "", "", nil, err
	}

	return string(jsonData), "", cloudwatchMetricData, nil
}

func GetNetworkmetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", elementType, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
						MetricName: aws.String(metricName),
						Namespace:  aws.String("AWS/RDS"),
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

func processedRawNetworkTransmitThroughputData(result *cloudwatch.GetMetricDataOutput) []NetworkTransmitThroughput {
	var processedData []NetworkTransmitThroughput

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		processedData = append(processedData, NetworkTransmitThroughput{
			Timestamp: *timestamp,
			Value:     value,
		})
	}

	return processedData
}

func init() {
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSNetworkTransmitThroughputCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
