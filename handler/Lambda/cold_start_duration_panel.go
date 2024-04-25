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

var AwsxLambdaColdStartCmd = &cobra.Command{
	Use:   "cold_start_duration_panel",
	Short: "get lambda cold start duration metrics data",
	Long:  `Command to get lambda cold start duration metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaColdStartData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda cold start duration data: ", err)
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

func GetLambdaColdStartData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := metricData.GetMetricFunctionNameData(clientAuth, instanceId, "LambdaInsights", "init_duration", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cold start duration data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Cold Start Duration"] = rawData

	return "", cloudwatchMetricData, nil
}

func init() {
	AwsxLambdaColdStartCmd.PersistentFlags().String("startTime", "", "Start time")
	AwsxLambdaColdStartCmd.PersistentFlags().String("endTime", "", "End time")
	AwsxLambdaColdStartCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
