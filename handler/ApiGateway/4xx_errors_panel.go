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

type Api4xxResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"4xx Errors"`
}

var AwsxApi4xxErrorCmd = &cobra.Command{
	Use:   "api_4xxerror_panel",
	Short: "get 4xxerror metrics data",
	Long:  `command to get 4xxerror metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApi4xxErrorData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API 4xx error data: ", err)
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

func GetApi4xxErrorData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := GetApi4xxErrorMetricValue(clientAuth, ApiName, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting 4xx error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["4XXError"] = metricValue

	result := process4xxErrorRawData(metricValue)


	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData,nil
}

func GetApi4xxErrorMetricValue(clientAuth *model.Auth, ApiName string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("error_4xx"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/ApiGateway"),
						MetricName: aws.String("4XXError"),
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
		return nil, err
	}

	return result, nil
}
func process4xxErrorRawData(result *cloudwatch.GetMetricDataOutput) Api4xxResult {
	var rawData Api4xxResult
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
	AwsxApi4xxErrorCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApi4xxErrorCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApi4xxErrorCmd.PersistentFlags().String("query", "", "query")
	AwsxApi4xxErrorCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApi4xxErrorCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApi4xxErrorCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApi4xxErrorCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApi4xxErrorCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApi4xxErrorCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApi4xxErrorCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApi4xxErrorCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApi4xxErrorCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApi4xxErrorCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApi4xxErrorCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApi4xxErrorCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApi4xxErrorCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApi4xxErrorCmd.PersistentFlags().String("ApiName", "", "api name")
}
