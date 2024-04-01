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

type ApiSuccessfulFailedResult struct {
	SuccessfulEvents float64 `json:"successfulEvents"`
	FailedEvents float64 `json:"failedEvents"`
}

var AwsxApiSuccessfulFailedCmd = &cobra.Command{
	Use:   "successful_and_failed_events_panel",
	Short: "get successful failed metrics data",
	Long:  `command to get successful failed metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiSuccessFailedData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API error data: ", err)
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

func GetApiSuccessFailedData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
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
	totalevents, err := GetApiTotalEventsMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting  error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Count"] = totalevents

	error_4xx, err := GetApiClientErrorMetricValue(clientAuth, startTime, endTime, ApiName,  cloudWatchClient)
	if err != nil {
		log.Println("Error in getting  error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["4xxError"] = error_4xx

	error_5xx, err := GetApiServerErrorsMetricValue(clientAuth, startTime, endTime, ApiName , cloudWatchClient)
	if err != nil {
		log.Println("Error in getting  error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["5xxError"] = error_5xx

	SuccessEvents := totalevents - error_4xx - error_5xx

	failedEvents := error_4xx - error_5xx

	jsonString, err := json.Marshal(ApiSuccessfulFailedResult{SuccessfulEvents: SuccessEvents, FailedEvents :failedEvents})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetApiTotalEventsMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time,ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
    input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ApiGateway"),
		MetricName: aws.String("Count"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("ApiName"),
				Value: aws.String(ApiName),
			},
		},
        StartTime: startTime,
        EndTime:   endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Sum")},
		Unit:       aws.String("Count"),
    }

    if cloudWatchClient == nil {
        cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
    }

    result, err := cloudWatchClient.GetMetricStatistics(input)
    if err != nil {
        return 0, err
    }

    if len(result.Datapoints) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// Assuming there's only one datapoint, return its Sum
	return aws.Float64Value(result.Datapoints[0].Sum), nil
}

func GetApiClientErrorMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ApiGateway"),
		MetricName: aws.String("4XXError"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("ApiName"),
				Value: aws.String(ApiName),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Sum")},
		Unit:       aws.String("Count"),
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}

	if len(result.Datapoints) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// Assuming there's only one datapoint, return its Sum
	return aws.Float64Value(result.Datapoints[0].Sum), nil
}

func GetApiServerErrorsMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ApiGateway"),
		MetricName: aws.String("5XXError"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("ApiName"),
				Value: aws.String(ApiName),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Sum")},
		Unit:       aws.String("Count"),
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}

	if len(result.Datapoints) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// Assuming there's only one datapoint, return its Sum
	return aws.Float64Value(result.Datapoints[0].Sum), nil
}

func init() {
	AwsxApiSuccessfulFailedCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiSuccessfulFailedCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiSuccessfulFailedCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiSuccessfulFailedCmd.PersistentFlags().String("ApiName", "", "api name")
}
