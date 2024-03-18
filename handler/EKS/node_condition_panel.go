package EKS

import (
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NodeConditionPanel struct {
	DiskPressureAvg   float64 `json:"disk_pressure_avg"`
	MemoryPressureAvg float64 `json:"memory_pressure_avg"`
	PIDPressureAvg    float64 `json:"pid_pressure_avg"`
}

func GetNodeConditionPanel(cmd *cobra.Command, clientAuth *model.Auth) (map[string]float64, *NodeConditionPanel, error) {
	clusterName, _ := cmd.PersistentFlags().GetString("clusterName")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return nil, nil, err
		}
		startTime = &parsedStartTime
	} else {
		// Use default start time if not provided
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	// Parse end time
	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return nil, nil, err
		}
		endTime = &parsedEndTime
	} else {
		// Use default end time if not provided
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Get node condition data
	nodeConditionData, err := GetNodeConditionData(clientAuth, clusterName, startTime, endTime)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pressure averages
	diskPressureAvg, memoryPressureAvg, pidPressureAvg := calculatePressureAverages(nodeConditionData)

	// Create NodeConditionPanel object
	nodeConditionPanel := &NodeConditionPanel{
		DiskPressureAvg:   diskPressureAvg,
		MemoryPressureAvg: memoryPressureAvg,
		PIDPressureAvg:    pidPressureAvg,
	}

	// Return map of field names and their corresponding values
	return map[string]float64{
		"disk_pressure":   diskPressureAvg,
		"memory_pressure": memoryPressureAvg,
		"pid_pressure":    pidPressureAvg,
	}, nodeConditionPanel, nil
}

func calculatePressureAverages(data []*cloudwatch.MetricDataResult) (float64, float64, float64) {
	var diskPressureTotal, memoryPressureTotal, pidPressureTotal, totalCount float64

	for _, result := range data {
		for _, value := range result.Values {
			if value != nil {
				switch *result.Id {
				case "diskPressure", "memoryPressure", "pidPressure":
					switch len(result.Values) {
					case 0:
						continue
					default:
						totalCount++
						switch *result.Id {
						case "diskPressure":
							diskPressureTotal += *value
						case "memoryPressure":
							memoryPressureTotal += *value
						case "pidPressure":
							pidPressureTotal += *value
						}
					}
				}
			}
		}
	}

	if totalCount == 0 {
		return 0, 0, 0
	}

	diskPressureAvg := diskPressureTotal / totalCount
	memoryPressureAvg := memoryPressureTotal / totalCount
	pidPressureAvg := pidPressureTotal / totalCount

	return diskPressureAvg, memoryPressureAvg, pidPressureAvg
}

func GetNodeConditionData(clientAuth *model.Auth, clusterName string, startTime, endTime *time.Time) ([]*cloudwatch.MetricDataResult, error) {
	// Define the metric names for disk pressure, memory pressure, and PID pressure
	diskPressureMetricName := "node_status_condition_disk_pressure"
	memoryPressureMetricName := "node_status_condition_memory_pressure"
	pidPressureMetricName := "node_status_condition_pid_pressure"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("diskPressure"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String(diskPressureMetricName),
						Namespace:  aws.String("ContainerInsights"), // Update with your namespace
					},
					Period: aws.Int64(300),        // Adjust the period as needed
					Stat:   aws.String("Average"), // Assuming you want average value
				},
			},
			{
				Id: aws.String("memoryPressure"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String(memoryPressureMetricName),
						Namespace:  aws.String("ContainerInsights"), // Update with your namespace
					},
					Period: aws.Int64(300),        // Adjust the period as needed
					Stat:   aws.String("Average"), // Assuming you want average value
				},
			},
			{
				Id: aws.String("pidPressure"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(clusterName),
							},
						},
						MetricName: aws.String(pidPressureMetricName),
						Namespace:  aws.String("ContainerInsights"), // Update with your namespace
					},
					Period: aws.Int64(300),        // Adjust the period as needed
					Stat:   aws.String("Average"), // Assuming you want average value
				},
			},
		},
	}

	// Get the CloudWatch client
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	// Call the GetMetricData API
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	// Return the metric data results
	return result.MetricDataResults, nil
}
