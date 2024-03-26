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

type DiskQueueDepth struct {
	Timestamp time.Time
	Value     float64
}

var AwsxRDSDiskQueueDepthCmd = &cobra.Command{
	Use:   "disk_queue_depth_panel",
	Short: "get disk queue depth metrics data",
	Long:  `command to get disk queue depth metrics data`,

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
			jsonResp, cloudwatchMetricResp, err, _ := GetRDSDiskQueueDepthPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk queue depth data: ", err)
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

func GetRDSDiskQueueDepthPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data for disk queue depth metric
	rawDiskQueueDepthData, err := GetDiskQueueDepthMetricData(clientAuth, elementType, startTime, endTime, "DiskQueueDepth", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk queue depth data: ", err)
		return "", "", nil, err
	}
	cloudwatchMetricData["DiskQueueDepth"] = rawDiskQueueDepthData

	// Process raw disk queue depth data
	resultDiskQueueDepth := processedRawDiskQueueDepthData(rawDiskQueueDepthData)
	jsonDiskQueueDepth, err := json.Marshal(resultDiskQueueDepth)
	if err != nil {
		log.Println("Error in marshalling json for disk queue depth data: ", err)
		return "", "", nil, err
	}

	return string(jsonDiskQueueDepth), string(jsonDiskQueueDepth), cloudwatchMetricData, nil
}

func GetDiskQueueDepthMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
					Stat:   aws.String("Average"), // Use "Average" for disk queue depth metrics
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

func processedRawDiskQueueDepthData(result *cloudwatch.GetMetricDataOutput) []DiskQueueDepth {
	var processedData []DiskQueueDepth

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		value := *result.MetricDataResults[0].Values[i]
		processedData = append(processedData, DiskQueueDepth{
			Timestamp: *timestamp,
			Value:     value,
		})
	}

	return processedData
}


func init() {
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSDiskQueueDepthCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

