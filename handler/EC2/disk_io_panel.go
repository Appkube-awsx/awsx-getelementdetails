package EC2

import (
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

type DiskIOPerformanceResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxEC2DiskIOPerformanceCmd = &cobra.Command{
	Use:   "disk_io_performance_panel",
	Short: "get disk I/O performance metrics data",
	Long:  `command to get disk I/O performance metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetDiskReadPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting disk I/O performance metrics: ", err)
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

func GetEC2DiskIOPerformancePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data for DiskReadBytes
	rawDataDiskReadBytes, err := GetDiskIOMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "DiskReadBytes", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data for DiskReadBytes: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DiskReadBytes"] = rawDataDiskReadBytes

	// Fetch raw data for DiskWriteBytes
	rawDataDiskWriteBytes, err := GetDiskIOMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "DiskWriteBytes", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data for DiskWriteBytes: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DiskWriteBytes"] = rawDataDiskWriteBytes

	//resultDiskReadBytes := processRawPanelRawData(rawDataDiskReadBytes)
	//resultDiskWriteBytes := processRawPanelRawData(rawDataDiskWriteBytes)
	//
	//// Calculate total disk I/O
	//totalDiskIO := make([]struct {
	//	Timestamp time.Time
	//	Value     float64
	//}, len(resultDiskReadBytes))
	//
	//for i := range resultDiskReadBytes.RawData {
	//	totalDiskIO[i].Timestamp = resultDiskReadBytes.RawData[i].Timestamp
	//	totalDiskIO[i].Value = resultDiskReadBytes.RawData[i].Value + resultDiskWriteBytes.RawData[i].Value
	//}
	//
	//jsonResults := map[string]interface{}{
	//	"DiskReadBytes":  resultDiskReadBytes,
	//	"DiskWriteBytes": resultDiskWriteBytes,
	//	"TotalDiskIO":    totalDiskIO,
	//}

	//jsonString, err := json.Marshal(jsonResults)
	//if err != nil {
	//	log.Println("Error in marshalling json in string: ", err)
	//	return "", nil, err
	//}

	return "", cloudwatchMetricData, nil
}

func GetDiskIOMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "AWS/EC2"

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
						MetricName: aws.String(metricName),
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

	if len(result.MetricDataResults) == 0 || len(result.MetricDataResults[0].Values) == 0 {
		return nil, fmt.Errorf("no data available for the specified time range")
	}

	// If there is only one value, return it
	if len(result.MetricDataResults[0].Values) == 1 {
		return result, nil
	}

	// If there are multiple values, calculate the average
	var sum float64
	for _, v := range result.MetricDataResults[0].Values {
		sum += aws.Float64Value(v)
	}
	// average := sum / float64(len(result.MetricDataResults[0].Values))

	return result, nil
}

//
//func processRawPanelRawData(result *cloudwatch.GetMetricDataOutput) DiskReadPanelData {
//	var rawData DiskReadPanelData
//
//	// Initialize an empty slice to store the raw data
//	rawData.RawData = []struct {
//		Timestamp time.Time
//		Value     float64
//	}{}
//
//	// Iterate over each metric data result
//	for _, metricDataResult := range result.MetricDataResults {
//		// Iterate over each timestamp and value pair in the current metric data result
//		for i, timestamp := range metricDataResult.Timestamps {
//			// Append the timestamp and value to the rawData slice
//			rawData.RawData = append(rawData.RawData, struct {
//				Timestamp time.Time
//				Value     float64
//			}{
//				Timestamp: *timestamp,
//				Value:     *metricDataResult.Values[i],
//			})
//		}
//	}
//
//	return rawData
//}

func init() {
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("query", "", "query")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEC2DiskIOPerformanceCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
