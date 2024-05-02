package ApiGateway

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type Api5xxResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"5xx Errors"`
// }

var AwsxApi5xxErrorCmd = &cobra.Command{
	Use:   "api_5xxerror_panel",
	Short: "get 5xxerror metrics data",
	Long:  `command to get 5xxerror metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApi5xxErrorData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API 5xx error data: ", err)
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

func GetApi5xxErrorData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	metricValue, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "5XXError", startTime, endTime, "Sum", "ApiName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting 5xx error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["5XXError"] = metricValue

	return "", cloudwatchMetricData, nil
}

// func process5xxErrorRawData(result *cloudwatch.GetMetricDataOutput) Api5xxResult {
// 	var rawData Api5xxResult
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
	AwsxApi5xxErrorCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApi5xxErrorCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApi5xxErrorCmd.PersistentFlags().String("query", "", "query")
	AwsxApi5xxErrorCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApi5xxErrorCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApi5xxErrorCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApi5xxErrorCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApi5xxErrorCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApi5xxErrorCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApi5xxErrorCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApi5xxErrorCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApi5xxErrorCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApi5xxErrorCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApi5xxErrorCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApi5xxErrorCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApi5xxErrorCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApi5xxErrorCmd.PersistentFlags().String("ApiName", "", "api name")
}
