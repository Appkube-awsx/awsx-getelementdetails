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

type DiskAvailablePanelData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxEc2DiskAvailableCmd = &cobra.Command{
	Use:   "disk_available_panel",
	Short: "get disk available metrics data",
	Long:  `command to get disk available metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskAvailablePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk available utilization: ", err)
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

func GetDiskAvailablePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")

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

	startTime, endTime, err := parseTimeRange(startTimeStr, endTimeStr)
	if err != nil {
		return "", nil, err
	}

	totalResult, usedResult, err := GetDiskTotalPanelMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting total and used disk space data: ", err)
		return "", nil, err
	}

	result := processDiskAvailablePanelRawData(totalResult, usedResult)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"Total": totalResult,
		"Used":  usedResult,
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetDiskTotalPanelMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, *cloudwatch.GetMetricDataOutput, error) {
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
				Id: aws.String("total"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceID),
							},
						},
						MetricName: aws.String("disk_total"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
			},
			{
				Id: aws.String("used"),
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
		return nil, nil, err
	}

	var totalResult, usedResult *cloudwatch.GetMetricDataOutput

	// Separate the total and used metric data
	for _, metricData := range result.MetricDataResults {
		if *metricData.Id == "total" {
			totalResult = &cloudwatch.GetMetricDataOutput{
				MetricDataResults: []*cloudwatch.MetricDataResult{metricData},
			}
		} else if *metricData.Id == "used" {
			usedResult = &cloudwatch.GetMetricDataOutput{
				MetricDataResults: []*cloudwatch.MetricDataResult{metricData},
			}
		}
	}

	return totalResult, usedResult, nil
}

func processDiskAvailablePanelRawData(totalResult, usedResult *cloudwatch.GetMetricDataOutput) DiskAvailablePanelData {
	var availableData DiskAvailablePanelData

	// Initialize an empty slice to store the raw data
	availableData.RawData = []struct {
		Timestamp time.Time
		Value     float64
	}{}

	// Get total disk space
	var total float64
	for _, metricDataResult := range totalResult.MetricDataResults {
		for _, value := range metricDataResult.Values {
			total += *value
		}
	}

	// Get used disk space
	var used float64
	for _, metricDataResult := range usedResult.MetricDataResults {
		for _, value := range metricDataResult.Values {
			used += *value
		}
	}

	// Calculate available disk space
	available := total - used

	// Append the calculated available disk space to the rawData slice
	availableData.RawData = append(availableData.RawData, struct {
		Timestamp time.Time
		Value     float64
	}{
		Timestamp: time.Now(), // Use current time as timestamp
		Value:     available,
	})

	return availableData
}

func parseTimeRange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
	var (
		startTime, endTime time.Time
		err                error
	)

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing start time: %v", err)
		}
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing end time: %v", err)
		}
	}

	return &startTime, &endTime, nil
}

func init() {
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2DiskAvailableCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
