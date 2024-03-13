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

type TimeseriesDataPoint struct {
	Timestamp    time.Time `json:"Timestamp"`
	Availability float64   `json:"Availability"`
}

var AwsxEKSServiceAvailabilityCmd = &cobra.Command{
	Use:   "service_availability_panel",
	Short: "get service availability metrics data",
	Long:  `command to get service availability metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetServiceAvailabilityData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Service availability data: ", err)
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


func GetServiceAvailabilityData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, []TimeseriesDataPoint, error) {
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

	rawData, err := GetServiceAvailabilityMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	var timeSeriesData []TimeseriesDataPoint
	for i, timestamp := range rawData.MetricDataResults[0].Timestamps {
		availability := ProcessServiceAvailabilityRawData(rawData, i)
		dataPoint := TimeseriesDataPoint{
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

func GetServiceAvailabilityMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("pod_status_running"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("SampleCount"),
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
						MetricName: aws.String("pod_status_pending"),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("SampleCount"),
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
						MetricName: aws.String("pod_status_ready"),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("SampleCount"),
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

func ProcessServiceAvailabilityRawData(result *cloudwatch.GetMetricDataOutput, index int) float64 {
	// Calculate service availability based on the metrics
	totalRunning := float64(0)
	totalPending := float64(0)
	totalReady := float64(0)

	for _, result := range result.MetricDataResults {
		if *result.Id == "m1" {
			for _, value := range result.Values {
				totalRunning += *value
			}
		} else if *result.Id == "m2" {
			for _, value := range result.Values {
				totalPending += *value
			}
		} else if *result.Id == "m3" {
			for _, value := range result.Values {
				totalReady += *value
			}
		}
	}

	// Calculate service availability
	totalPods := totalRunning + totalPending + totalReady
	if totalPods > 0 {
		return (totalReady / totalPods) * 100
	} else {
		return 0
	}
}

func init() {
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSServiceAvailabilityCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}