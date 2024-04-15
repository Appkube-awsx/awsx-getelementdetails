package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

type Results struct {
	CurrentUsage float64 `json:"currentUsage"`
	AverageUsage float64 `json:"averageUsage"`
	MaxUsage     float64 `json:"maxUsage"`
}

var AwsxRDSMemoryUtilizationCmd = &cobra.Command{
	Use:   "memory_utilization_panel",
	Short: "get memory utilization metrics data",
	Long:  `command to get memory utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetRDSMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetRDSMemoryUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-15 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	// Retrieve instance class of the RDS instance
	instanceClass, err := GetRDSInstanceClass(clientAuth)
	if err != nil {
		log.Println("Error getting RDS instance class: ", err)
		return "", nil, err
	}

	// Determine total memory based on instance class
	totalMemoryBytes := GetTotalMemory(instanceClass)
	totalMemoryGB := convertBytesToGB(float64(totalMemoryBytes))
	fmt.Println("Total Memory (GB):", totalMemoryGB)

	// Fetch CloudWatch metric data for current, average, and maximum memory usage
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	// Get current usage
	currentUsage, err := GetRDSMemoryUtilizationMetricData(clientAuth, startTime, endTime, "SampleCount", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting current usage: ", err)
		return "", nil, err
	}
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CurrentUsage"] = currentUsage
	} else {
		log.Println("No data available for current usage")
	}

	// Get average usage
	averageUsage, err := GetRDSMemoryUtilizationMetricData(clientAuth, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average usage: ", err)
		return "", nil, err
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["AverageUsage"] = averageUsage
	} else {
		log.Println("No data available for average usage")
	}

	// Get maximum usage
	maxUsage, err := GetRDSMemoryUtilizationMetricData(clientAuth, startTime, endTime, "Maximum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting maximum usage: ", err)
		return "", nil, err
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["MaxUsage"] = maxUsage
	} else {
		log.Println("No data available for maximum usage")
	}

	// Create JSON output
	jsonOutput := make(map[string]float64)
	if len(currentUsage.MetricDataResults) > 0 && len(currentUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["CurrentUsage"] = convertBytesToGB(*currentUsage.MetricDataResults[0].Values[0])
	}
	if len(averageUsage.MetricDataResults) > 0 && len(averageUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["AverageUsage"] = convertBytesToGB(*averageUsage.MetricDataResults[0].Values[0])
	}
	if len(maxUsage.MetricDataResults) > 0 && len(maxUsage.MetricDataResults[0].Values) > 0 {
		jsonOutput["MaxUsage"] = convertBytesToGB(*maxUsage.MetricDataResults[0].Values[0])
	}

	// Convert JSON output to string
	jsonResult, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", nil, err
	}

	return string(jsonResult), cloudwatchMetricData, nil
}

func GetRDSMemoryUtilizationMetricData(clientAuth *model.Auth, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	// Get metric data for memory utilization
	input := &cloudwatch.GetMetricDataInput{
		StartTime: startTime,
		EndTime:   endTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("FreeableMemory"),
						Namespace:  aws.String("AWS/RDS"),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
				},
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}
	// Retrieve metric data
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("Memory Utilization Data:", result)

	return result, nil
}

func GetRDSInstanceClass(clientAuth *model.Auth) (string, error) {
	// Initialize RDS client
	rdsClient := awsclient.GetClient(*clientAuth, awsclient.RDS_CLIENT).(*rds.RDS)

	// Get DB instance details
	input := &rds.DescribeDBInstancesInput{}

	result, err := rdsClient.DescribeDBInstances(input)
	if err != nil {
		return "", err
	}

	// Assuming a single RDS instance for simplicity, extract the instance class
	instanceClass := *result.DBInstances[0].DBInstanceClass
	fmt.Println("Instance Class:", instanceClass)

	return instanceClass, nil
}

func GetTotalMemory(instanceClass string) int64 {
	// Determine total memory based on the instance class
	// You can create a mapping between instance classes and their corresponding memory sizes.
	// For example:
	switch instanceClass {
	case "db.t4g.medium":
		return 4 * 1024 * 1024 * 1024 // 4 GB
	case "db.t3.medium":
		return 2 * 1024 * 1024 * 1024 // 2 GB
	// Add cases for other instance classes as needed
	default:
		return 0 // Default value if instance class is unknown
	}
}

func convertBytesToGB(bytes float64) float64 {
	// Convert memory data from bytes to GB
	return bytes / (1024 * 1024 * 1024)
}

func init() {
	AwsxRDSMemoryUtilizationCmd.PersistentFlags().String("startTime", "", "Start time for metrics retrieval")
	AwsxRDSMemoryUtilizationCmd.PersistentFlags().String("endTime", "", "End time for metrics retrieval")
}
