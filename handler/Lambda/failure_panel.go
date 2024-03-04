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

type FailureResult struct {
	Value float64 `json:"Value"`
}

var AwsxLambdaFailureCmd = &cobra.Command{
	Use:   "failure_panel",
	Short: "get failure metrics data",
	Long:  `command to get failure metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaErrorData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda failure data : ", err)
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

func GetLambdaFailureData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {

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
	failureCount, err := GetLambdaFailureCount(clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting failure count: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["FailureCount"] = failureCount

	// Debug prints
	log.Printf("Failure Count: %f", failureCount)

	jsonString, err := json.Marshal(ErrorResult{Value: failureCount})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaFailureCount(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	// First, retrieve the total number of invocations
	totalInvocations, invocationsErr := GetTotalLambdaErrorInvocations(clientAuth, startTime, endTime, cloudWatchClient)
	if invocationsErr != nil {
		return 0, invocationsErr
	}
	log.Printf("Invocation Value: %f", totalInvocations)

	// Then, retrieve the number of errors
	totalErrors, err := GetTotalLambdaErrors(clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		return 0, err
	}
	log.Printf("Error Value: %f", totalErrors)
	// Calculate failure count: total invocations - total errors
	failureCount := totalInvocations - totalErrors
	log.Printf("failure Value: %f", failureCount)
	return failureCount, nil
}

func GetTotalLambdaErrorInvocations(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // Period of 5 minutes
		Statistics: []*string{aws.String("Sum")},
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

	// Extract the sum value from the first datapoint
	sumValue := aws.Float64Value(result.Datapoints[0].Sum)

	return sumValue, nil
}

func GetTotalLambdaErrors(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Errors"),
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // Period of 5 minutes
		Statistics: []*string{aws.String("Sum")},
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

	// Extract the sum value from the first datapoint
	sumValue := aws.Float64Value(result.Datapoints[0].Sum)

	return sumValue, nil
}

func init() {
	AwsxLambdaFailureCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaFailureCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaFailureCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaFailureCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaFailureCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaFailureCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaFailureCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaFailureCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaFailureCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaFailureCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaFailureCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaFailureCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaFailureCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaFailureCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaFailureCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaFailureCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
