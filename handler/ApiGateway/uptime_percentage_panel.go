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

type MetricResults struct {
	UptimePercentage float64 `json:"uptimePercentage"`
}

var AwsxApiUptimeCmd = &cobra.Command{
	Use:   "api_uptime_panel",
	Short: "get uptime metrics data",
	Long:  `command to get uptime metrics data`,

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
			jsonResp, uptimeMetricResp, err := GetApiUptimeData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API uptime data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(uptimeMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetApiUptimeData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
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
	totalRequests, err := GetTotalRequestsMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting total requests metric value: ", err)
		return "", nil, err
	}

	clientErrors, err := GetClientErrorsMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting client errors metric value: ", err)
		return "", nil, err
	}

	serverErrors, err := GetServerErrorsMetricValue(clientAuth, startTime, endTime, ApiName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting server errors metric value: ", err)
		return "", nil, err
	}

	// Calculate uptime percentage
	uptimePercentage := ((totalRequests - clientErrors - serverErrors) / totalRequests) * 100

	cloudwatchMetricData["UptimePercentage"] = uptimePercentage

	// Debug prints
	log.Printf("Uptime Percentage: %f", uptimePercentage)

	jsonString, err := json.Marshal(MetricResults{UptimePercentage: uptimePercentage})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetTotalRequestsMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ApiGateway"),
		MetricName: aws.String("Count"),
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

func GetClientErrorsMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
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

func GetServerErrorsMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, ApiName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
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
	AwsxApiUptimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiUptimeCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiUptimeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiUptimeCmd.PersistentFlags().String("ApiName", "", "api name")
}
