package EKS

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

type NetworkAvailabilityResult struct {
	Availability float64 `json:"Availability"`
}

type TimeSeriesDataPoint struct {
	Timestamp    time.Time `json:"Timestamp"`
	Availability float64   `json:"Availability"`
}

var AwsxEKSNetworkAvailabilityCmd = &cobra.Command{
	Use:   "network_availability_panel",
	Short: "get network_availability graph metrics data",
	Long:  `command to get network_availability graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkAvailabilityData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Network availability data: ", err)
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

func GetNetworkAvailabilityData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, []TimeSeriesDataPoint, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

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

	rawData, err := GetNetworkAvailabilityMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	var timeSeriesData []TimeSeriesDataPoint
	for i, timestamp := range rawData.MetricDataResults[0].Timestamps {
		availability := ProcessNetworkAvailabilityRawData(rawData, i)
		dataPoint := TimeSeriesDataPoint{
			Timestamp:    *timestamp,
			Availability: availability,
		}
		timeSeriesData = append(timeSeriesData, dataPoint)
	}

	jsonString, err := json.Marshal(timeSeriesData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), timeSeriesData, nil
}

func GetNetworkAvailabilityMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	elmType := "ContainerInsights"
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
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("node_interface_network_tx_dropped"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("m2"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("node_interface_network_rx_dropped"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("m3"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("pod_network_rx_bytes"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("m4"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("pod_network_tx_bytes"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
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

func ProcessNetworkAvailabilityRawData(result *cloudwatch.GetMetricDataOutput, index int) float64 {
	// Calculate network availability based on the metrics
	totalTxDropped := float64(0)
	totalRxDropped := float64(0)
	totalRxBytes := float64(0)
	totalTxBytes := float64(0)

	for _, result := range result.MetricDataResults {
		if *result.Id == "m1" {
			for _, value := range result.Values {
				totalTxDropped += *value
			}
		} else if *result.Id == "m2" {
			for _, value := range result.Values {
				totalRxDropped += *value
			}
		} else if *result.Id == "m3" {
			for _, value := range result.Values {
				totalRxBytes += *value
			}
		} else if *result.Id == "m4" {
			for _, value := range result.Values {
				totalTxBytes += *value
			}
		}
	}

	// Calculate network availability
	if totalTxBytes > 0 && totalTxBytes > totalTxDropped && totalRxBytes > 0 && totalRxBytes > totalRxDropped {
		return 100 * (1 - ((totalTxDropped + totalRxDropped) / (totalTxBytes + totalRxBytes)))
	} else {
		return 0
	}
}

func init() {
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNetworkAvailabilityCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}