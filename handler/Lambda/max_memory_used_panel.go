package Lambda

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type MemoryData struct {
// 	FunctionName string
// 	MemoryUnit   float64
// }

var AwsxLambdaMaxMemoryCmd = &cobra.Command{
	Use:   "max_memory_used_panel",
	Short: "get lambda memory metrics data",
	Long:  `Command to get lambda memory metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running from child command")
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
			jsonResp, cloudwatchMetricResp, err := GetLambdaMaxMemoryData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda max memory used data: ", err)
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

func GetLambdaMaxMemoryData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "LambdaInsights", "used_memory_max", startTime, endTime, "Maximum", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Printf("Error in getting lambda memory metric data for function: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Max Memory Used (MB)"] = rawData
	// result := processMaxMemoryRawData(rawData, functionName)
	// // result.FunctionName = functionName

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling json in string: ", err)
	// 	return "", nil, err
	// }

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaMaxMemoryCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaMaxMemoryCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaMaxMemoryCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
