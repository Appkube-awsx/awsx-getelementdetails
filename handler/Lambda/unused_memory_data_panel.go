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

type UnusedMemoryData struct {
	FunctionName             string
	AvgUnusedAllocatedMemory float64
	MaxMemoryUsedAvg         float64
}

var LambdaMemoryMetricsCmd = &cobra.Command{
	Use:   "unused_memory_data_panel",
	Short: "Get Lambda unused memory data",
	Long:  `Command to get Lambda Unused memory data`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running Lambda memory metrics command")

		responseType, _ := cmd.PersistentFlags().GetString("responseType")

		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Fatalf("Error getting client auth: %v", err)
		}
		if !authFlag {
			log.Println("Authentication failed")
			return
		}

		jsonResp, cloudwatchMetricResp, err := GetLambdaUnusedMemoryPanel(cmd, clientAuth)
		if err != nil {
			log.Fatalf("Error getting Lambda memory metric data: %v", err)
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaUnusedMemoryPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, *map[string]*cloudwatch.GetMetricDataOutput, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	functionName := "List-Org-Github"
	fmt.Println("getting function", functionName)

	startTime, endTime, err := parseTimeFlags(startTimeStr, endTimeStr)
	if err != nil {
		return "", nil, err
	}

	cloudwatchMetricData := make(map[string]*cloudwatch.GetMetricDataOutput)

	rawData, err := GetLambdaMemoryMetricData(clientAuth, startTime, endTime, functionName)
	if err != nil {
		return "", nil, err
	}
	cloudwatchMetricData[functionName] = rawData

	processedDataList := processMemoryRawData(rawData, functionName)

	jsonString, err := json.Marshal(processedDataList)
	if err != nil {
		log.Fatalf("Error marshalling JSON response: %v", err)
	}
	return string(jsonString), &cloudwatchMetricData, nil
}

func GetLambdaMemoryMetricData(clientAuth *model.Auth, startTime, endTime time.Time, functionName string) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		StartTime: &startTime,
		EndTime:   &endTime,
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
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("m2"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("total_memory"),
						Namespace:  aws.String("LambdaInsights"),
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
			},
		},
	}

	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processMemoryRawData(result *cloudwatch.GetMetricDataOutput, functionName string) []*UnusedMemoryData {
	var memoryDataList []*UnusedMemoryData

	if len(result.MetricDataResults) == 0 {
		return memoryDataList
	}

	numDataPoints := len(result.MetricDataResults[0].Timestamps)
	if numDataPoints == 0 {
		return memoryDataList
	}

	for i := 0; i < numDataPoints; i++ {
		memoryData := &UnusedMemoryData{}

		usedMemory := *result.MetricDataResults[0].Values[i]
		totalMemory := *result.MetricDataResults[1].Values[i]

		memoryData.FunctionName = functionName // Set the FunctionName here
		memoryData.AvgUnusedAllocatedMemory = totalMemory - usedMemory
		memoryData.MaxMemoryUsedAvg = usedMemory

		memoryDataList = append(memoryDataList, memoryData)
	}

	return memoryDataList
}


func parseTimeFlags(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		startTime = time.Now().Add(-5 * time.Minute)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	} else {
		endTime = time.Now()
	}

	return startTime, endTime, nil
}

func init() {
	LambdaMemoryMetricsCmd.PersistentFlags().String("startTime", "", "Start time")
	LambdaMemoryMetricsCmd.PersistentFlags().String("endTime", "", "End time")
	LambdaMemoryMetricsCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
