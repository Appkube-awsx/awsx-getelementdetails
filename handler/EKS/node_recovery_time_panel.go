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

var AwsxEKSNodeRecoveryPanelCmd = &cobra.Command{
	Use:   "node_recovery_time_panel",
	Short: "get node recovery time metrics data",
	Long:  `command to get node recovery time metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNodeRecoveryTime(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting node recovery time data: ", err)
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

type NodeRecoveryData struct {
	Timestamp    time.Time     `json:"timestamp"`
	RecoveryTime time.Duration `json:"recovery_time"`
}

func GetNodeRecoveryTime(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, []NodeRecoveryData, error) {
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

	// Fetch node ready metric data
	nodeReadyData, err := GetNodeReadyMetricData(clientAuth, instanceId, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error fetching node ready metric data: ", err)
		return "", nil, err
	}

	// Process node ready data
	recoveryTimeSeries := ProcessNodeReadyData(nodeReadyData)

	// Check if recoveryTimeSeries is empty
	if len(recoveryTimeSeries) == 0 {
		return "No node recovery events detected", nil, nil
	}

	// Marshal recovery time series data to JSON
	jsonData, err := json.Marshal(recoveryTimeSeries)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return "", nil, err
	}

	return string(jsonData), recoveryTimeSeries, nil
}

func GetNodeReadyMetricData(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_status_condition_ready"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Maximum"),
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

func ProcessNodeReadyData(result *cloudwatch.GetMetricDataOutput) []NodeRecoveryData {
	var recoveryTimeSeries []NodeRecoveryData

	for i := 1; i < len(result.MetricDataResults[0].Timestamps); i++ {
		currentTimestamp := *result.MetricDataResults[0].Timestamps[i]
		previousTimestamp := *result.MetricDataResults[0].Timestamps[i-1]

		if *result.MetricDataResults[0].Values[i-1] == 0 && *result.MetricDataResults[0].Values[i] == 1 {
			recoveryTime := currentTimestamp.Sub(previousTimestamp)

			recoveryTimeSeries = append(recoveryTimeSeries, NodeRecoveryData{
				Timestamp:    currentTimestamp,
				RecoveryTime: recoveryTime,
			})
		}
	}

	return recoveryTimeSeries
}
