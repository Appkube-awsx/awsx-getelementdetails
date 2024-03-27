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

type WriteIOPS struct {
	WriteIOPS []struct {
		Timestamp time.Time
		Value     float64
	} `json:"write_iops"`
}

var AwsxRDSWriteIOPSCmd = &cobra.Command{
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
			jsonResp, _, err := GetRDSWriteIOPSPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting write IOPS log volume data: ", err)
				return
			}
			fmt.Println("Write IOPS Log Volume Data:")
			fmt.Println(jsonResp)
		}

	},
}

func GetRDSWriteIOPSPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawWriteIOPSData, err := GetWriteIOPSMetricData(clientAuth, instanceID, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting write IOPS log volume data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["WriteIOPS"] = rawWriteIOPSData

	result := processRawWriteIOPSData(rawWriteIOPSData)
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetWriteIOPSMetricData(clientAuth *model.Auth, instanceID string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", instanceID, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("WriteIOPS"),
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

func processRawWriteIOPSData(result *cloudwatch.GetMetricDataOutput) WriteIOPS {
	var rawData WriteIOPS
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
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSWriteIOPSCmd.PersistentFlags().String("endTime", "", "endcl time")
}
