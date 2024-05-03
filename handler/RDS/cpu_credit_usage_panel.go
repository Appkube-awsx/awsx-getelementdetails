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

// type CreditUsageResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"CPU_Credit_Usage"`
// }

var AwsxRDSCPUCreditUsageCmd = &cobra.Command{
	Use:   "cpu_credit_usage_panel",
	Short: "Get CPU credit usage metrics data for RDS instances",
	Long:  `Command to get CPU credit usage metrics data for RDS instances`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUCreditUsagePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting CPU credit usage: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case. It prints JSON
				fmt.Println(jsonResp)
			}
		}
	},
}

func GetCPUCreditUsagePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := comman_function.GetMetricData(clientAuth, instanceId, "AWS/RDS", "CPUCreditUsage", startTime, endTime, "Sum", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu credit usage data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["CPU_Credit_Usage"] = rawData
	return "", cloudwatchMetricData, nil

}

// func processRawCreditUsageData(result *cloudwatch.GetMetricDataOutput) CreditUsageResult {
// 	var rawData CreditUsageResult
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
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSCPUCreditUsageCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
