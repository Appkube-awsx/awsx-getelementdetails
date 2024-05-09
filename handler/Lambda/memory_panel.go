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

var AwsxLambdaMemoryUsageCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetLambdaMemoryUsageData(cmd, clientAuth, nil)
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

func GetLambdaMemoryUsageData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID := "appkube-ecommerce-api-dev-updateProduct"
	//metricName := "Invocations"
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	fmt.Println(elementId)
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	// elementId, err = comman_function.GetCmdbData(cmd)
	// if err != nil {
	// 	return "", nil, fmt.Errorf("error getting element ID: %v", err)
	// }
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	Average, err := comman_function.GetMetricData(clientAuth, instanceID, elementType+"Insights", "used_memory_max", startTime, endTime, "Average", "function_name", cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(Average.MetricDataResults) > 0 && len(Average.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["Average"] = Average
	} else {
		log.Println("No data available for Average ")
	}

	Maximum, err := comman_function.GetMetricData(clientAuth, instanceID, elementType+"Insights", "used_memory_max", startTime, endTime, "Maximum", "function_name", cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(Maximum.MetricDataResults) > 0 && len(Maximum.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["Maximum"] = Maximum
	} else {
		log.Println("No data available for Maximum Usage")
	}
	Minimum, err := comman_function.GetMetricData(clientAuth, instanceID, elementType+"Insights", "used_memory_max", startTime, endTime, "Minimum", "function_name", cloudWatchClient)
	if err != nil {
		return "", nil, err
	}
	if len(Minimum.MetricDataResults) > 0 && len(Minimum.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["Minimum"] = Minimum
	} else {
		log.Println("No data available for Minimum Usage")
	}

	jsonOutput := make(map[string]float64)

	if len(Average.MetricDataResults) > 0 && len(Average.MetricDataResults[0].Values) > 0 {
		jsonOutput["Average"] = *Average.MetricDataResults[0].Values[0]
	}
	if len(Maximum.MetricDataResults) > 0 && len(Maximum.MetricDataResults[0].Values) > 0 {
		jsonOutput["Maximum"] = *Maximum.MetricDataResults[0].Values[0]
	}
	if len(Minimum.MetricDataResults) > 0 && len(Minimum.MetricDataResults[0].Values) > 0 {
		jsonOutput["Minimum"] = *Minimum.MetricDataResults[0].Values[0]
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}
	return string(jsonString), cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaMemoryUsageCmd.PersistentFlags().String("elementId", "", "Lambda function name or ID")
	AwsxLambdaMemoryUsageCmd.PersistentFlags().String("elementType", "", "Element type")
	AwsxLambdaMemoryUsageCmd.PersistentFlags().String("startTime", "", "Start time for metrics collection (RFC333)")
}
