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

type MemoryData struct {
	FunctionName string
	MemoryUnit   float64
}

var AwsxLambdaMaxMemoryCmd = &cobra.Command{
	Use:   "max_memory_used_panel",
	Short: "get lambda memory metrics data",
	Long:  `Command to get lambda memory metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running from child command")
		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			return
		}
		if !authFlag {
			log.Println("Authentication failed")
			return
		}

		responseType, _ := cmd.PersistentFlags().GetString("responseType")
		if responseType != "json" && responseType != "frame" {
			log.Println("Invalid response type. Valid options are 'json' or 'frame'.")
			return
		}

		jsonResp, cloudwatchMetricResp, err := GetLambdaMaxMemoryData(cmd, clientAuth, nil)
		if err != nil {
			log.Println("Error getting lambda max memory used data: ", err)
			return
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaMaxMemoryData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, *map[string]*cloudwatch.GetMetricDataOutput, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
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
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	functionName := "List-Org-Github"
	fmt.Println("getting function", functionName)
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := GetLambdaMaxMemoryMetricData(clientAuth, startTime, endTime, cloudWatchClient, functionName)
	if err != nil {
		log.Printf("Error in getting lambda memory metric data for function %s: %v\n", functionName, err)
		return "", nil, err
	}
	cloudwatchMetricData["Max Memory Used (MB)"] = rawData
	result := processMaxMemoryRawData(rawData, functionName)
	// result.FunctionName = functionName

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), &cloudwatchMetricData, nil
}

func GetLambdaMaxMemoryMetricData(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch, functionName string) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("used_memory_max"),
						Namespace:  aws.String("LambdaInsights"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("function_name"),
								Value: aws.String(functionName),
							},
						},
					},
					Period: aws.Int64(300), // 5 minutes
					Stat:   aws.String("Maximum"),
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

func processMaxMemoryRawData(result *cloudwatch.GetMetricDataOutput, functionName string) []*MemoryData {
	var memoryDataArray []*MemoryData

	numDataPoints := len(result.MetricDataResults[0].Timestamps)
	fmt.Println("Number of data points:", numDataPoints)

	if numDataPoints > 0 {
		for i := 0; i < numDataPoints; i++ {
			memoryData := &MemoryData{
				FunctionName: functionName,
				MemoryUnit:   *result.MetricDataResults[0].Values[i],
			}
			memoryDataArray = append(memoryDataArray, memoryData)
		}
	}

	return memoryDataArray
}

func init() {
	AwsxLambdaMaxMemoryCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaMaxMemoryCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaMaxMemoryCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
