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

type DiskUsedPanelData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Disk_Used"`
}

var AwsxEc2DiskUsedCmd = &cobra.Command{
	Use:   "disk_used_panel",
	Short: "get disk used metrics data",
	Long:  `command to get disk used metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskUsedPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk read  utilization: ", err)
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

func GetDiskUsedPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := GetDiskUsedPanelMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting disk used data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Disk_Used"] = rawData

	result := processDiskUsedPanelRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetDiskUsedPanelMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("disk_used"),
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

func processDiskUsedPanelRawData(result *cloudwatch.GetMetricDataOutput) DiskUsedPanelData {
	var rawData DiskUsedPanelData

	// Initialize an empty slice to store the raw data
	rawData.RawData = []struct {
		Timestamp time.Time
		Value     float64
	}{}

	// Iterate over each metric data result
	for _, metricDataResult := range result.MetricDataResults {
		// Iterate over each timestamp and value pair in the current metric data result
		for i, timestamp := range metricDataResult.Timestamps {
			// Append the timestamp and value to the rawData slice
			rawData.RawData = append(rawData.RawData, struct {
				Timestamp time.Time
				Value     float64
			}{
				Timestamp: *timestamp,
				Value:     *metricDataResult.Values[i],
			})
		}
	}

	return rawData
}

func init() {
	AwsxEc2DiskUsedCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2DiskUsedCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
