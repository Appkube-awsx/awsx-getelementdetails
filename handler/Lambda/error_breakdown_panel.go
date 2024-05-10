package Lambda

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type BreakdownResult struct {
	ErrorPercentage  float64 `json:"ErrorPecentage"`
	Value            float64 `json:"Value"`
	PercentageChange float64 `json:"PercentageChange"`
	ChangeType       string  `json:"ChangeType"`
}

var AwsxLambdaErrorBreakdownCmd = &cobra.Command{
	Use:   "error_breakdown_panel",
	Short: "get error metrics data",
	Long:  `command to get error metrics data`,

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
		// if authFlag {
		// 	panel, err := GetLambdaTopErrorsMessagesEvents(cmd, clientAuth, nil)
		// 	if err != nil {
		// 		return
		// 	}
		// 	fmt.Println(panel)

		// }
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err := GetErrorBreakdownData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda error breakdown data : ", err)
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

func GetErrorBreakdownData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]interface{}, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	fmt.Println(instanceId)

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	// instanceId, err = comman_function.GetCmdbData(cmd)

	// if err != nil {
	// 	return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	// }
	cloudwatchMetricData := map[string]interface{}{}
	InvocationInput := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300),
		Statistics: []*string{aws.String("Sum")},
	}

	ErrorInput := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Errors"),
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300),
		Statistics: []*string{aws.String("Sum")},
	}

	// Fetch raw data for last month and current month
	lastMonthStartTime := startTime.AddDate(0, -1, 0)
	lastMonthEndTime := endTime.AddDate(0, -1, 0)

	InvocationInputLastMonth := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		StartTime:  &lastMonthStartTime,
		EndTime:    &lastMonthEndTime,
		Period:     aws.Int64(300),
		Statistics: []*string{aws.String("Sum")},
	}

	ErrorInputLastMonth := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Errors"),
		StartTime:  &lastMonthStartTime,
		EndTime:    &lastMonthEndTime,
		Period:     aws.Int64(300),
		Statistics: []*string{aws.String("Sum")},
	}
	fmt.Println("date", startTime, endTime, lastMonthStartTime, lastMonthEndTime)
	lastMonthInvocations, err := GetLambdaBreakdownData(InvocationInputLastMonth, clientAuth, &lastMonthStartTime, &lastMonthEndTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting error metric value for last month: ", err)
		return "", nil, err
	}

	currentMonthInvocations, err := GetLambdaBreakdownData(InvocationInput, clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting error metric value for current month: ", err)
		return "", nil, err
	}

	ErrorCount, err := GetLambdaBreakdownData(ErrorInput, clientAuth, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting error metric value for error: ", err)
		return "", nil, err
	}

	lastMonthErrorCount, err := GetLambdaBreakdownData(ErrorInputLastMonth, clientAuth, &lastMonthStartTime, &lastMonthEndTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting error metric value for last month: ", err)
		return "", nil, err
	}

	fmt.Println(lastMonthInvocations, currentMonthInvocations, ErrorCount, lastMonthErrorCount)
	// Calculate percentage change
	errorPercentage := (ErrorCount / currentMonthInvocations) * 100
	errorPercentageRounded := math.Round(errorPercentage*100) / 100

	percentageChange := ((ErrorCount - lastMonthErrorCount) / lastMonthErrorCount) * 100
	percentageChangeRounded := math.Round(percentageChange*100) / 100

	// Determine if it's an increment or decrement
	changeType := "increment"
	if percentageChange < 0 {
		changeType = "decrement"
	}

	cloudwatchMetricData["LastMonthErrorCount"] = lastMonthErrorCount
	cloudwatchMetricData["ErrorCount"] = ErrorCount
	cloudwatchMetricData["PercentageChange"] = fmt.Sprintf("%.2f%% %s", percentageChange, changeType)
	cloudwatchMetricData["ErrorPercentage"] = fmt.Sprintf("%.2f%%", errorPercentage)

	jsonString, err := json.Marshal(BreakdownResult{ErrorPercentage: errorPercentageRounded, Value: ErrorCount, PercentageChange: percentageChangeRounded, ChangeType: changeType})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaBreakdownData(input *cloudwatch.GetMetricStatisticsInput, clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}
	// fmt.Println(result)

	if len(result.Datapoints) == 0 {
		return 0, fmt.Errorf("no data available for the specified time range")
	}

	// Sum up the values from all the datapoints
	totalInvocations := 0.0
	for _, dp := range result.Datapoints {
		totalInvocations += aws.Float64Value(dp.Sum)
	}

	return totalInvocations, nil
}
func init() {
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaErrorBreakdownCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
