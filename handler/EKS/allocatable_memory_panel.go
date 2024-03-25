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

type TimeSeriesMemData struct {
	Timestamp      time.Time
	AllocatableMem float64
	// ReservedMem    float64
}

type AllocateMemResult struct {
	AllocatableMemory []TimeSeriesMemData `json:"AllocatableMemory"`
}

var AwsxEKSAllocatableMemCmd = &cobra.Command{
	Use:   "allocatable_mem_panel",
	Short: "get allocatable memory metrics data",
	Long:  `command to get allocatable memory metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetAllocatableMemData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting allocatable memory: ", err)
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

func GetAllocatableMemData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
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

	// Fetch raw data
	rawData, err := GetAllocatableMemMetricData(clientAuth, instanceId, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}

	// Process the raw data if needed
	result := processMemRawData(rawData)
	// cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
	// 	"AllocatableMemory": &cloudwatch.GetMetricDataOutput{
	// 		MetricDataResults: []*cloudwatch.MetricDataResult{
	// 			{
	// 				Timestamps: make([]*time.Time, len(result.AllocatableMemory)),
	// 				Values:     make([]*float64, len(result.AllocatableMemory)),
	// 			},
	// 		},
	// 	},
	// }

	// // Assign the processed data to cloudwatchMetricData
	// for i, data := range result.AllocatableMemory {
	// 	cloudwatchMetricData["AllocatableMemory"].MetricDataResults[0].Timestamps[i] = &data.Timestamp
	// 	cloudwatchMetricData["AllocatableMemory"].MetricDataResults[0].Values[i] = &data.AllocatableMem
	// }

	timestamps := make([]time.Time, len(result.AllocatableMemory))
	values := make([]float64, len(result.AllocatableMemory))

	fmt.Println("timeeeeeeeee", timestamps)
	fmt.Println("timeeeeeeeee", values)

	// Populate the slices with actual data
	for i, data := range result.AllocatableMemory {
		// Assigning values directly to slices without taking their addresses
		timestamps[i] = data.Timestamp
		values[i] = data.AllocatableMem
	}

	// Initialize the MetricDataResults slice
	metricDataResults := make([]*cloudwatch.MetricDataResult, len(result.AllocatableMemory))

	// Populate the MetricDataResults with actual data
	for i := range result.AllocatableMemory {
		metricDataResults[i] = &cloudwatch.MetricDataResult{
			Timestamps: []*time.Time{&timestamps[i]},
			Values:     []*float64{&values[i]},
		}
	}

	// Assign the processed data to cloudwatchMetricData
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"AllocatableMemory": {
			MetricDataResults: metricDataResults,
		},
	}

	// Log only the allocatable memory and its corresponding timestamp
	// for _, data := range result.AllocatableMemory {
	// 	log.Printf("Timestamp: %v, Allocatable Memory: %v", data.Timestamp, data.AllocatableMem)
	// }
	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetAllocatableMemMetricData(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_memory_limit"),
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
						MetricName: aws.String("node_memory_reserved_capacity"),
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
	// fmt.Println("result",result)
	// fmt.Println("instanceId",instanceId)
	// fmt.Println("elmType",elmType)
	// fmt.Println("input",input)


	return result, nil
}

func processMemRawData(result *cloudwatch.GetMetricDataOutput) AllocateMemResult {
	var rawData AllocateMemResult
	rawData.AllocatableMemory = make([]TimeSeriesMemData, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.AllocatableMemory[i].Timestamp = *timestamp
		memLimit := *result.MetricDataResults[0].Values[i]
		reservedCapacity := *result.MetricDataResults[1].Values[i]
		fmt.Println("memlimit",memLimit)
		fmt.Println("reserved capacity",reservedCapacity)
		allocatableMem := memLimit - reservedCapacity

		// Only include the calculated allocatable memory in the result
		rawData.AllocatableMemory[i].AllocatableMem = allocatableMem
	}
	// fmt.Println("raw data",rawData)
	return rawData
}
func init() {
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSAllocatableMemCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
