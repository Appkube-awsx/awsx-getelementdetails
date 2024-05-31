package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type CPUReservedResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU_Reservation"`
// }

var AwsxCpuReservedPanelCmd = &cobra.Command{
	Use:   "cpu_reserved_panel",
	Short: "get cpu reserved metrics data",
	Long:  `command to get cpu reserved metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetEC2CPUReservationData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu reserved data : ", err)
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

func GetEC2CPUReservationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	CPU_Reservation, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CpuReserved", startTime, endTime, "Sum", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu reservation raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_Reservation"] = CPU_Reservation
	CPU_Utilization, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/"+elementType, "CPUUtilization", startTime, endTime, "Sum", "InstanceId", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu reservation raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CPU_Utilization"] = CPU_Utilization
	return "", cloudwatchMetricData, nil
}

// func processCPUReservedRawData(result *cloudwatch.GetMetricDataOutput) CPUReservedResult {
// 	var rawData CPUReservedResult
// 	rawData.RawData = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.RawData[i].Timestamp = *timestamp
// 		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
// 	}

// 	return rawData
// }

func init() {
	comman_function.InitAwsCmdFlags(AwsxCpuReservedPanelCmd)
}
