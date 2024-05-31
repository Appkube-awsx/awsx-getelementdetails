package EC2

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

// DiskUtilization holds the used and free space information
type DiskUtilization struct {
	UsedSpace float64 `json:"usedSpace"`
	FreeSpace float64 `json:"freeSpace"`
}

var AwsxEc2DiskUtilizationCmd = &cobra.Command{
	Use:   "disk_used_panel",
	Short: "get disk used metrics data",
	Long:  `command to get disk used metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskUtilizationData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk read utilization: ", err)
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

func GetDiskUtilizationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	//instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	instanceId := "i-0f095714b7c326e6f"
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	// instanceId, err = comman_function.GetCmdbData(cmd)
	// if err != nil {
	// 	return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	// }

	totalDiskSpace := 100.0 // Assuming total disk space is 100 (representing 100%)
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "CWAgent", "disk_used_percent", startTime, endTime, "Average", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk used data: ", err)
		return "", nil, err
	}

	diskUtilization := DiskUtilization{}

	// Assuming there is only one result and one value
	if len(rawData.MetricDataResults) > 0 && len(rawData.MetricDataResults[0].Values) > 0 {
		percentage := *rawData.MetricDataResults[0].Values[0]
		diskUtilization.UsedSpace = (percentage / 100) * totalDiskSpace
		diskUtilization.FreeSpace = totalDiskSpace - diskUtilization.UsedSpace
		//fmt.Printf("Used Space: %.2f, Free Space: %.2f\n", diskUtilization.UsedSpace, diskUtilization.FreeSpace)
	} else {
		log.Println("No data found")
	}

	cloudwatchMetricData["Used"] = rawData

	// Return the struct as JSON
	jsonString, err := json.Marshal(diskUtilization)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2DiskUtilizationCmd)
}
