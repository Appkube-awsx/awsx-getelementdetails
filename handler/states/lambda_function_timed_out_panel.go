package States

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxStepLambdaFunctionTimedOutCmd = &cobra.Command{
	Use:   "lambda_function_timed_out_panel",
	Short: "get lambda function timed out metrics data",
	Long:  `command to get lambda function timed out metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetStepLambdaFunctionTimedOut(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting function time out data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetStepLambdaFunctionTimedOut(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "States", "LambdaFunctionsTimedOut", startTime, endTime, "Sum", "LambdaFunctionArn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting lambda function data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["LambdaFunctionsTimedOut"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxStepLambdaFunctionTimedOutCmd)
}
