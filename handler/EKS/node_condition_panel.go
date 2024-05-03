package EKS

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
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

var AwsxEKSNodeConditionCmd = &cobra.Command{
	Use:   "node_condition_panel",
	Short: "get node condition metrics data",
	Long:  `command to get node condition metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNodeConditionPanel(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting Node condition data: ", err)
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

func GetNodeConditionPanel(cmd *cobra.Command, clientAuth *model.Auth) (map[string]float64, *NodeConditionPanel, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	// Get node condition data
	nodeConditionData, err := GetNodeConditionData(clientAuth, instanceId, startTime, endTime)
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

func GetNodeConditionData(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time) ([]*cloudwatch.MetricDataResult, error) {
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
								Value: aws.String(instanceId),
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
								Value: aws.String(instanceId),
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
								Name:  aws.String("instanceId"),
								Value: aws.String(instanceId),
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

func init() {
	AwsxEKSNodeConditionCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeConditionCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
