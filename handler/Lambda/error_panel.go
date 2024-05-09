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

func GetLambdaErrorData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]interface{}, error) {
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

	cloudwatchMetricData := map[string]interface{}{}

	// Fetch raw data for last month and current month
	lastMonthStartTime := startTime.AddDate(0, -1, 0)
	lastMonthEndTime := endTime.AddDate(0, -1, 0)
	lastMonthMemory, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/Lambda", "Errors", &lastMonthStartTime, &lastMonthEndTime, "Sum", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting error metric value for last month: ", err)
		return "", nil, err
	}

	lastMonthValue := float64(0)
	if len(lastMonthMemory.MetricDataResults) > 0 && len(lastMonthMemory.MetricDataResults[0].Values) > 0 {
		lastMonthValue = *lastMonthMemory.MetricDataResults[0].Values[0]
	}
	cloudwatchMetricData["LastMonthMemory"] = lastMonthValue

	currentMonthMemory, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/Lambda", "Errors", startTime, endTime, "Sum", "FunctionName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting error metric value for current month: ", err)
		return "", nil, err
	}
    fmt.Println(currentMonthMemory)
	currentMonthValue := float64(0)
	if len(currentMonthMemory.MetricDataResults) > 0 && len(currentMonthMemory.MetricDataResults[0].Values) > 0 {
		currentMonthValue = *currentMonthMemory.MetricDataResults[0].Values[0]
	}
	cloudwatchMetricData["CurrentMemory"] = currentMonthValue

	// Calculate percentage change
	var percentageChange float64
	if lastMonthValue != 0 {
		percentageChange = ((currentMonthValue - lastMonthValue) / lastMonthValue) * 100
	} 

	// Determine if it's an increment or decrement
	changeType := "increment"
	if percentageChange < 0 {
		changeType = "decrement"
	}

	cloudwatchMetricData["PercentageChange"] = percentageChange
	cloudwatchMetricData["ChangeType"] = changeType

	return "", cloudwatchMetricData, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxLambdaErrorCmd)
}
