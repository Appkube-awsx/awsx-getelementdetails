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

// type CPUSurplusCreditBalanceResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU_Surplus_Credit_Balance"`
// }

// Define a CLI command to get CPU Surplus Credit Balance for RDS instances
var AwsxRDSCPUSurplusCreditBalanceCmd = &cobra.Command{
	Use:   "cpu_surplus_credit_balance",
	Short: "Get CPU Surplus Credit Balance metrics data for RDS instances",
	Long:  `Command to get CPU Surplus Credit Balance metrics data for RDS instances`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUSurplusCreditBalance(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting CPU Surplus Credit Balance: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case. It prints JSON.
				fmt.Println(jsonResp)
			}
		}

	},
}

// Function to get CPU Surplus Credit Balance metrics data
func GetCPUSurplusCreditBalance(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "CPUSurplusCreditBalance", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu surplus credit balance data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["CPU_Surplus_Credit_Balance"] = rawData

	return "", cloudwatchMetricData, nil
}

// // Function to process raw CPU Surplus Credit Balance metric data
// func CPUsurplusCreditbalanceRawData(result *cloudwatch.GetMetricDataOutput) CPUSurplusCreditBalanceResult {
// 	var rawData CPUSurplusCreditBalanceResult
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
	comman_function.InitAwsCmdFlags(AwsxRDSCPUSurplusCreditBalanceCmd)
}
