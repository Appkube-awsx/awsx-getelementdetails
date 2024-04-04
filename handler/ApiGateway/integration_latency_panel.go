package ApiGateway

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type ApiIntegrationLatencyResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"IntegrationLatency"`
}

var AwsxApiIntegrationLatencyCmd = &cobra.Command{
	Use:   "api_integration_latency_panel",
	Short: "get integration latency metrics data",
	Long:  `command to get integration latency metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiIntegrationLatencyData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API integration latency data: ", err)
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

func GetApiIntegrationLatencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	ApiName := "dev-hrms"
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", nil, err
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
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	// Debug prints
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	metricValue, err := GetApiIntegrationLatencyMetricValue(clientAuth, startTime, endTime, ApiName, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting latency metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["IntegrationLatency"] = metricValue

	result := processIntegrationLatencyRawData(metricValue)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApiIntegrationLatencyMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("integration_latency"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/ApiGateway"),
						MetricName: aws.String("IntegrationLatency"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ApiName"),
								Value: aws.String(ApiName),
							},
						},
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
				ReturnData: aws.Bool(true),
			},
		},
		StartTime: startTime,
		EndTime:   endTime,
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
func processIntegrationLatencyRawData(result *cloudwatch.GetMetricDataOutput) ApiIntegrationLatencyResult {
	var rawData ApiIntegrationLatencyResult
	rawData.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("query", "", "query")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("ApiName", "", "api name")
}
