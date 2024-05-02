package Lambda

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

type Request struct {
	Value float64 `json:"Value"`
}

var AwsxLambdaRequestCmd = &cobra.Command{
	Use:   "request_panel",
	Short: "get request metrics data",
	Long:  `command to get request metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaRequestData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda cpu received data : ", err)
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

func GetLambdaRequestData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
	//functionName := "CW-agent-installation-automation"

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
	rawData, err := GetLambdaRequestMetricValue(clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu used metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory"] = rawData

	// Debug prints
	log.Printf("Memory Metric Value: %f", rawData)

	jsonString, err := json.Marshal(rawData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaRequestMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/Lambda"),
						MetricName: aws.String("UrlRequestCount"),
						// Dimensions: []*cloudwatch.Dimension{
						// 	{
						// 		//Name:  aws.String("function_name"),
						// 		//Value: aws.String(functionName),
						// 	},
						// },
					},
					Period: aws.Int64(300), // 5 minutes
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
		return 0, err
	}

	if len(result.MetricDataResults) == 0 || len(result.MetricDataResults[0].Values) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// If there is only one value, return it
	if len(result.MetricDataResults[0].Values) == 1 {
		return aws.Float64Value(result.MetricDataResults[0].Values[0]), nil
	}

	// If there are multiple values, calculate the average
	var sum float64
	for _, v := range result.MetricDataResults[0].Values {
		sum += aws.Float64Value(v)
	}
	return sum / float64(len(result.MetricDataResults[0].Values)), nil
}

func init() {
	AwsxLambdaRequestCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaRequestCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaRequestCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaRequestCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaRequestCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaRequestCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaRequestCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaRequestCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaRequestCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaRequestCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaRequestCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaRequestCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaRequestCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaRequestCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaRequestCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaRequestCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaRequestCmd.PersistentFlags().String("functionName", "", "Lambda function name")
}
