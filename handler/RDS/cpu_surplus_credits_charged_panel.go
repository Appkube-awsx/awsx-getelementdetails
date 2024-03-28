package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// Define a struct to hold the result of the CPU Surplus Credit Balance query
type CPUSurplusCreditChargedResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"CPU_Surplus_Credit_Balance"`
}

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

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	rawData, err := GetCPUSurplusCreditChargedMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["CPU_Surplus_Credit_Balance"] = rawData

	result := CPUsurplusCreditChargedRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

// Function to retrieve CPU Surplus Credit Balance metric data
func GetCPUSurplusCreditChargedMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting CPU Surplus Credit Balance metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)
	elmType := "AWS/RDS"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cpuSurplusCreditsCharged"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("DBInstanceIdentifier"),
								Value: aws.String("postgresql"), // Ensure instanceID is the identifier of your RDS instance
							},
						},
						MetricName: aws.String("CPUSurplusCreditsCharged"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),        // 5 minutes (in seconds)
					Stat:   aws.String("Average"), // You can use 'Average', 'Sum', 'Minimum', 'Maximum'
				},
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

// Function to process raw CPU Surplus Credit Charged metric data
func CPUsurplusCreditChargedRawData(result *cloudwatch.GetMetricDataOutput) CPUSurplusCreditChargedResult {
	var rawData CPUSurplusCreditChargedResult
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
