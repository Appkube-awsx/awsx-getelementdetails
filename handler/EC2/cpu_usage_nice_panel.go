package EC2

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

type CpuUsageNice struct {
	CPU_Nice []struct {
		Timestamp time.Time
		Value     float64
	} `json:"CpuUsageNice"`
}

var AwsxEc2CpuUsageNiceCmd = &cobra.Command{
	Use:   "cpu_usage_nice_utilization_panel",
	Short: "get cpu usage nice utilization metrics data",
	Long:  `command to get cpu usage nice utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetCPUUsageNicePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetCPUUsageNicePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data
	rawData, err := GetCpuUsageNiceUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu usage nice data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CpuUsageNice"] = rawData

	result := processedRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetCpuUsageNiceUtilizationMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "CWAgent"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceID),
							},
						},
						MetricName: aws.String("cpu_usage_nice"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
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
func processedRawData(result *cloudwatch.GetMetricDataOutput) CpuUsageNice {
	var rawData CpuUsageNice
	rawData.CPU_Nice = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.CPU_Nice[i].Timestamp = *timestamp
		rawData.CPU_Nice[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}
func init() {
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUsageNiceCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
