package ApiGateway

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type CacheMissResult struct {
// 	TimeSeries []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"timeSeries"`
// }

var AwsxApiCacheMissCmd = &cobra.Command{
	Use:   "cache_miss_count_panel",
	Short: "get cache miss count metrics data",
	Long:  `command to get cache miss count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiCacheMissData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API cache miss count data: ", err)
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

func GetApiCacheMissData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	metricValue, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "CacheMissCount", startTime, endTime, "Sum", "ApiName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting API cache miss count metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CacheMiss"] = metricValue

	var totalSum float64
	for _, value := range metricValue.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{request count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil
}

// func processCacheMissRawData(result *cloudwatch.GetMetricDataOutput) CacheMissResult {
// 	var timeSeries []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		timeSeries = append(timeSeries, struct {
// 			Timestamp time.Time
// 			Value     float64
// 		}{
// 			Timestamp: *timestamp,
// 			Value:     *result.MetricDataResults[0].Values[i],
// 		})
// 	}

// 	return CacheMissResult{TimeSeries: timeSeries}
// }

func init() {
	AwsxApiCacheMissCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiCacheMissCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiCacheMissCmd.PersistentFlags().String("query", "", "query")
	AwsxApiCacheMissCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiCacheMissCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiCacheMissCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiCacheMissCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiCacheMissCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiCacheMissCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiCacheMissCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiCacheMissCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiCacheMissCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiCacheMissCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiCacheMissCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiCacheMissCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiCacheMissCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiCacheMissCmd.PersistentFlags().String("ApiName", "", "api name")
}
