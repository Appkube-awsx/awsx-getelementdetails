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

// type TransactionLogsDiskResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Transaction_Logs_Disk_Usage"`
// }

var AwsxRDSTransactionLogsDiskCmd = &cobra.Command{
	Use:   "transaction_logs_disk_usage_panel",
	Short: "get transation logs disk usage metrics data",
	Long:  `command to get transaction logs disk usage metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetTransactionLogsDiskUsagePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting logs disk usage: ", err)
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

func GetTransactionLogsDiskUsagePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

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

	rawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "AWS/RDS", "TransactionLogsDiskUsage", startTime, endTime, "Average", "DBInstanceIdentifier", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting transaction logs disk usage data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Transaction_Logs_Disk_Usage"] = rawData

	return "", cloudwatchMetricData, nil

}

// func processTransactionLogsDiskRawData(result *cloudwatch.GetMetricDataOutput) TransactionLogsDiskResult {
// 	var rawData TransactionLogsDiskResult
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
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSTransactionLogsDiskCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
