package Lambda

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type MetricResult struct {
//     Value float64 `json:"Value"`
// }

var AwsxLambdaMemoryCmd = &cobra.Command{
	Use:   "memory_used_panel",
	Short: "get memory metrics data",
	Long:  `command to get memory metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaMemoryData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda memory data : ", err)
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

func GetLambdaMemoryData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	metricValue, err := metricData.GetMetricFunctionNameData(clientAuth, instanceId, "LambdaInsights", "total_memory", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting memory metric value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Memory"] = metricValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaMemoryCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaMemoryCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaMemoryCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaMemoryCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaMemoryCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaMemoryCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaMemoryCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaMemoryCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaMemoryCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaMemoryCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaMemoryCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaMemoryCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaMemoryCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaMemoryCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaMemoryCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaMemoryCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaMemoryCmd.PersistentFlags().String("functionName", "", "Lambda function name")
}
