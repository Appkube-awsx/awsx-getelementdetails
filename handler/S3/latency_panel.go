package S3

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxS3LatencyCmd = &cobra.Command{
	Use:   "latency_panel",
	Short: "get latency metrics data for s3",
	Long:  `command to get latency metrics data for s3`,

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
			jsonResp, cloudwatchMetricResp, err := GetLatencyPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting latency: ", err)
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

func GetLatencyPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")

	bucketName, _ := cmd.PersistentFlags().GetString("bucketName")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	bucketName, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting bucketName: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	TotalRequestLatencyData, err := GetS3MetricsData(clientAuth, bucketName, "AWS/"+elementType, "TotalRequestLatency", startTime, endTime, "Average", "BucketName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting Total Request latancy data: ", err)
		return "", nil, err
	}

	if len(TotalRequestLatencyData.MetricDataResults) > 0 && len(TotalRequestLatencyData.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["TotalRequestLatencyData"] = TotalRequestLatencyData
	} else {
		log.Println("No data available for TotalRequestLatencyData")
	}

	FirstByteLatency, err := GetS3MetricsData(clientAuth, bucketName, "AWS/"+elementType, "FirstByteLatency", startTime, endTime, "Average", "BucketName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting First Byte Latency data: ", err)
		return "", nil, err
	}

	if len(FirstByteLatency.MetricDataResults) > 0 && len(FirstByteLatency.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["FirstByteLatency"] = FirstByteLatency
	} else {
		log.Println("No data available for FirstByteLatency")
	}

	jsonOutput := make(map[string]float64)
	if len(TotalRequestLatencyData.MetricDataResults) > 0 && len(TotalRequestLatencyData.MetricDataResults[0].Values) > 0 {
		jsonOutput["TotalRequestLatencyData"] = *TotalRequestLatencyData.MetricDataResults[0].Values[0]
	}
	if len(FirstByteLatency.MetricDataResults) > 0 && len(FirstByteLatency.MetricDataResults[0].Values) > 0 {
		jsonOutput["FirstByteLatency"] = *FirstByteLatency.MetricDataResults[0].Values[0]
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetS3MetricsData(clientAuth *model.Auth, BucketName, elementType string, metricName string, startTime, endTime *time.Time, statistic string, dimensionsName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", BucketName, elementType, startTime, endTime)
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
								Name:  aws.String(dimensionsName),
								Value: aws.String(BucketName),
							},
							{
								Name:  aws.String("FilterId"),
								Value: aws.String("hello"),
							},
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elementType),
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

	return result, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxS3LatencyCmd)
}
