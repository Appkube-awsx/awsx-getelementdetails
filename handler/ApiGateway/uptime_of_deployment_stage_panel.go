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
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type MetricResultss struct {
	UptimePercentage   float64 `json:"uptimePercentage"`
	DowntimePercentage float64 `json:"downtimePercentage"`
}

var AwsxApiDeploymentCmd = &cobra.Command{
	Use:   "api_downtime_deployment_panel",
	Short: "Get uptime and downtime deployment metrics data for API stages",
	Long:  `Command to get uptime and downtime deployment metrics data for API stages`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, err := GetApiUptimedata(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting API uptime data: ", err)
				return
			}

			if responseType == "frame" {
				fmt.Println(jsonResp)
			} else {
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetApiUptimedata(cmd *cobra.Command, clientAuth *model.Auth) (string, error) {
	apiID := "i3mdnxvgrf"
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime, err := parseTime(startTimeStr, endTimeStr)
	if err != nil {
		return "", err
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	stages, err := GetStagesForAPI(clientAuth, apiID)
	if err != nil {
		return "", err
	}

	cloudwatchMetricData := make(map[string]MetricResultss)

	for _, stage := range stages {
		log.Printf("Fetching metrics for stage: %s", stage)
		totalRequests, err := GetMetricValue(clientAuth, startTime, endTime, apiID, stage, "Count", "Sum")
		if err != nil {
			log.Printf("Error in getting total requests metric value for stage %s: %v", stage, err)
			return "", err
		}

		clientErrors, err := GetMetricValue(clientAuth, startTime, endTime, apiID, stage, "4XXError", "Sum")
		if err != nil {
			log.Printf("Error in getting client errors metric value for stage %s: %v", stage, err)
			return "", err
		}

		serverErrors, err := GetMetricValue(clientAuth, startTime, endTime, apiID, stage, "5XXError", "Sum")
		if err != nil {
			log.Printf("Error in getting server errors metric value for stage %s: %v", stage, err)
			return "", err
		}

		uptimePercentage := calculateUptimePercentage(totalRequests, clientErrors, serverErrors)
		downtimePercentage := 100 - uptimePercentage

		cloudwatchMetricData[stage] = MetricResultss{
			UptimePercentage:   uptimePercentage,
			DowntimePercentage: downtimePercentage,
		}

		// log.Printf("Uptime Percentage for stage %s: %f", stage, uptimePercentage)
		// log.Printf("Downtime Percentage for stage %s: %f", stage, downtimePercentage)
	}

	jsonString, err := json.Marshal(cloudwatchMetricData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", err
	}

	return string(jsonString), nil
}

func parseTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing start time: %v", err)
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing end time: %v", err)
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	return startTime, endTime, nil
}

func GetStagesForAPI(clientAuth *model.Auth, apiID string) ([]string, error) {
	apiGatewayClient := awsclient.GetClient(*clientAuth, awsclient.APIGATEWAY_CLIENT).(*apigateway.APIGateway)

	params := &apigateway.GetStagesInput{
		RestApiId: aws.String(apiID),
	}

	resp, err := apiGatewayClient.GetStages(params)
	if err != nil {
		return nil, err
	}

	stages := make([]string, len(resp.Item))
	for i, stage := range resp.Item {
		stages[i] = aws.StringValue(stage.StageName)
	}

	return stages, nil
}

func GetMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, apiID, stage, metricName, statistic string) (float64, error) {
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	apiName := "dev-hrms"
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ApiGateway"),
		MetricName: aws.String(metricName),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("ApiName"),
				Value: aws.String(apiName),
			},
			{
				Name:  aws.String("Stage"),
				Value: aws.String(stage),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String(statistic)},
		Unit:       aws.String("Count"),
	}
	// fmt.Println("Input for fetching metrics:", input)
	result, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}

	if len(result.Datapoints) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// Assuming Sum statistic for simplicity, you can modify this based on your requirements
	sum := 0.0
	for _, datapoint := range result.Datapoints {
		sum += aws.Float64Value(datapoint.Sum)
	}

	return sum, nil
}

func calculateUptimePercentage(totalRequests, clientErrors, serverErrors float64) float64 {
	total := totalRequests + clientErrors + serverErrors
	if total == 0 {
		return 100
	}
	uptimePercentage := (totalRequests / total) * 100
	fmt.Println("Uptime percentage:", uptimePercentage)
	return uptimePercentage
}

func init() {
	AwsxApiDeploymentCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxApiDeploymentCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxApiDeploymentCmd.PersistentFlags().String("responseType", "", "Response type: json/frame")
	AwsxApiDeploymentCmd.PersistentFlags().String("ApiName", "", "API Name")
}
