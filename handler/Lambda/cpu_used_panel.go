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

// type CpuResult struct {
// 	Value float64 `json:"Value"`
// }

var AwsxLambdaCpuCmd = &cobra.Command{
	Use:   "cpu_panel",
	Short: "get cpu metrics data",
	Long:  `command to get cpu metrics data`,

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
				log.Println("Error getting lambda cpu data : ", err)
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

func GetLambdaCpuData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	CpuUsedValue, err := commanFunction.GetMetricFunctionNameData(clientAuth, instanceId, "LambdaInsights", "cpu_total_time", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu used value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CpuUsedValue"] = CpuUsedValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaCpuCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaCpuCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaCpuCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaCpuCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaCpuCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaCpuCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaCpuCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaCpuCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaCpuCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaCpuCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaCpuCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaCpuCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaCpuCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaCpuCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaCpuCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaCpuCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
