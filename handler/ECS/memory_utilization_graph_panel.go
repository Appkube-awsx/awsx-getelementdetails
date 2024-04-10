package ECS

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/Appkube-awsx/awsx-common/authenticate"
    "github.com/Appkube-awsx/awsx-common/awsclient"
    "github.com/Appkube-awsx/awsx-common/cmdb"
    "github.com/Appkube-awsx/awsx-common/config"
    "github.com/Appkube-awsx/awsx-common/model"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "github.com/spf13/cobra"
)

type MemoryUtilizationGraphResult struct {
    RawData []struct {
        Timestamp time.Time
        Value     float64
    } `json:"RawData"`
}

var AwsxECSMemoryUtilizationGraphCmd = &cobra.Command{
    Use:   "memory_utilization_graph_panel",
    Short: "get memory utilization graph metrics data",
    Long:  `command to get memory utilization graph metrics data`,

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
            jsonResp, cloudwatchMetricResp, err := GetMemoryUtilizationGraphData(cmd, clientAuth, nil)
            if err != nil {
                log.Println("Error getting memory utilization graph data : ", err)
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

func GetMemoryUtilizationGraphData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
    elementId, _ := cmd.PersistentFlags().GetString("elementId")
    cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
    instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
    elementType, _ := cmd.PersistentFlags().GetString("elementType")
    startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
    endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

    if elementId != "" {
        log.Println("getting cloud-element data from cmdb")
        apiUrl := cmdbApiUrl
        if cmdbApiUrl == "" {
            log.Println("using default cmdb url")
            apiUrl = config.CmdbUrl
        }
        log.Println("cmdb url: " + apiUrl)
        cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
        if err != nil {
            return "", nil, err
        }
        instanceId = cmdbData.InstanceId

    }
    fmt.Println("instanceId", instanceId)

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

    cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

    // Fetch raw data
    rawData, err := GetMemoryUtilizationGraphMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
    if err != nil {
        log.Println("Error in getting raw data: ", err)
        return "", nil, err
    }
    cloudwatchMetricData["Memory Utilization (GB)"] = rawData

    result := processMemoryUtilizationGraphRawData(rawData)

    jsonString, err := json.Marshal(result)
    if err != nil {
        log.Println("Error in marshalling json in string: ", err)
        return "", nil, err
    }

    return string(jsonString), cloudwatchMetricData, nil
}

func GetMemoryUtilizationGraphMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {

    elmType := "ECS/ContainerInsights"
    input := &cloudwatch.GetMetricDataInput{
        EndTime:   endTime,
        StartTime: startTime,
        MetricDataQueries: []*cloudwatch.MetricDataQuery{
            {
                Id: aws.String("m1"),
                MetricStat: &cloudwatch.MetricStat{
                    Metric: &cloudwatch.Metric{
                        Dimensions: []*cloudwatch.Dimension{
                            {
                                Name:  aws.String("ClusterName"),
                                Value: aws.String(instanceId),
                            },
                        },
                        MetricName: aws.String("MemoryUtilized"),
                        Namespace:  aws.String(elmType),
                    },
                    Period: aws.Int64(900),
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

func processMemoryUtilizationGraphRawData(result *cloudwatch.GetMetricDataOutput) MemoryUtilizationGraphResult {
    var rawData MemoryUtilizationGraphResult
    rawData.RawData = make([]struct {
        Timestamp time.Time
        Value     float64
    }, len(result.MetricDataResults[0].Timestamps))

    // Convert memory utilized from MB to GB
    for i, timestamp := range result.MetricDataResults[0].Timestamps {
        rawData.RawData[i].Timestamp = *timestamp
        // Convert MB to GB
        rawData.RawData[i].Value = (*result.MetricDataResults[0].Values[i]) / 1024.0
    }

    return rawData
}

func init() {
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("elementId", "", "element id")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("elementType", "", "element type")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("query", "", "query")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("zone", "", "aws region")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("startTime", "", "start time")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("endTime", "", "endcl time")
    AwsxECSMemoryUtilizationGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
