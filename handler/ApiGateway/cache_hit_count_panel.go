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

// type CacheHitsResult struct {
// 	TimeSeries []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"timeSeries"`
// }

var AwsxApiCacheHitsCmd = &cobra.Command{
	Use:   "cache_hit_count_panel",
	Short: "get cache hits metrics data",
	Long:  `command to get cache hits metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiCacheHitsData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API cache hits data: ", err)
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

func GetApiCacheHitsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "CacheHitCount", startTime, endTime, "Sum", "ApiName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting API cache hits metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CacheHits"] = metricValue

	var totalSum float64
	for _, value := range metricValue.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{request count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil
}

// func processCacheHitsRawData(result *cloudwatch.GetMetricDataOutput) CacheHitsResult {
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

// 	return CacheHitsResult{TimeSeries: timeSeries}
// }

func init() {
	AwsxApiCacheHitsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiCacheHitsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiCacheHitsCmd.PersistentFlags().String("query", "", "query")
	AwsxApiCacheHitsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiCacheHitsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiCacheHitsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiCacheHitsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiCacheHitsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiCacheHitsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiCacheHitsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiCacheHitsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiCacheHitsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiCacheHitsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiCacheHitsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiCacheHitsCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiCacheHitsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiCacheHitsCmd.PersistentFlags().String("ApiName", "", "api name")
}
