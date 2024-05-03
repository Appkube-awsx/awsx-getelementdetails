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

// Define a struct to hold the result of the CPU Surplus Credit Balance query
// type CPUSurplusCreditChargedResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU_Surplus_Credit_Charged"`
// }

// Define a CLI command to get CPU Surplus Credit Balance for RDS instances
var AwsxRDSSurplusCreditsChargedCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetCPUSurplusCreditCharged(cmd, clientAuth, nil)
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
func GetCPUSurplusCreditCharged(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	startTime, endTime, err := comman_function.ParseTimes(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}
	instanceId, err = comman_function.GetCmdbData(cmd)

	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "CPUSurplusCreditsCharged", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu surplus credits charged  data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["CPU_Surplus_Credit_Balance"] = rawData

	return "", cloudwatchMetricData, nil
}

// Function to process raw CPU Surplus Credit Charged metric data
// func CPUsurplusCreditChargedRawData(result *cloudwatch.GetMetricDataOutput) CPUSurplusCreditChargedResult {
// 	var rawData CPUSurplusCreditChargedResult
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
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSSurplusCreditsChargedCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
