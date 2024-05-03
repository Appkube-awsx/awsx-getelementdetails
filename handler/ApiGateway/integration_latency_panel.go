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

// type ApiIntegrationLatencyResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"IntegrationLatency"`
// }

var AwsxApiIntegrationLatencyCmd = &cobra.Command{
	Use:   "api_integration_latency_panel",
	Short: "get integration latency metrics data",
	Long:  `command to get integration latency metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiIntegrationLatencyData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API integration latency data: ", err)
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

func GetApiIntegrationLatencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "IntegrationLatency", startTime, endTime, "Average", "ApiName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting latency metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["IntegrationLatency"] = metricValue

	return "", cloudwatchMetricData, nil
}

// func processIntegrationLatencyRawData(result *cloudwatch.GetMetricDataOutput) ApiIntegrationLatencyResult {
// 	var rawData ApiIntegrationLatencyResult
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
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("query", "", "query")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiIntegrationLatencyCmd.PersistentFlags().String("ApiName", "", "api name")
}
