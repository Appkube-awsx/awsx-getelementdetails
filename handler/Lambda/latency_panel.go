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

type LatencyResult struct {
	Value float64 `json:"Value"`
}

var AwsxLambdaLatencyCmd = &cobra.Command{
	Use:   "latency_panel",
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
			jsonResp, cloudwatchMetricResp, err := GetLambdaLatencyData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda latency data : ", err)
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

func GetLambdaLatencyData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	avgLatencyValue, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/Lambda", "Duration", startTime, endTime, "Average", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average latency value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["AverageLatency"] = avgLatencyValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaLatencyCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaLatencyCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaLatencyCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaLatencyCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaLatencyCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaLatencyCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaLatencyCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaLatencyCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaLatencyCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaLatencyCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaLatencyCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaLatencyCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaLatencyCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaLatencyCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaLatencyCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaLatencyCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
