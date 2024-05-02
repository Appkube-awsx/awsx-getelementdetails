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

// type Api4xxResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"4xx Errors"`
// }

var AwsxApi4xxErrorCmd = &cobra.Command{
	Use:   "api_4xxerror_panel",
	Short: "get 4xxerror metrics data",
	Long:  `command to get 4xxerror metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApi4xxErrorData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting API 4xx error data: ", err)
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

func GetApi4xxErrorData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "4XXError", startTime, endTime, "Sum", "ApiName", cloudWatchClient)
	//metricValue, err := GetApi4xxErrorMetricValue(clientAuth, ApiName, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting 4xx error metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["4XXError"] = metricValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxApi4xxErrorCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApi4xxErrorCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApi4xxErrorCmd.PersistentFlags().String("query", "", "query")
	AwsxApi4xxErrorCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApi4xxErrorCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApi4xxErrorCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApi4xxErrorCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApi4xxErrorCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApi4xxErrorCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApi4xxErrorCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApi4xxErrorCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApi4xxErrorCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApi4xxErrorCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApi4xxErrorCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApi4xxErrorCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApi4xxErrorCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApi4xxErrorCmd.PersistentFlags().String("ApiName", "", "api name")
}
