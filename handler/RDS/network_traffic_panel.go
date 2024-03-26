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

type NetworkTraffic struct {
	Inbound []struct {
		Timestamp time.Time
		Value     float64
	} `json:"network_inbound"`
	Outbound []struct {
		Timestamp time.Time
		Value     float64
	} `json:"network_outbound"`
}

var AwsxRDSNetworkTrafficCmd = &cobra.Command{
	Use:   "network_traffic_panel",
	Short: "get network traffic metrics data",
	Long:  `command to get network traffic metrics data`,

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
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSNetworkTrafficPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network traffic data: ", err)
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

func GetRDSNetworkTrafficPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
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

	// Fetch raw data for inbound and outbound metrics separately
	rawInboundData, err := GetNetworkMetricData(clientAuth, elementType, startTime, endTime, "NetworkReceiveThroughput", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting network inbound data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["Inbound Traffic"] = rawInboundData

	rawOutboundData, err := GetNetworkMetricData(clientAuth, elementType, startTime, endTime, "NetworkTransmitThroughput", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting network outbound data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["Outbound Traffic"] = rawOutboundData

	// Process raw inbound data
	resultInbound := processedRawData(rawInboundData)
	jsonInbound, err := json.Marshal(resultInbound)
	if err != nil {
		log.Println("Error in marshalling json for inbound data: ", err)
		return "", "", nil, err
	}

	// Process raw outbound data
	resultOutbound := processedRawData(rawOutboundData)
	jsonOutbound, err := json.Marshal(resultOutbound)
	if err != nil {
		log.Println("Error in marshalling json for outbound data: ", err)
		return "", "", nil, err
	}
	return string(jsonInbound), string(jsonOutbound), cloudwatchMetricData, nil
}

func GetNetworkMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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

func processedRawData(result *cloudwatch.GetMetricDataOutput) []struct {
	Timestamp time.Time
	Value     float64
} {
	var processedData []struct {
		Timestamp time.Time
		Value     float64
	}

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		// Convert bytes per second to megabytes per second
		valueMBPS := value / (1024 * 1024)
		processedData = append(processedData, struct {
			Timestamp time.Time
			Value     float64
		}{Timestamp: *timestamp, Value: valueMBPS})
	}

	return processedData
}

func init() {
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSNetworkTrafficCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
