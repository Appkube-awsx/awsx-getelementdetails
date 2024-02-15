package EKS

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
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

func GetNodeCapacityPanel(cmd *cobra.Command, clientAuth *model.Auth) (*NodeCapacityPanel, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	startTime, endTime := NodeCapacityparseTime(startTimeStr, endTimeStr)

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cpuUsageRawData, err := GetNodeCapacityMetricData(clientAuth, clusterName, "", startTime, endTime, NodeCPUMetricName)
	if err != nil {
		return nil, err
	}

	memoryUsageRawData, err := GetNodeCapacityMetricData(clientAuth, clusterName, "", startTime, endTime, NodeMemoryMetricName)
	if err != nil {
		return nil, err
	}

	storageAvailRawData, err := GetNodeCapacityMetricData(clientAuth, clusterName, "", startTime, endTime, NodeStorageMetricName)
	if err != nil {
		return nil, err
	}

	totalCPU := 100.0       // Assuming 100% CPU
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

func GetNodeCapacityMetricData(clientAuth *model.Auth, clusterName, namespace string, startTime, endTime *time.Time, metricName string) (*cloudwatch.GetMetricDataOutput, error) {
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
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String("ContainerInsights"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"), 
				},
			},
		},
	}

	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
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
