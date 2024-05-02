package Lambda

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type NetReceived struct {
// 	Value float64 `json:"Value"`
// }

var AwsxLambdaNetReceivedCmd = &cobra.Command{
	Use:   "net_received_panel",
	Short: "get net received metrics data",
	Long:  `command to get net received metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaNetReceivedData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda net received data : ", err)
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

func GetLambdaNetReceivedData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := commanFunction.GetMetricData(clientAuth, instanceId, "LambdaInsights", "rx_bytes", startTime, endTime, "Average", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting net received metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory"] = metricValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaNetReceivedCmd.PersistentFlags().String("functionName", "", "Lambda function name")
}
