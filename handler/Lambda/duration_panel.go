package Lambda

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxLambdaDurationCmd = &cobra.Command{
	Use:   "duration_panel",
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
			jsonResp, cloudwatchMetricResp, err := GetLambdaDurationData(cmd, clientAuth, nil)
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

func GetLambdaDurationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID := "appkube-ecommerce-api-dev-getAllOrders"
	metricName := "Duration"
	dimensionsName := "FunctionName"
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
	average, err := comman_function.GetMetricData(clientAuth, instanceID, "AWS/"+elementType, metricName, startTime, endTime, "Average", dimensionsName, cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(average.MetricDataResults) > 0 && len(average.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = average
	} else {
		log.Println("No data available for average Usage")
	}
	minimum, err := comman_function.GetMetricData(clientAuth, instanceID, "AWS/"+elementType, metricName, startTime, endTime, "Minimum", dimensionsName, cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(minimum.MetricDataResults) > 0 && len(minimum.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MinUsage"] = minimum
	} else {
		log.Println("No data available for minimum Usage")
	}
	maximum, err := comman_function.GetMetricData(clientAuth, instanceID, "AWS/"+elementType, metricName, startTime, endTime, "Maximum", dimensionsName, cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(maximum.MetricDataResults) > 0 && len(maximum.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MaxUsage"] = maximum
	} else {
		log.Println("No data available for maximum Usage")
	}

	jsonOutput := make(map[string]float64)

	if len(average.MetricDataResults) > 0 && len(average.MetricDataResults[0].Values) > 0 {
		jsonOutput["AverageUsage"] = *average.MetricDataResults[0].Values[0]
	}
	if len(minimum.MetricDataResults) > 0 && len(minimum.MetricDataResults[0].Values) > 0 {
		jsonOutput["MinUsage"] = *minimum.MetricDataResults[0].Values[0]
	}
	if len(maximum.MetricDataResults) > 0 && len(maximum.MetricDataResults[0].Values) > 0 {
		jsonOutput["MaxUsage"] = *maximum.MetricDataResults[0].Values[0]
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaDurationCmd.PersistentFlags().String("elementId", "", "Lambda function name or ID")
	AwsxLambdaDurationCmd.PersistentFlags().String("elementType", "", "Element type")
	AwsxLambdaDurationCmd.PersistentFlags().String("startTime", "", "Start time for metrics collection (RFC333)")
}
