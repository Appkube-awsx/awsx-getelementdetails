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

type SuccessFailureResult struct {
	SuccessRate float64 `json:"SuccessRate"`
	FailureRate float64 `json:"FailureRate"`
}

var AwsxLambdaSuccessFailureCmd = &cobra.Command{
	Use:   "successfailure_panel",
	Short: "get successfailure metrics data",
	Long:  `command to get successfailure metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaSuccessFailureData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda throttling data : ", err)
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

func GetLambdaSuccessFailureData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, *SuccessFailureResult, error) {
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

	// Fetch raw data
	invocations, err := GetMetricData(cloudWatchClient, clientAuth, "Invocations", startTime, endTime)
	if err != nil {
		log.Println("Error retrieving invocations:", err)
		return "", nil, err
	}

	errors, err := GetMetricData(cloudWatchClient, clientAuth, "Errors", startTime, endTime)
	if err != nil {
		log.Println("Error retrieving errors:", err)
		return "", nil, err
	}

	// Calculate success rate
	successRate := (invocations - errors) / invocations * 100

	// Calculate failure rate
	failureRate := (errors / invocations) * 100

	result := &SuccessFailureResult{
		SuccessRate: successRate,
		FailureRate: failureRate,
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		
		return "", nil, err
	}

	return string(jsonString), result, nil
}

func GetMetricData(cloudWatchClient *cloudwatch.CloudWatch, clientAuth *model.Auth, metricName string, startTime, endTime *time.Time) (float64, error) {
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String(metricName),
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300),               // Period of 5 minutes
		Statistics: []*string{aws.String("Sum")}, // Sum the metric over the specified period
	}

	result, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}

	if len(result.Datapoints) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// Extract the sum of the metric from the first datapoint
	metricValue := aws.Float64Value(result.Datapoints[0].Sum)

	return metricValue, nil
}

func init() {
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaSuccessFailureCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
