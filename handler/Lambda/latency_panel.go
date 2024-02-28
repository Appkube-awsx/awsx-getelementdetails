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

type LatencyResult struct {
    Value float64 `json:"Value"`
}

var AwsxLambdaLatencyCmd = &cobra.Command{
    Use:   "latency_panel",
    Short: "get latency metrics data",
    Long:  `command to get latency metrics data`,

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
            jsonResp, cloudwatchMetricResp, err := GetLambdaLatencyData(cmd, clientAuth, nil)
            if err != nil {
                log.Println("Error getting lambda latency data : ", err)
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

func GetLambdaLatencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]float64, error) {
    startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
    endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

    var startTime, endTime *time.Time

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
    avgLatencyValue, err := GetAverageLambdaLatencyMetricValue(clientAuth, startTime, endTime, cloudWatchClient)
    if err != nil {
        log.Println("Error in getting average latency value: ", err)
        return "", nil, err
    }
    cloudwatchMetricData["AverageLatency"] = avgLatencyValue

    // Debug prints
    log.Printf("Average Latency Value: %f", avgLatencyValue)

    jsonString, err := json.Marshal(LatencyResult{Value: avgLatencyValue})
    if err != nil {
        log.Println("Error in marshalling json in string: ", err)
        return "", nil, err
    }

    return string(jsonString), cloudwatchMetricData, nil
}

func GetAverageLambdaLatencyMetricValue(clientAuth *model.Auth, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
    input := &cloudwatch.GetMetricStatisticsInput{
        Namespace:  aws.String("AWS/Lambda"),
        MetricName: aws.String("Duration"),
        StartTime:  startTime,
        EndTime:    endTime,
        Period:     aws.Int64(300), // Adjust period as needed (e.g., 5 minutes)
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

    // Extract the average value from the first datapoint
    averageLatencyValue := aws.Float64Value(result.Datapoints[0].Average)

    return averageLatencyValue, nil
}

func init() {
    AwsxLambdaLatencyCmd.PersistentFlags().String("elementId", "", "element id")
    AwsxLambdaLatencyCmd.PersistentFlags().String("elementType", "", "element type")
    AwsxLambdaLatencyCmd.PersistentFlags().String("query", "", "query")
    AwsxLambdaLatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
    AwsxLambdaLatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
    AwsxLambdaLatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
    AwsxLambdaLatencyCmd.PersistentFlags().String("zone", "", "aws region")
    AwsxLambdaLatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
    AwsxLambdaLatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
    AwsxLambdaLatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
    AwsxLambdaLatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
    AwsxLambdaLatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
    AwsxLambdaLatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
    AwsxLambdaLatencyCmd.PersistentFlags().String("startTime", "", "start time")
    AwsxLambdaLatencyCmd.PersistentFlags().String("endTime", "", "end time")
    AwsxLambdaLatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
