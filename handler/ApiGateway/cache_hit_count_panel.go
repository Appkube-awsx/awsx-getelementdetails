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

type CacheHitsResult struct {
	TimeSeries []struct {
		Timestamp time.Time
		Value     float64
	} `json:"timeSeries"`
}

var AwsxApiCacheHitsCmd = &cobra.Command{
	Use:   "cache_hit_count_panel",
	Short: "get cache hits metrics data",
	Long:  `command to get cache hits metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiCacheHitsData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API cache hits data: ", err)
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

func GetApiCacheHitsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := GetApiCacheHitsMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting API cache hits metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CacheHits"] = metricValue

	result := processCacheHitsRawData(metricValue)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApiCacheHitsMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cacheHits"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/ApiGateway"),
						MetricName: aws.String("CacheHitCount"),
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

func processCacheHitsRawData(result *cloudwatch.GetMetricDataOutput) CacheHitsResult {
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

	return CacheHitsResult{TimeSeries: timeSeries}
}

func init() {
	AwsxApiCacheHitsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiCacheHitsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiCacheHitsCmd.PersistentFlags().String("query", "", "query")
	AwsxApiCacheHitsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiCacheHitsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiCacheHitsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiCacheHitsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiCacheHitsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiCacheHitsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiCacheHitsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiCacheHitsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiCacheHitsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiCacheHitsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiCacheHitsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiCacheHitsCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiCacheHitsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiCacheHitsCmd.PersistentFlags().String("ApiName", "", "api name")
}
