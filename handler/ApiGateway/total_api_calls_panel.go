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

// type ApiCallsResult struct {
// 	TimeSeries []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"timeSeries"`
// }

var AwsxApiCallsCmd = &cobra.Command{
	Use:   "total_api_calls_panel",
	Short: "get total API calls metrics data",
	Long:  `command to get total API calls metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetApiCallsData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting total API calls data: ", err)
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

func GetApiCallsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/ApiGateway", "Count", startTime, endTime, "Sum", "ApiName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting total API calls metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["TotalApiCalls"] = metricValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxApiCallsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxApiCallsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxApiCallsCmd.PersistentFlags().String("query", "", "query")
	AwsxApiCallsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxApiCallsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxApiCallsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxApiCallsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxApiCallsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxApiCallsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxApiCallsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxApiCallsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxApiCallsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxApiCallsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxApiCallsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxApiCallsCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxApiCallsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxApiCallsCmd.PersistentFlags().String("ApiName", "", "api name")
}
