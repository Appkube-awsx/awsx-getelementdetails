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

// type ErrorResult struct {
// 	Value float64 `json:"Value"`
// }

var AwsxLambdaErrorCmd = &cobra.Command{
	Use:   "error_panel",
	Short: "get error metrics data",
	Long:  `command to get error metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaErrorData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda errors data : ", err)
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

func GetLambdaErrorData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	avgErrorValue, err := commanFunction.GetMetricFunctionNameData(clientAuth, instanceId, "AWS/Lambda", "Errors", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average error value: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["AverageErrors"] = avgErrorValue

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaErrorCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaErrorCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaErrorCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaErrorCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaErrorCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaErrorCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaErrorCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaErrorCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaErrorCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaErrorCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaErrorCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaErrorCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaErrorCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaErrorCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaErrorCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxLambdaErrorCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
