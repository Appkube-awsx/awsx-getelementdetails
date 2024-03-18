package EKS

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

type DiskIOPerformanceResult struct {
    TotalOps []struct {
        Timestamp time.Time
        Value     float64
    } `json:"total_ops"`
}

var AwsxEKSDiskIOPerformanceCmd = &cobra.Command{
    Use:   "disk_io_performance_panel",
    Short: "get disk I/O performance metrics data",
    Long:  `command to get disk I/O performance metrics data`,

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
            jsonResp, cloudwatchMetricResp, err := GetEKSDiskIOPerformancePanel(cmd, clientAuth, nil)
            if err != nil {
                log.Println("Error getting disk I/O performance metrics: ", err)
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

func GetEKSDiskIOPerformancePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

    startTime, endTime := parseTime(startTimeStr, endTimeStr)

    log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

    cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

    totalOpsRawData, err := GetDiskIOPerformanceMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
    if err != nil {
        log.Println("Error fetching total operations raw data: ", err)
        return "", nil, err
    }
    cloudwatchMetricData["TotalOps"] = totalOpsRawData

    result := processDiskIOPerformance(totalOpsRawData)

    jsonString, err := json.Marshal(result)
    if err != nil {
        log.Println("Error marshalling JSON: ", err)
        return "", nil, err
    }

    return string(jsonString), cloudwatchMetricData, nil
}

func GetDiskIOPerformanceMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
    // Define your metric query for disk I/O performance here
    elmType := "ContainerInsights"
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
                        MetricName: aws.String("node_diskio_io_serviced_total"), 
                        Namespace:  aws.String(elmType),
                    },
                    Period: aws.Int64(300),
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

func processDiskIOPerformance(totalOpsRawData *cloudwatch.GetMetricDataOutput) DiskIOPerformanceResult {
    var result DiskIOPerformanceResult

    result.TotalOps = make([]struct {
        Timestamp time.Time
        Value     float64
    }, len(totalOpsRawData.MetricDataResults[0].Timestamps))
    for i, timestamp := range totalOpsRawData.MetricDataResults[0].Timestamps {
        result.TotalOps[i].Timestamp = *timestamp
        result.TotalOps[i].Value = *totalOpsRawData.MetricDataResults[0].Values[i]
    }

    return result
}

func init() {
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("elementId", "", "element id")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("elementType", "", "element type")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("query", "", "query")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("vaultToken", "", "vault token")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("zone", "", "aws region")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("accessKey", "", "aws access key")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("secretKey", "", "aws secret key")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("externalId", "", "aws external id")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("instanceId", "", "instance id")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("startTime", "", "start time")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("endTime", "", "endcl time")
    AwsxEKSDiskIOPerformanceCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
