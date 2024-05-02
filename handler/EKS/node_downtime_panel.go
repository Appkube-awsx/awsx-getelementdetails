package EKS

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NodeDowntimeDataPoint struct {
	Timestamp    time.Time `json:"Timestamp"`
	NodeDowntime float64   `json:"NodeDowntime"`
}

var AwsxEKSNodeDowntimeCmd = &cobra.Command{
	Use:   "node_downtime_panel",
	Short: "get node downtime metrics data",
	Long:  `command to get node downtime metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNodeDowntimePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Node downtime data: ", err)
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

func GetNodeDowntimePanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, []NodeDowntimeDataPoint, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	nodeMetrics, err := GetNodeDowntimeMetrics(clientAuth, instanceId, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting node metrics: ", err)
		return "", nil, err
	}

	// Calculate node downtime data points
	var downtimeData []NodeDowntimeDataPoint
	for i := 0; i < len(nodeMetrics.MetricDataResults[0].Values); i++ {
		downtime := 0.0
		if *nodeMetrics.MetricDataResults[0].Values[i] <= 0 {
			downtime = 1.0
		}
		dataPoint := NodeDowntimeDataPoint{
			Timestamp:    *nodeMetrics.MetricDataResults[0].Timestamps[i],
			NodeDowntime: downtime,
		}
		downtimeData = append(downtimeData, dataPoint)
	}

	jsonString, err := json.Marshal(downtimeData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), downtimeData, nil
}

func GetNodeDowntimeMetrics(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	elmType := "ContainerInsights"
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("cpu_utilization"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("node_cpu_utilization"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),

					Stat: aws.String("Average"),
				},
			},
			{
				Id: aws.String("memory_utilization"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("node_memory_utilization"),
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

	return result, nil
}
