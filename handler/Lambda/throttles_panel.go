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

type ThrottleResult struct {
    Value float64 `json:"Value"`
}

var AwsxLambdaThrottleCmd = &cobra.Command{
    Use:   "throttles_panel",
    Short: "get throttle metrics data",
    Long:  `command to get throttle metrics data`,

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
            jsonResp, cloudwatchMetricResp, err := GetLambdaThrottleData(cmd, clientAuth, nil)
            if err != nil {
                log.Println("Error getting lambda throttle data : ", err)
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

func GetLambdaThrottleData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
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
    maxThrottleValue, err := GetMaxLambdaThrottleMetricValue(clientAuth, startTime, endTime, cloudWatchClient)
    if err != nil {
        log.Println("Error in getting maximum throttle value: ", err)
        return "", nil, err
    }
    cloudwatchMetricData["MaxThrottles"] = maxThrottleValue

    // Debug prints
    log.Printf("Maximum Throttle Value: %f", maxThrottleValue)

    jsonString, err := json.Marshal(ThrottleResult{Value: maxThrottleValue})
    if err != nil {
        log.Println("Error in marshalling json in string: ", err)
        return "", nil, err
    }

    return string(jsonString), cloudwatchMetricData, nil
}

func GetMaxLambdaThrottleMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
    input := &cloudwatch.GetMetricStatisticsInput{
        Namespace:  aws.String("AWS/Lambda"),
        MetricName: aws.String("Throttles"),
        StartTime:  startTime,
        EndTime:    endTime,
        Period:     aws.Int64(2592000), 
        Statistics: []*string{aws.String("Average")},
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

    // Extract the maximum value from the first datapoint
    maxThrottleValue := aws.Float64Value(result.Datapoints[0].Maximum)

    return maxThrottleValue, nil
}

func init() {
    AwsxLambdaThrottleCmd.PersistentFlags().String("elementId", "", "element id")
    AwsxLambdaThrottleCmd.PersistentFlags().String("elementType", "", "element type")
    AwsxLambdaThrottleCmd.PersistentFlags().String("query", "", "query")
    AwsxLambdaThrottleCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
    AwsxLambdaThrottleCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
    AwsxLambdaThrottleCmd.PersistentFlags().String("vaultToken", "", "vault token")
    AwsxLambdaThrottleCmd.PersistentFlags().String("zone", "", "aws region")
    AwsxLambdaThrottleCmd.PersistentFlags().String("accessKey", "", "aws access key")
    AwsxLambdaThrottleCmd.PersistentFlags().String("secretKey", "", "aws secret key")
    AwsxLambdaThrottleCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
    AwsxLambdaThrottleCmd.PersistentFlags().String("externalId", "", "aws external id")
    AwsxLambdaThrottleCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
    AwsxLambdaThrottleCmd.PersistentFlags().String("instanceId", "", "instance id")
    AwsxLambdaThrottleCmd.PersistentFlags().String("startTime", "", "start time")
    AwsxLambdaThrottleCmd.PersistentFlags().String("endTime", "", "end time")
    AwsxLambdaThrottleCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
