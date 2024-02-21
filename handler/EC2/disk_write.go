package EC2

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

type DiskWritePanelData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxEc2DiskWriteCmd = &cobra.Command{
	Use:   "disk_write_utilization_panel",
	Short: "get disk write utilization metrics data",
	Long:  `command to get disk write utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskWritePanel(cmd, clientAuth, nil)
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

func GetDiskWritePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID, _ := cmd.PersistentFlags().GetString("instanceId")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
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
	rawData, err := GetDiskWritePanelMetricData(clientAuth, instanceID, namespace, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RawData"] = rawData

	result := processDiskWritePanelRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetDiskWritePanelMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "AWS/EC2"
	if elementType == "EC2" {
		elmType = "AWS/" + elementType
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
						MetricName: aws.String("DiskWriteBytes"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
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

func processDiskWritePanelRawData(result *cloudwatch.GetMetricDataOutput) DiskWritePanelData {
	var rawDatas DiskWritePanelData
	rawDatas.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawDatas.RawData[i].Timestamp = *timestamp
		rawDatas.RawData[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawDatas
}

func init() {
	AwsxEc2DiskWriteCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2DiskWriteCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
