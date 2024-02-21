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

type NodeCapacityMetrics struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	StorageAvail float64 `json:"storage_avail"`
}

const (
	NodeCPUMetricName     = "node_cpu_utilization"
	NodeMemoryMetricName  = "node_memory_utilization"
	NodeStorageMetricName = "node_filesystem_utilization"
)

type NodeCapacityPanel struct {
	RawData  map[string]*cloudwatch.GetMetricDataOutput `json:"raw_data"`
	JsonData string                                     `json:"json_data"`
}

var AwsxEKSNodeCapacityCmd = &cobra.Command{
	Use:   "node_capacity_panel",
	Short: "get node capacity metrics data",
	Long:  `command to get node capacity metrics data`,

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
			nodeCapacityPanel, err := GetNodeCapacityPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting node capacity data: ", err)
				return
			}

			jsonResp := nodeCapacityPanel.JsonData
			cloudwatchMetricResp := nodeCapacityPanel.RawData

			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNodeCapacityPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (*NodeCapacityPanel, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	// elementType, _ := cmd.PersistentFlags().GetString("elementType")
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
			return nil, err
		}
		instanceId = cmdbData.InstanceId

	}

	startTime, endTime := NodeCapacityparseTime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cpuUsageRawData, err := GetNodeCapacityMetricData(clientAuth, instanceId, "", startTime, endTime, NodeCPUMetricName, cloudWatchClient)
	if err != nil {
		return nil, err
	}

	memoryUsageRawData, err := GetNodeCapacityMetricData(clientAuth, instanceId, "", startTime, endTime, NodeMemoryMetricName, cloudWatchClient)
	if err != nil {
		return nil, err
	}

	storageAvailRawData, err := GetNodeCapacityMetricData(clientAuth, instanceId, "", startTime, endTime, NodeStorageMetricName, cloudWatchClient)
	if err != nil {
		return nil, err
	}

	totalCPU := 100.0 // Assuming 100% CPU
	totalMemory := 100.0
	totalStorage := 100.0

	nodeCapacity := NodeCapacityMetrics{
		CPUUsage:     calculateCPUUsage(cpuUsageRawData, totalCPU),
		MemoryUsage:  calculateMemoryUsage(memoryUsageRawData, totalMemory),
		StorageAvail: calculateStorageAvailability(storageAvailRawData, totalStorage),
	}

	jsonData, err := json.Marshal(nodeCapacity)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return nil, err
	}

	nodeCapacityPanel := &NodeCapacityPanel{
		RawData:  map[string]*cloudwatch.GetMetricDataOutput{"cpu": cpuUsageRawData, "memory": memoryUsageRawData, "storage": storageAvailRawData},
		JsonData: string(jsonData),
	}

	return nodeCapacityPanel, nil
}

func calculateCPUUsage(data *cloudwatch.GetMetricDataOutput, totalCPU float64) float64 {
	// Sum up CPU usage values
	var sum float64
	for _, result := range data.MetricDataResults {
		for _, value := range result.Values {
			sum += *value
		}
	}
	// Calculate average CPU usage percentage
	return (sum / float64(len(data.MetricDataResults))) / totalCPU
}

func calculateMemoryUsage(data *cloudwatch.GetMetricDataOutput, totalMemory float64) float64 {
	// Sum up memory usage values
	var sum float64
	for _, result := range data.MetricDataResults {
		for _, value := range result.Values {
			sum += *value
		}
	}
	return (sum / float64(len(data.MetricDataResults))) / totalMemory
}

func calculateStorageAvailability(data *cloudwatch.GetMetricDataOutput, totalStorage float64) float64 {
	var sum float64
	for _, result := range data.MetricDataResults {
		for _, value := range result.Values {
			sum += *value
		}
	}
	// Calculate average storage availability percentage
	return (sum / float64(len(data.MetricDataResults))) / totalStorage
}

func GetNodeCapacityMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String(metricName),
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

func NodeCapacityparseTime(startTimeStr, endTimeStr string) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		} else {
			startTime = &parsedStartTime
		}
	} else {
		now := time.Now()
		startTime = &now
		minusFiveMinutes := now.Add(-5 * time.Minute)
		startTime = &minusFiveMinutes
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
		} else {
			endTime = &parsedEndTime
		}
	} else {
		now := time.Now()
		endTime = &now
	}

	return startTime, endTime
}

func init() {
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}