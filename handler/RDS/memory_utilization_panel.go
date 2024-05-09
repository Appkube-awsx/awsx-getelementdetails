package RDS

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

// type Results struct {
// 	CurrentUsage float64 `json:"currentUsage"`
// 	AverageUsage float64 `json:"averageUsage"`
// 	MaxUsage     float64 `json:"maxUsage"`
// }

var AwsxRDSMemoryUtilizationCmd = &cobra.Command{
	Use:   "memory_utilization_panel",
	Short: "get memory utilization metrics data",
	Long:  `command to get memory utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization: ", err)
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

func GetRDSMemoryUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	
	// Fetch CloudWatch metric data for current, average, and maximum memory usage
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	// Get current usage
	currentUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "FreeableMemory", startTime, endTime, "SampleCount", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting current usage: ", err)
		return "", nil, err
	}
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CurrentUsage"] = currentUsage
	} else {
		log.Println("No data available for current usage")
	}

	// Get average usage
	averageUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "FreeableMemory", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average usage: ", err)
		return "", nil, err
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = averageUsage
	} else {
		log.Println("No data available for average usage")
	}

	// Get maximum usage
	maxUsage, err :=comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "FreeableMemory", startTime, endTime, "Maximum", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting maximum usage: ", err)
		return "", nil, err
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MaxUsage"] = maxUsage
	} else {
		log.Println("No data available for maximum usage")
	}

	// Create JSON output
	jsonOutput := make(map[string]float64)
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["CurrentUsage"] = *currentUsage.MetricDataResults[0].Values[0]
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["AverageUsage"] = *averageUsage.MetricDataResults[0].Values[0]
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["MaxUsage"] = *maxUsage.MetricDataResults[0].Values[0]
	}

	// Convert JSON output to string
	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}



func init() {
	comman_function.InitAwsCmdFlags(AwsxRDSMemoryUtilizationCmd )
}
