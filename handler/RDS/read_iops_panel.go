package RDS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type ReadIOPS struct {
// 	ReadIOPS []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"read_iops"`
// }

var AwsxRDSReadIOPSCmd = &cobra.Command{
	Use:   "read_iops_panel",
	Short: "Get read IOPS metrics data",
	Long:  `Command to get read IOPS metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, cloudwatchMetricResp, err := GetRDSReadIOPSPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network transmit throughput data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case: print JSON
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetRDSReadIOPSPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "ReadIOPS", startTime, endTime, "Sum", "DBInstanceIdentifier", cloudWatchClient)

	if err != nil {
		log.Println("Error in getting read iops data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["ReadIOPS"] = rawData

	return "", cloudwatchMetricData, nil
}

// func processRawReadIOPSData(result *cloudwatch.GetMetricDataOutput) ReadIOPS {
// 	var rawData ReadIOPS
// 	rawData.ReadIOPS = make([]struct {
// 		Timestamp time.Time
// 		Value     float64
// 	}, len(result.MetricDataResults[0].Timestamps))

// 	for i, timestamp := range result.MetricDataResults[0].Timestamps {
// 		rawData.ReadIOPS[i].Timestamp = *timestamp
// 		rawData.ReadIOPS[i].Value = *result.MetricDataResults[0].Values[i]
// 	}

// 	return rawData
// }

func init() {
	comman_function.InitAwsCmdFlags(AwsxRDSReadIOPSCmd)
}
