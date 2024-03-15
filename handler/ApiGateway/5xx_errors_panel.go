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

type Api5xxResult struct {
	Value float64 `json:"Value"`
}

var AwsxApi5xxErrorCmd = &cobra.Command{
	Use:   "api_5xxerror_panel",
	Short: "get 5xxerror metrics data",
	Long:  `command to get 5xxerror metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApi5xxErrorData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API 5xx error data: ", err)
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

func GetApi5xxErrorData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
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

	cloudwatchMetricData := map[string]float64{}

	// Fetch raw data
	metricValue, err := GetApi5xxErrorMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting 5xx error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["5XXError"] = metricValue

	// Debug prints
	log.Printf("5XXError Metric Value: %f", metricValue)

	jsonString, err := json.Marshal(MetricResult{Value: metricValue})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApi5xxErrorMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("error_5xx"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/ApiGateway"),
						MetricName: aws.String("5XXError"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ApiName"),
								Value: aws.String(ApiName),
							},
						},
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Sum"), // Use Sum statistic to get total count
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
		return 0, err
	}

	if len(result.MetricDataResults) == 0 || len(result.MetricDataResults[0].Values) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// If there is only one value, return it
	if len(result.MetricDataResults[0].Values) == 1 {
		return aws.Float64Value(result.MetricDataResults[0].Values[0]), nil
	}

	// If there are multiple values, calculate the sum
	var sum float64
	for _, v := range result.MetricDataResults[0].Values {
		sum += aws.Float64Value(v)
	}
	return sum, nil
}

func init() {
	AwsxApi5xxErrorCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApi5xxErrorCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApi5xxErrorCmd.PersistentFlags().String("query", "", "query")
	AwsxApi5xxErrorCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApi5xxErrorCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApi5xxErrorCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApi5xxErrorCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApi5xxErrorCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApi5xxErrorCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApi5xxErrorCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApi5xxErrorCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApi5xxErrorCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApi5xxErrorCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApi5xxErrorCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApi5xxErrorCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApi5xxErrorCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApi5xxErrorCmd.PersistentFlags().String("ApiName", "", "api name")
}
