package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type TimeSeriesData struct {
	Timestamp time.Time
	Latency   float64
}

var AwsxRDSLatencyAnalysisCmd = &cobra.Command{
	Use:   "latency_analysis_panel",
	Short: "get latency analysis data",
	Long:  `command to get latency analysis data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSLatencyAnalysisData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting latency analysis data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetRDSLatencyAnalysisData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
    elementId, _ := cmd.PersistentFlags().GetString("elementId")
    elementType, _ := cmd.PersistentFlags().GetString("elementType")
    cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
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
    }

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

    log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

    rawLatencyData, err := GetMetricDatas(clientAuth, elementType, startTime, endTime, cloudWatchClient)
    if err != nil {
        log.Println("Error in getting latency data: ", err)
        return "", nil, err
    }

    latencyResult := processRawLatencyData(rawLatencyData)

    // Create a new GetMetricDataOutput instance
    output := &cloudwatch.GetMetricDataOutput{
        MetricDataResults: []*cloudwatch.MetricDataResult{
            {
                Timestamps: make([]*time.Time, len(latencyResult)),
                Values:     make([]*float64, len(latencyResult)),
            },
        },
    }

    // Populate the output with processed latency data
    for i, data := range latencyResult {
        output.MetricDataResults[0].Timestamps[i] = &data.Timestamp
        output.MetricDataResults[0].Values[i] = &data.Latency
    }

    // Prepare cloudwatchMetricData
    cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
        "LatencyAnalysis": output,
    }

    latencyJSON, err := json.Marshal(latencyResult)
    if err != nil {
        log.Println("Error in marshalling latency data to JSON: ", err)
        return "", nil, err
    }

    return string(latencyJSON), cloudwatchMetricData, nil
}


func GetMetricDatas(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for elementType %s in namespace AWS/RDS from %v to %v", elementType, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("ReadLatency"),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("m2"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("WriteLatency"),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(60),
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
	// fmt.Println("hekllo", result)
	return result, nil
}
func processRawLatencyData(result *cloudwatch.GetMetricDataOutput) []TimeSeriesData {
	var processedData []TimeSeriesData

	// Assuming both read and write metrics have the same number of data points
	for i := range result.MetricDataResults[0].Timestamps {
		readLatency := *result.MetricDataResults[0].Values[i]
		writeLatency := *result.MetricDataResults[1].Values[i]

		timestamp := *result.MetricDataResults[0].Timestamps[i]

		// Calculate combined latency (sum of read and write latencies)
		totalLatency := readLatency + writeLatency

		processedData = append(processedData, TimeSeriesData{
			Timestamp: timestamp,
			Latency:   totalLatency,
		})
	}

	sort.Slice(processedData, func(i, j int) bool {
		return processedData[i].Timestamp.Before(processedData[j].Timestamp)
	})

	return processedData
}


func init() {
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSLatencyAnalysisCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

