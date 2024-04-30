package RDS

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type TransactionLogsGenerationResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Transaction_Logs_Generation"`
// }

var AwsxRDSTransactionLogsGenCmd = &cobra.Command{
	Use:   "transaction_logs_generation_panel",
	Short: "get transation logs generation metrics data",
	Long:  `command to get transaction logs generation metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetTransactionLogsGenerationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting transaction logs generation: ", err)
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

func GetTransactionLogsGenerationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricDatabaseData(clientAuth, instanceId, "AWS/RDS", "TransactionLogsGeneration", startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting transaction logs generation data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Transaction_Logs_Generation"] = rawData

	return "", cloudwatchMetricData, nil

}

// func processTransactionLogRawData(result *cloudwatch.GetMetricDataOutput) TransactionLogsGenerationResult {
// 	var rawData TransactionLogsGenerationResult
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
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSTransactionLogsGenCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
