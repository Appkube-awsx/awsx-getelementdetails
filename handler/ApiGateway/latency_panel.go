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

type ApiLatency struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Latency "`
}

var AwsxApiLatencyCmd = &cobra.Command{
	Use:   "api_latency_panel",
	Short: "get latency metrics data",
	Long:  `command to get latency metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiLatencyData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API latency data: ", err)
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

func GetApiLatencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := GetApiLatencyMetricValue(clientAuth, ApiName, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting latency metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Latency"] = metricValue
	result := processLatencyRawData(metricValue)
	// Debug prints
	//log.Printf("Latency Metric Value: %f", metricValue)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApiLatencyMetricValue(clientAuth *model.Auth, ApiName string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("latency"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/ApiGateway"),
						MetricName: aws.String("Latency"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ApiName"),
								Value: aws.String(ApiName),
							},
						},
					},
					Period: aws.Int64(300), 
					Stat:   aws.String("Sum"),
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

func processLatencyRawData(result *cloudwatch.GetMetricDataOutput) ApiLatency {
	var rawData ApiLatency
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
	AwsxApiLatencyCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiLatencyCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiLatencyCmd.PersistentFlags().String("query", "", "query")
	AwsxApiLatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiLatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiLatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiLatencyCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiLatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiLatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiLatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiLatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiLatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiLatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiLatencyCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiLatencyCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiLatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiLatencyCmd.PersistentFlags().String("ApiName", "", "api name")
}
