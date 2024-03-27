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

type ThrottleData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxLambdaThrottleCmd = &cobra.Command{
	Use:   "throttle_panel",
	Short: "get lambda throttle metrics data",
	Long:  `Command to get lambda throttle metrics data`,

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

		jsonResp, cloudwatchMetricResp, err := GetLambdaThrottleData(cmd, clientAuth, nil)
		if err != nil {
			log.Println("Error getting lambda throttle data: ", err)
			return
		}

		if responseType == "frame" {
			fmt.Println(cloudwatchMetricResp)
		} else {
			fmt.Println(jsonResp)
		}
	},
}

func GetLambdaThrottleData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	cloudwatchMetricData := make(map[string]*cloudwatch.GetMetricDataOutput)

	// Fetch raw data
	rawData, err := GetLambdaThrottleMetricData(clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	cloudwatchMetricData["Throttling"] = rawData

	// Generate JSON response
	jsonString, err := json.Marshal(rawData)
	if err != nil {
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, err
}


func GetLambdaThrottleMetricData(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("throttle_query"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("Throttles"),
						Namespace:  aws.String("AWS/Lambda"),
					},
					Period: aws.Int64(300), // 5 minutes
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

	return result, nil
}

func init() {
	AwsxLambdaThrottleCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaThrottleCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaThrottleCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
