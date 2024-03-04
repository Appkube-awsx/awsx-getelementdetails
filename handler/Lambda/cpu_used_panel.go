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

type CpuResult struct {
	Value float64 `json:"Value"`
}

var AwsxLambdaCpuCmd = &cobra.Command{
	Use:   "cpu_panel",
	Short: "get cpu metrics data",
	Long:  `command to get cpu metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaLatencyData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda cpu data : ", err)
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

func GetLambdaCpuData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
	functionName := "CW-agent-installation-automation"

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
	CpuUsedValue, err := GetLambdaCpuMetricValue(clientAuth, startTime, endTime, functionName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu used value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CpuUsedValue"] = CpuUsedValue

	// Debug prints
	log.Printf("Raw Cpu Value: %f", CpuUsedValue)

	jsonString, err := json.Marshal(CpuResult{Value: CpuUsedValue})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaCpuMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, functionName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cpu_total_time"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("LambdaInsights"),
						MetricName: aws.String("cpu_total_time"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("function_name"),
								Value: aws.String(functionName),
							},
						},
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
	AwsxLambdaCpuCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaCpuCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaCpuCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaCpuCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaCpuCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaCpuCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaCpuCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaCpuCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaCpuCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaCpuCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaCpuCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaCpuCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaCpuCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaCpuCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaCpuCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaCpuCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
