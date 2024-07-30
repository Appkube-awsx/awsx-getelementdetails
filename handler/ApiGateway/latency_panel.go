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

// type ApiLatency struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Latency "`
// }

var AwsxApiLatencyCmd = &cobra.Command{
	Use:   "api_latency_panel",
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
			jsonResp, cloudwatchMetricResp, err := GetApiLatencyData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API latency data: ", err)
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

func GetApiLatencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "Latency", startTime, endTime, "Sum", "ApiName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting latency metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Latency"] = metricValue

	var totalSum float64
	for _, value := range metricValue.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{request count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil
}

// func processLatencyRawData(result *cloudwatch.GetMetricDataOutput) ApiLatency {
// 	var rawData ApiLatency
// 	rawData.RawData = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.RawData[i].Timestamp = *timestamp
// 		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
// 	}
// 	return rawData
// }

func init() {
	comman_function.InitAwsCmdFlags(AwsxApiLatencyCmd)
}
