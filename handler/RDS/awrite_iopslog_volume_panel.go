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

type WriteIOPSLogVolume struct {
	WriteIOPS []struct {
		Timestamp time.Time
		Value     float64
	} `json:"write_iops_log_volume"`
}

var AwsxRDSWriteIOPSLogVolumeCmd = &cobra.Command{
	Use:   "write_iops_log_volume_panel",
	Short: "Get write IOPS log volume metrics data",
	Long:  `Command to get write IOPS log volume metrics data`,

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
			jsonResp, _, err := GetRDSWriteIOPSLogVolumePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting write IOPS log volume data: ", err)
				return
			}
			fmt.Println("Write IOPS Log Volume Data:")
			fmt.Println(jsonResp)
		}

	},
}

func GetRDSWriteIOPSLogVolumePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	instanceID, _ := cmd.PersistentFlags().GetString("instanceId")
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

	rawWriteIOPSLogVolumeData, err := GetWriteIOPSLogVolumeMetricData(clientAuth, instanceID, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting write IOPS log volume data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["WriteIOPSLogVolume"] = rawWriteIOPSLogVolumeData

	result := processRawWriteIOPSLogVolumeData(rawWriteIOPSLogVolumeData)
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetWriteIOPSLogVolumeMetricData(clientAuth *model.Auth, instanceID string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", instanceID, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("writeIOPSLogVolume"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
						MetricName: aws.String("WriteIOPS"),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(300), // 5 minutes (in seconds)
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

func processRawWriteIOPSLogVolumeData(result *cloudwatch.GetMetricDataOutput) WriteIOPSLogVolume {
	var rawData WriteIOPSLogVolume
	rawData.WriteIOPS = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.WriteIOPS[i].Timestamp = *timestamp
		rawData.WriteIOPS[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSWriteIOPSLogVolumeCmd.PersistentFlags().String("endTime", "", "endcl time")
}
