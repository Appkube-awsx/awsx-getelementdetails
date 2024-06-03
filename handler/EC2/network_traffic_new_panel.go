package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NetworkTrafficALL struct {
	InboundTraffic  float64 `json:"inbound_traffic"`
	OutboundTraffic float64 `json:"outbound_traffic"`
}

var AwsxEC2NetworkTrafficCmdAllinstances = &cobra.Command{
	Use:   "network_traffic_new_panel",
	Short: "Get network traffic metrics data",
	Long:  `Command to get network traffic metrics data`,

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
			networkTraffic, err := GetNetworkTrafficNewPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network traffic data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println("This cli is for only json output")
			} else {
				jsonResp, err := json.Marshal(networkTraffic)
				if err != nil {
					log.Println("Error marshalling network traffic data: ", err)
					return
				}
				fmt.Println(string(jsonResp))
			}
		}
	},
}

func GetNetworkTrafficNewPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, error) {
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
			return "", err
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
			return "", err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Fetch raw data for inbound and outbound metrics separately
	rawInboundData, err := GetAllNetworkMetricData(clientAuth, elementType, startTime, endTime, "NetworkIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting network inbound data: ", err)
		return "", err
	}

	rawOutboundData, err := GetAllNetworkMetricData(clientAuth, elementType, startTime, endTime, "NetworkOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting network outbound data: ", err)
		return "", err
	}

	// Process raw inbound data
	resultInbound := processedAllRawData(rawInboundData)
	inboundTraffic := sumTraffic(resultInbound)

	// Process raw outbound data
	resultOutbound := processedAllRawData(rawOutboundData)
	outboundTraffic := sumTraffic(resultOutbound)

	networkTraffic := NetworkTrafficALL{
		InboundTraffic:  inboundTraffic,
		OutboundTraffic: outboundTraffic,
	}
	jsondata, err := json.Marshal(networkTraffic)
	return string(jsondata), nil
}

func GetAllNetworkMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/EC2 from %v to %v", elementType, startTime, endTime)

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
						Namespace:  aws.String("AWS/EC2"),
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

func processedAllRawData(result *cloudwatch.GetMetricDataOutput) []struct {
	Timestamp time.Time
	Value     float64
} {
	var processedData []struct {
		Timestamp time.Time
		Value     float64
	}

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		// Convert bytes to megabytes
		valueMB := value / (1024 * 1024)
		processedData = append(processedData, struct {
			Timestamp time.Time
			Value     float64
		}{Timestamp: *timestamp, Value: valueMB})
	}

	return processedData
}

func sumTraffic(data []struct {
	Timestamp time.Time
	Value     float64
}) float64 {
	var totalTraffic float64
	for _, entry := range data {
		totalTraffic += entry.Value
	}
	return totalTraffic
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEC2NetworkTrafficCmdAllinstances)
}
