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

type ExecutionTimeData struct {
	FunctionName string
	ResponseTime float64
	Duration     float64
}

var LambdaExecutionTimeCmd = &cobra.Command{
	Use:   "execution_time_panel",
	Short: "Get Lambda function execution time metrics",
	Long:  `Command to get Lambda function execution time metrics`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running Lambda execution time metrics command")

		responseType, _ := cmd.PersistentFlags().GetString("responseType")

		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Fatalf("Error getting client auth: %v", err)
		}
		if !authFlag {
			log.Println("Authentication failed")
			return
		}

		jsonResp, cloudwatchMetricResp, err := GetLambdaExecutionTimePanel(cmd, clientAuth)
		if err != nil {
			log.Fatalf("Error getting Lambda execution time metric data: %v", err)
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaExecutionTimePanel(cmd *cobra.Command, clientAuth *model.Auth) (string, *map[string]*cloudwatch.GetMetricDataOutput, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	functionName := "List-Org-Github"
	fmt.Println("getting function", functionName)

	startTime, endTime, err := parseTimeFlag(startTimeStr, endTimeStr)
	if err != nil {
		return "", nil, err
	}

	cloudwatchMetricData := make(map[string]*cloudwatch.GetMetricDataOutput)

	rawData, err := GetLambdaExecutionTimeMetricData(clientAuth, startTime, endTime, functionName)
	if err != nil {
		return "", nil, err
	}
	cloudwatchMetricData[functionName] = rawData

	processedDataList := processExecutionTimeRawData(rawData, functionName)

	jsonString, err := json.Marshal(processedDataList)
	if err != nil {
		log.Fatalf("Error marshalling JSON response: %v", err)
	}
	return string(jsonString), &cloudwatchMetricData, nil
}

func GetLambdaExecutionTimeMetricData(clientAuth *model.Auth, startTime, endTime time.Time, functionName string) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		StartTime: &startTime,
		EndTime:   &endTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("Duration"),
						Namespace:  aws.String("AWS/Lambda"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("FunctionName"),
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
						MetricName: aws.String("Duration"),
						Namespace:  aws.String("AWS/Lambda"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("FunctionName"),
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

func processExecutionTimeRawData(result *cloudwatch.GetMetricDataOutput, functionName string) []*ExecutionTimeData {
	var executionTimeDataList []*ExecutionTimeData

	if len(result.MetricDataResults) == 0 {
		return executionTimeDataList
	}

	numDataPoints := len(result.MetricDataResults[0].Timestamps)
	if numDataPoints == 0 {
		return executionTimeDataList
	}

	for i := 0; i < numDataPoints; i++ {
		executionTimeData := &ExecutionTimeData{}

		duration := *result.MetricDataResults[0].Values[i]
		// initDuration := *result.MetricDataResults[1].Values[i]

		executionTimeData.FunctionName = functionName
		executionTimeData.ResponseTime = duration 
		executionTimeData.Duration = duration

		executionTimeDataList = append(executionTimeDataList, executionTimeData)
	}

	return executionTimeDataList
}

func parseTimeFlag(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
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
	LambdaExecutionTimeCmd.PersistentFlags().String("startTime", "", "Start time")
	LambdaExecutionTimeCmd.PersistentFlags().String("endTime", "", "End time")
	LambdaExecutionTimeCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
