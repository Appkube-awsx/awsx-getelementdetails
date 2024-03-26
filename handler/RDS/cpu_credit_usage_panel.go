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

type CreditUsageResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"CPU_Credit_Usage"`
}

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
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	if elementId != "" {
		log.Println("Getting cloud-element data from CMDB")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("Using default CMDB URL")
			apiUrl = config.CmdbUrl
		}
		log.Println("CMDB URL: " + apiUrl)
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

	rawData, err := GetCPUCreditUsageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	cloudwatchMetricData["CPU_Credit_Usage"] = rawData

	result := processRawCreditUsageData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling JSON in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetCPUCreditUsageMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)
	elmType := "AWS/RDS"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cpuCreditUsage"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("DBInstanceIdentifier"),
								Value: aws.String("postgresql"),
							},
						},
						MetricName: aws.String("CPUCreditUsage"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
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

func processRawCreditUsageData(result *cloudwatch.GetMetricDataOutput) DBResult {
	var rawData DBResult
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
