package ApiGateway

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"

	// "github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
)

type APIGatewayLatency struct {
	Latency []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Response Time"`
}

var ApiResponseTimeCmd = &cobra.Command{
	Use:   "api_response_time_panel",
	Short: "Get API response time metrics data",
	Long:  `Command to get API response time metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command...")
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
			jsonResp, cloudwatchMetricResp, err := GetApiResponseTimePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API response time metrics: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case. It prints JSON
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetApiResponseTimePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	apiName := "dev-appkube-ecommerce-api"
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute) // Default: 5 minutes ago
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now() // Default: Current time
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := GetApiGatewayLatencyMetricData(clientAuth, apiName, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting API response time data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Response Time"] = rawData

	result := processTheRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling JSON to string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApiGatewayLatencyMetricData(clientAuth *model.Auth, apiName string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for API %s latency from %v to %v", apiName, startTime, endTime)

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
								Name:  aws.String("ApiName"),
								Value: aws.String(apiName),
							},
						},
						MetricName: aws.String("Latency"),
						Namespace:  aws.String("AWS/ApiGateway"),
					},
					Period: aws.Int64(60),         // Period of data in seconds
					Stat:   aws.String("Average"), // Average latency over the period
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

func processTheRawData(result *cloudwatch.GetMetricDataOutput) APIGatewayLatency {
	var rawData APIGatewayLatency
	rawData.Latency = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.Latency[i].Timestamp = *timestamp
		rawData.Latency[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	ApiResponseTimeCmd.PersistentFlags().String("startTime", "", "Start Time in RFC3339 format")
	ApiResponseTimeCmd.PersistentFlags().String("endTime", "", "End Time in RFC3339 format")
	ApiResponseTimeCmd.PersistentFlags().String("responseType", "", "Response type: json/frame")
}
