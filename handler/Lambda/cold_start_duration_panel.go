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

type ColdStartData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxLambdaColdStartCmd = &cobra.Command{
	Use:   "cold_start_duration_panel",
	Short: "get lambda cold start duration metrics data",
	Long:  `Command to get lambda cold start duration metrics data`,

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

		jsonResp, cloudwatchMetricResp, err := GetLambdaColdStartData(cmd, clientAuth, nil)
		if err != nil {
			log.Println("Error getting lambda cold start duration data: ", err)
			return
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaColdStartData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]interface{}, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	functionName := "List-Org-Github"

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

	cloudwatchMetricData := make(map[string]interface{})

	// Fetch raw data
	rawData, err := GetLambdaColdStartMetricData(clientAuth, startTime, endTime,functionName, cloudWatchClient)
	if err != nil {
		return "", nil, err
	}

	// Check if there are data points available
	if len(rawData.MetricDataResults) > 0 && len(rawData.MetricDataResults[0].Values) > 0 {
		// Check if there's only one data point
		if len(rawData.MetricDataResults[0].Values) == 1 {
			cloudwatchMetricData["Cold Start Duration"] = rawData.MetricDataResults[0].Values[0]
		} else {
			// If there are multiple data points, store the entire output
			cloudwatchMetricData["Cold Start Duration"] = rawData
		}
	} else {
		// No data points available
		cloudwatchMetricData["Cold Start Duration"] = "No data available for the specified time range and metric"
	}

	// Generate JSON response
	jsonString, err := json.Marshal(cloudwatchMetricData)
	if err != nil {
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaColdStartMetricData(clientAuth *model.Auth, startTime, endTime *time.Time,functionName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cold_start_duration"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("init_duration"),
						Namespace:  aws.String("LambdaInsights"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("function_name"),
								Value: aws.String(functionName),
							},
						},
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
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
	fmt.Println("result", result)
	return result, nil
}

func init() {
	AwsxLambdaColdStartCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaColdStartCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaColdStartCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
