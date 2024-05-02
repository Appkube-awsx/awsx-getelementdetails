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

// type GraphMemoryData struct {
// 	FunctionName string
// 	MemoryUnit   float64
// }

var AwsxLambdaMaxMemoryGraphCmd = &cobra.Command{
	Use:   "max_memory_used_graph_panel",
	Short: "get lambda memory used data",
	Long:  `Command to get lambda memory used data`,

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

			jsonResp, cloudwatchMetricResp, err := GetLambdaMaxMemoryGraphData(cmd, clientAuth, nil)
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

func GetLambdaMaxMemoryGraphData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "LambdaInsights", "used_memory_max", startTime, endTime, "Maximum", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Printf("Error in getting lambda memory metric data for function: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Max Memory Used (MB)"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaMaxMemoryGraphCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaMaxMemoryGraphCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaMaxMemoryGraphCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
