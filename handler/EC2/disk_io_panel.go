package EC2

import (
	//"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type DiskIOPerformanceResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"RawData"`
// }

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
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := comman_function.ParseTimes(cmd)

	
		if err != nil {
			return "", nil, fmt.Errorf("error parsing time: %v", err)
		}
		instanceId, err = comman_function.GetCmdbData(cmd)

		
		if err != nil {
			return "", nil, fmt.Errorf("error getting instance ID: %v", err)
		}	

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data for DiskReadBytes
	rawDataDiskReadBytes, err := comman_function.GetMetricData(clientAuth, instanceId, "CWAgent","diskio_read_bytes", startTime, endTime, "Sum","InstanceId",  cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data for DiskReadBytes: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["DiskReadBytes"] = rawDataDiskReadBytes

	// Fetch raw data for DiskWriteBytes
	rawDataDiskWriteBytes, err := comman_function.GetMetricData(clientAuth, instanceId,  "CWAgent","diskio_write_bytes", startTime, endTime, "Sum","InstanceId",   cloudWatchClient)
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
		//"DiskReadBytes":  resultDiskReadBytes,
	//"DiskWriteBytes": resultDiskWriteBytes,
		//"TotalDiskIO":    totalDiskIO,
	//}

	//jsonString, err := json.Marshal(jsonResults)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	//return "", cloudwatchMetricData, nil
	var totalSum float64
	for _, value := range rawDataDiskReadBytes.MetricDataResults {
		for _, datum := range value.Values {
			totalSum += *datum
		}
	}
	totalSumStr := fmt.Sprintf("{disk io count: %f}", totalSum)
	return totalSumStr, cloudwatchMetricData, nil
}

// func GetDiskIOMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
// 	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

// 	elmType := "AWS/EC2"

// 	input := &cloudwatch.GetMetricDataInput{
// 		EndTime:   endTime,
// 		StartTime: startTime,
// 		MetricDataQueries: []*cloudwatch.MetricDataQuery{
// 			{
// 				Id: aws.String("m1"),
// 				MetricStat: &cloudwatch.MetricStat{
// 					Metric: &cloudwatch.Metric{
// 						Dimensions: []*cloudwatch.Dimension{
// 							{
// 								Name:  aws.String("InstanceId"),
// 								Value: aws.String(instanceID),
// 							},
// 						},
// 						MetricName: aws.String(metricName),
// 						Namespace:  aws.String(elmType),
// 					},
// 					Period: aws.Int64(300),
// 					Stat:   aws.String(statistic),
// 				},
// 			},
// 		},
	// }
	// if cloudWatchClient == nil {
	// 	cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	// }

	// result, err := cloudWatchClient.GetMetricData(input)
// if err != nil {
// 		log.Println("Error in marshalling json: ", err)
//  		return "", nil, err
//  	}
//  }

	// if len(result.MetricDataResults) == 0 || len(result.MetricDataResults[0].Values) == 0 {
	// 	return nil, fmt.Errorf("no data available for the specified time range")
	// }

	// // If there is only one value, return it
	// if len(result.MetricDataResults[0].Values) == 1 {
	// 	return result, nil
	// }

	// If there are multiple values, calculate the average
// 	var sum float64
// 	for _, v := range result.MetricDataResults[0].Values {
// 		sum += aws.Float64Value(v)
// 	}
// 	// average := sum / float64(len(result.MetricDataResults[0].Values))

// 	return result, nil
//}

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
	comman_function.InitAwsCmdFlags(AwsxEC2DiskIOPerformanceCmd)
}
