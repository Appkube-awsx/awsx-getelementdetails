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

type CacheMissResult struct {
	TimeSeries []struct {
		Timestamp time.Time
		Value     float64
	} `json:"timeSeries"`
}

var AwsxApiCacheMissCmd = &cobra.Command{
	Use:   "cache_miss_count_panel",
	Short: "get cache miss count metrics data",
	Long:  `command to get cache miss count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiCacheMissData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API cache miss count data: ", err)
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

func GetApiCacheMissData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := GetApiCacheMissMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting API cache miss count metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CacheMiss"] = metricValue

	result := processCacheMissRawData(metricValue)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApiCacheMissMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cacheMiss"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/ApiGateway"),
						MetricName: aws.String("CacheMissCount"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ApiName"),
								Value: aws.String(ApiName),
							},
						},
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Sum"), // Sum to get the total count
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

func processCacheMissRawData(result *cloudwatch.GetMetricDataOutput) CacheMissResult {
	var timeSeries []struct {
		Timestamp time.Time
		Value     float64
	}

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		timeSeries = append(timeSeries, struct {
			Timestamp time.Time
			Value     float64
		}{
			Timestamp: *timestamp,
			Value:     *result.MetricDataResults[0].Values[i],
		})
	}

	return CacheMissResult{TimeSeries: timeSeries}
}

func init() {
	AwsxApiCacheMissCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiCacheMissCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiCacheMissCmd.PersistentFlags().String("query", "", "query")
	AwsxApiCacheMissCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiCacheMissCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiCacheMissCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiCacheMissCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiCacheMissCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiCacheMissCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiCacheMissCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiCacheMissCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiCacheMissCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiCacheMissCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiCacheMissCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiCacheMissCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiCacheMissCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiCacheMissCmd.PersistentFlags().String("ApiName", "", "api name")
}
