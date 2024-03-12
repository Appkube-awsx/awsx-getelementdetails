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

type CpuUsageIdle struct {
	CPU_Idle []struct {
		Timestamp time.Time
		Value     float64
	} `json:"CpuUsageIdle"`
}

var AwsxEc2CpuUsageIdleCmd = &cobra.Command{
	Use:   "cpu_usage_Idle_utilization_panel",
	Short: "get cpu usage idle utilization metrics data",
	Long:  `command to get cpu usage idle utilization metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command..")
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
			jsonResp, cloudwatchMetricResp, err := GetCPUUsageIdlePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu usage idle utilization: ", err)
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

func GetCPUUsageIdlePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := GetCPUUsageIdleMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting cpu usage idle data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CpuUsageIdle"] = rawData

	result := processTheRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetCPUUsageIdleMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "CWAgent"
	if elementType == "EC2" {
		elmType = "CWAgent"
	}
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
						MetricName: aws.String("cpu_usage_idle"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Average"),
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

func processTheRawData(result *cloudwatch.GetMetricDataOutput) CpuUsageIdle {
	var rawData CpuUsageIdle
	rawData.CPU_Idle = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.CPU_Idle[i].Timestamp = *timestamp
		rawData.CPU_Idle[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}
func init() {
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUsageIdleCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
