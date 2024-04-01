package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type TransactionLogsGenerationResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Transaction_Logs_Generation"`
}

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
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, err
		}
		instanceId = cmdbData.InstanceId

	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)

			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)

			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := GetTransactionLogsGenerationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["Transaction_Logs_Generation"] = rawData

	result := processTransactionLogRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil

}

func GetTransactionLogsGenerationMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)
	elmType := "AWS/RDS"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("transactionloggeneration"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{

						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("DBInstanceIdentifier"),
								Value: aws.String("postgresql"), // Ensure instanceID is the identifier of your RDS instance
							},
						},
						MetricName: aws.String("TransactionLogsGeneration"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
				//ReturnData: aws.Bool(true),
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processTransactionLogRawData(result *cloudwatch.GetMetricDataOutput) TransactionLogsGenerationResult {
	var rawData TransactionLogsGenerationResult
	rawData.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

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
