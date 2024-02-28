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

type TimeSeriesData struct {
	Timestamp      time.Time
	AllocatableCPU float64
}

type AllocateResult struct {
	RawData []TimeSeriesData `json:"RawData"`
}

var AwsxEKSAllocatableCpuCmd = &cobra.Command{
	Use:   "allocatable_cpu_panel",
	Short: "get allocatable cpu metrics data",
	Long:  `command to get allocatable cpu metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetAllocatableCPUData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting allocatable cpu: ", err)
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

func GetAllocatableCPUData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Parse start time if provided
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

	// Debug prints
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	

	// Fetch raw data
	rawData, err := GetAllocatableCPUMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	
	// Process the raw data if needed
	result := processCPURawData(rawData)
	// log.Println(result)
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"RawData": &cloudwatch.GetMetricDataOutput{
			MetricDataResults: []*cloudwatch.MetricDataResult{
				{
					Timestamps: make([]*time.Time, len(result.RawData)),
					Values:     make([]*float64, len(result.RawData)),
				},
			},
		},
	}
	
	// Assign the processed data to cloudwatchMetricData
	for i, data := range result.RawData {
		cloudwatchMetricData["RawData"].MetricDataResults[0].Timestamps[i] = &data.Timestamp
		cloudwatchMetricData["RawData"].MetricDataResults[0].Values[i] = &data.AllocatableCPU
	}
	
	// log.Printf("CloudWatch Metric Data: %+v", cloudwatchMetricData)

	// Log only the allocatable CPU and its corresponding timestamp
	for _, data := range result.RawData {
		// log.Println(data)
		log.Printf("Timestamp: %v, Allocatable CPU: %v", data.Timestamp, data.AllocatableCPU)
	}
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetAllocatableCPUMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_cpu_limit"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60), 
					Stat:   aws.String("Average"),
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
						MetricName: aws.String("node_cpu_reserved_capacity"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60), 
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
	// log.Println (result)
	return result, nil
}

func processCPURawData(result *cloudwatch.GetMetricDataOutput) AllocateResult {
	var rawData AllocateResult
	rawData.RawData = make([]TimeSeriesData, len(result.MetricDataResults[0].Timestamps))

	// Assuming the two metrics have the same number of data points
	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		cpuLimit := *result.MetricDataResults[0].Values[i]
		reservedCapacity := *result.MetricDataResults[1].Values[i]

		// Log the values for troubleshooting
		// log.Printf("Timestamp: %v, cpuLimit: %v, reservedCapacity: %v", *timestamp, cpuLimit, reservedCapacity)

		allocatableCPU := cpuLimit - reservedCapacity
		// log.Println(allocatableCPU)

		// Only include the calculated allocatable CPU in the result
		rawData.RawData[i].AllocatableCPU = allocatableCPU
	}
	// log.Println(rawData)
	return rawData
}

func init() {
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSAllocatableCpuCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
