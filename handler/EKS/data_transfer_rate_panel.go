package EKS

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

type DataTransferRateDataPoint struct {
	Timestamp   time.Time `json:"Timestamp"`
	TransferIn  float64   `json:"TransferIn"`
	TransferOut float64   `json:"TransferOut"`
}

var AwsxEksDataTransferRateCmd = &cobra.Command{
	Use:   "data_transfer_rate_panel",
	Short: "get EKS data transfer rate metrics data",
	Long:  `command to get EKS data transfer rate metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetEksDataTransferRatePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting EKS data transfer rate data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetEksDataTransferRatePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, []DataTransferRateDataPoint, error) {
	clusterName := "myClustTT"
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-1 * time.Hour) // Default start time: 1 hour ago
		startTime = &defaultStartTime
	}

	// Parse end time if provided
	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now() // Default end time: Now
		endTime = &defaultEndTime
	}

	// Get EKS cluster metrics
	eksMetrics, err := GetEksMetrics(clientAuth, clusterName, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EKS metrics: ", err)
		return "", nil, err
	}

	// Calculate data transfer rate data points
	var dataTransferRateData []DataTransferRateDataPoint
	for i := 0; i < len(eksMetrics.MetricDataResults[0].Timestamps); i++ {
		dataPoint := DataTransferRateDataPoint{
			Timestamp:   *eksMetrics.MetricDataResults[0].Timestamps[i],
			TransferIn:  *eksMetrics.MetricDataResults[0].Values[i],
			TransferOut: *eksMetrics.MetricDataResults[1].Values[i],
		}
		dataTransferRateData = append(dataTransferRateData, dataPoint)
	}

	jsonString, err := json.Marshal(dataTransferRateData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), dataTransferRateData, nil
}

func GetEksMetrics(clientAuth *model.Auth, clusterName string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("transfer_in"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("node_interface_network_rx_dropped"),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(300), // 5-minute period
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("transfer_out"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String("node_interface_network_tx_dropped"),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(300), // 5-minute period
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
	fmt.Println("nulllll", result)
	fmt.Println("input", input)

	return result, nil
}

func init() {
	AwsxEksDataTransferRateCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("query", "", "query")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEksDataTransferRateCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
