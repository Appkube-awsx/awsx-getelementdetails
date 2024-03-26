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

type NetworkReceiveThroughput struct {
	Timestamp time.Time
	Value     float64
}

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
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSNetworkReceiveThroughputPanel(cmd, clientAuth, nil)
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

func GetRDSNetworkReceiveThroughputPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data for network receive throughput metric
	rawData, err := GetNetworkMetricdata(clientAuth, elementType, startTime, endTime, "NetworkReceiveThroughput", cloudWatchClient)
	if err != nil {
		log.Println("Error getting network receive throughput data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["NetworkReceiveThroughput"] = rawData

	// Process raw data
	result := processedRawNetworkReceiveThroughputData(rawData)
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling JSON for network receive throughput data: ", err)
		return "", "", nil, err
	}

	return string(jsonData), "", cloudwatchMetricData, nil
}

func GetNetworkMetricdata(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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

func processedRawNetworkReceiveThroughputData(result *cloudwatch.GetMetricDataOutput) []NetworkReceiveThroughput {
	var processedData []NetworkReceiveThroughput

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		processedData = append(processedData, NetworkReceiveThroughput{
			Timestamp: *timestamp,
			Value:     value,
		})
	}

	return processedData
}

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

