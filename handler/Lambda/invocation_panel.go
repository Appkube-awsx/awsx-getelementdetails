package Lambda

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxLambdaInvocationCmd = &cobra.Command{
	Use:   "invocation_panel",
	Short: "Get duration metrics data for a Lambda function",
	Long:  `Command to get duration metrics data for a Lambda function`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running duration panel command")

		var authFlag bool
		var clientAuth *model.Auth
		var err error
		authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaInvocationData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetLambdaInvocationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID := "appkube-ecommerce-api-dev-updateProduct"
	//metricName := "Invocations"
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	fmt.Println(elementId)
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	// elementId, err = comman_function.GetCmdbData(cmd)
	// if err != nil {
	// 	return "", nil, fmt.Errorf("error getting element ID: %v", err)
	// }
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	success, err := comman_function.GetMetricData(clientAuth, instanceID, "AWS/Lambda", "Invocations", startTime, endTime, "Sum", "FunctionName", cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(success.MetricDataResults) > 0 && len(success.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["Successfull"] = success
	} else {
		log.Println("No data available for Successfull ")
	}

	errordata, err := comman_function.GetMetricData(clientAuth, instanceID, "AWS/Lambda", "Errors", startTime, endTime, "Sum", "FunctionName", cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(errordata.MetricDataResults) > 0 && len(errordata.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["Error"] = errordata
	} else {
		log.Println("No data available for Error Usage")
	}
	coldstart, err := comman_function.GetMetricData(clientAuth, instanceID, "LambdaInsights", "init_duration", startTime, endTime, "Sum", "function_name", cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(coldstart.MetricDataResults) > 0 && len(coldstart.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["ColdStart"] = coldstart
	} else {
		log.Println("No data available for ColdStart")
	}

	jsonOutput := make(map[string]float64)

	if len(success.MetricDataResults) > 0 && len(success.MetricDataResults[0].Values) > 0 {
		jsonOutput["Successfull"] = *success.MetricDataResults[0].Values[0]
	}
	if len(errordata.MetricDataResults) > 0 && len(errordata.MetricDataResults[0].Values) > 0 {
		jsonOutput["Error"] = *errordata.MetricDataResults[0].Values[0]
	}
	if len(coldstart.MetricDataResults) > 0 && len(coldstart.MetricDataResults[0].Values) > 0 {
		jsonOutput["ColdStart"] = *coldstart.MetricDataResults[0].Values[0]
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}
	return string(jsonString), cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaInvocationCmd.PersistentFlags().String("elementId", "", "Lambda function name or ID")
	AwsxLambdaInvocationCmd.PersistentFlags().String("elementType", "", "Element type")
	AwsxLambdaInvocationCmd.PersistentFlags().String("startTime", "", "Start time for metrics collection (RFC333)")
}
