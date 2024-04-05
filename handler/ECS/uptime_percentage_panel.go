package ECS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type MetricResults struct {
	TimeSeriesData map[string]string `json:"timeSeriesData"`
}

var AwsxECSUptimeCmd = &cobra.Command{
	Use:   "ecs_uptime_panel",
	Short: "get uptime metrics data for ECS",
	Long:  `command to get uptime metrics data for ECS`,

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
			jsonResp, _, err := GetECSUptimeData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting ECS uptime data: ", err)
				return
			}
			fmt.Println(jsonResp)
		}
	},
}

func GetECSUptimeData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]string, error) {
	ClusterName := "cluster-01-02-2024"
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

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
	totalTaskCount, totalServiceCount, err := GetECSTaskAndServiceCount(clientAuth, startTime, endTime, ClusterName, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting ECS task and service count metrics: ", err)
		return "", nil, err
	}

	// Calculate uptime percentage
	uptimePercentage := (totalTaskCount / totalServiceCount) * 100
	if uptimePercentage > 100 {
		uptimePercentage = 100
	}

	timeSeriesData := map[string]string{
		"uptimePercentage": fmt.Sprintf("%.2f%%", uptimePercentage),
	}

	log.Printf("Uptime Percentage: %f", uptimePercentage)

	jsonString, err := json.Marshal(MetricResults{TimeSeriesData: timeSeriesData})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), timeSeriesData, nil
}

func GetECSTaskAndServiceCount(clientAuth *model.Auth, startTime, endTime *time.Time, ClusterName string, cloudWatchClient *cloudwatch.CloudWatch) (float64, float64, error) {
	taskCountInput := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("ECS/ContainerInsights"),
		MetricName: aws.String("TaskCount"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("ClusterName"),
				Value: aws.String(ClusterName),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Average")},
		Unit:       aws.String("Count"),
	}

	serviceCountInput := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("ECS/ContainerInsights"),
		MetricName: aws.String("ServiceCount"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("ClusterName"),
				Value: aws.String(ClusterName),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Average")},
		Unit:       aws.String("Count"),
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	taskCountResult, err := cloudWatchClient.GetMetricStatistics(taskCountInput)
	if err != nil {
		return 0, 0, err
	}

	serviceCountResult, err := cloudWatchClient.GetMetricStatistics(serviceCountInput)
	if err != nil {
		return 0, 0, err
	}

	totalTaskCount := 0.0
	totalServiceCount := 0.0

	for _, dp := range taskCountResult.Datapoints {
		totalTaskCount += *dp.Average
	}

	for _, dp := range serviceCountResult.Datapoints {
		totalServiceCount += *dp.Average
	}

	return totalTaskCount, totalServiceCount, nil
}

func init() {
	AwsxECSUptimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSUptimeCmd.PersistentFlags().String("endTime", "", "end time")
}
