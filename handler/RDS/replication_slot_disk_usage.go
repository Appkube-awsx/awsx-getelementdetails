package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type ReplicationSlotDiskUsage struct {
	Timestamp time.Time
	Value     float64
}

var AwsxRDSReplicationSlotDiskUsageCmd = &cobra.Command{
	Use:   "replication_slot_disk_usage_panel",
	Short: "get replication slot disk usage metrics data",
	Long:  `command to get replication slot disk usage metrics data`,

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
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSReplicationSlotDiskUsagePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting replication slot disk usage data: ", err)
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

func GetRDSReplicationSlotDiskUsagePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", "", nil, err
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
			return "", "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data for replication slot disk usage metric
	rawReplicationSlotDiskUsageData, err := GetReplicationSlotDiskUsageMetricData(clientAuth, elementType, startTime, endTime, "ReplicationSlotDiskUsage", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting replication slot disk usage data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["ReplicationSlotDiskUsage"] = rawReplicationSlotDiskUsageData

	// Process raw replication slot disk usage data
	resultReplicationSlotDiskUsage := processedRawReplicationSlotDiskUsageData(rawReplicationSlotDiskUsageData)
	jsonReplicationSlotDiskUsage, err := json.Marshal(resultReplicationSlotDiskUsage)
	if err != nil {
		log.Println("Error in marshalling json for replication slot disk usage data: ", err)
		return "", "", nil, err
	}

	return string(jsonReplicationSlotDiskUsage), string(jsonReplicationSlotDiskUsage), cloudwatchMetricData, nil
}

func GetReplicationSlotDiskUsageMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace AWS/RDS from %v to %v", elementType, startTime, endTime)

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
						MetricName: aws.String(metricName),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Average"), // Use "Average" for replication slot disk usage metrics
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

func processedRawReplicationSlotDiskUsageData(result *cloudwatch.GetMetricDataOutput) []ReplicationSlotDiskUsage {
	var processedData []ReplicationSlotDiskUsage

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		processedData = append(processedData, ReplicationSlotDiskUsage{
			Timestamp: *timestamp,
			Value:     value,
		})
	}

	return processedData
}

func init() {
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSReplicationSlotDiskUsageCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

