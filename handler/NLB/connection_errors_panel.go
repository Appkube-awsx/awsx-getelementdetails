package NLB

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
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

type TimeSeriesData struct {
	Timestamp time.Time
	Count     float64
}

type ConnectionErrorsResult struct {
	ConnectionErrors []TimeSeriesData `json:"ConnectionErrors"`
}

var AwsxNLBConnectionErrorsCmd = &cobra.Command{
	Use:   "connection_errors_panel",
	Short: "Get NLB connection errors metrics data",
	Long:  `Command to get NLB connection errors metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			jsonResp, cloudwatchMetricResp, err := GetNLBConnectionErrorsData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB connection errors: ", err)
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

func GetNLBConnectionErrorsData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	//loadBalancerArn, _ := cmd.PersistentFlags().GetString("loadBalancerArn")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	// elementType, _ := cmd.PersistentFlags().GetString("elementType")
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
		defaultStartTime := time.Now().Add(-1 * time.Hour) // Default to 1 hour ago
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

	// Debug prints
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Fetch raw data
	rawData, err := GetNLBConnectionErrorsMetricData(clientAuth, instanceId, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	// Process the raw data if needed
	result := processConnectionErrorsRawData(rawData)

	timestamps := make([]time.Time, len(result.ConnectionErrors))
	values := make([]float64, len(result.ConnectionErrors))

	// Populate the slices with actual data
	for i, data := range result.ConnectionErrors {
		// Assigning values directly to slices without taking their addresses
		timestamps[i] = data.Timestamp
		values[i] = data.Count
	}

	// Initialize the MetricDataResults slice
	metricDataResults := make([]*cloudwatch.MetricDataResult, len(result.ConnectionErrors))

	// Populate the MetricDataResults with actual data
	for i := range result.ConnectionErrors {
		metricDataResults[i] = &cloudwatch.MetricDataResult{
			Timestamps: []*time.Time{&timestamps[i]},
			Values:     []*float64{&values[i]},
		}
	}

	// Assign the processed data to cloudwatchMetricData
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"Connection Errors": {
			MetricDataResults: metricDataResults,
		},
	}

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNLBConnectionErrorsMetricData(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	elmType := "AWS/NetworkELB"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("connectionErrors_TCP_Client_Reset_Count"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("LoadBalancer"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("TCP_Client_Reset_Count"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300), // 5-minute period
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("connectionErrors_TCP_ELB_Reset_Count"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("LoadBalancer"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("TCP_ELB_Reset_Count"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300), // 5-minute period
					Stat:   aws.String("Sum"),
				},
			},
			{
				Id: aws.String("connectionErrors_TCP_Target_Reset_Count"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("LoadBalancer"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("TCP_Target_Reset_Count"),
						Namespace:  aws.String(elmType),
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

	return result, nil
}

func processConnectionErrorsRawData(result *cloudwatch.GetMetricDataOutput) ConnectionErrorsResult {
	var rawData ConnectionErrorsResult
	rawData.ConnectionErrors = make([]TimeSeriesData, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.ConnectionErrors[i].Timestamp = *timestamp
		count := *result.MetricDataResults[0].Values[i] + *result.MetricDataResults[1].Values[i] + *result.MetricDataResults[2].Values[i]
		rawData.ConnectionErrors[i].Count = count
	}

	// Sort the data based on timestamps in ascending order
	sort.Slice(rawData.ConnectionErrors, func(i, j int) bool {
		return rawData.ConnectionErrors[i].Timestamp.Before(rawData.ConnectionErrors[j].Timestamp)
	})

	return rawData
}

func init() {
	AwsxNLBConnectionErrorsCmd.PersistentFlags().String("instanceId", "", "Instance ID")
	AwsxNLBConnectionErrorsCmd.PersistentFlags().String("startTime", "", "Start time (RFC3339 format)")
	AwsxNLBConnectionErrorsCmd.PersistentFlags().String("endTime", "", "End time (RFC3339 format)")
	AwsxNLBConnectionErrorsCmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
}
