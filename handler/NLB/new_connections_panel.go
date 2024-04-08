package NLB

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

type NewConnectionsData struct {
	NewConnections []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NewConnections"`
}

var AwsxNLBNewConnectionsCmd = &cobra.Command{
	Use:   "nlb_new_connections_panel",
	Short: "Get NLB new connections metrics data",
	Long:  `Command to get NLB new connections metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command..")
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
			jsonResp, cloudwatchMetricResp, err := GetNLBNewConnectionsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB new connections: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// Default case. It prints JSON
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNLBNewConnectionsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	nlbArn, _ := cmd.PersistentFlags().GetString("nlbArn")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

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

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := GetNLBNewConnectionsMetricData(clientAuth, nlbArn, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB new connections data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NewConnections"] = rawData

	result := processNLBNewConnectionsRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNLBNewConnectionsMetricData(clientAuth *model.Auth, nlbArn string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for NLB %s from %v to %v", nlbArn, startTime, endTime)

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
								Name:  aws.String("LoadBalancer"),
								Value: aws.String("net/a0affec9643ca40c5a4e837eab2f07fb/f623f27b6210158f"),
							},
						},
						MetricName: aws.String("NewFlowCount"),
						Namespace:  aws.String("AWS/NetworkELB"),
					},
					Period: aws.Int64(60),
					Stat:   aws.String(statistic),
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

func processNLBNewConnectionsRawData(result *cloudwatch.GetMetricDataOutput) NewConnectionsData {
	var rawData NewConnectionsData
	rawData.NewConnections = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.NewConnections[i].Timestamp = *timestamp
		rawData.NewConnections[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxNLBNewConnectionsCmd.PersistentFlags().String("nlbArn", "", "NLB ARN")
	AwsxNLBNewConnectionsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBNewConnectionsCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNLBNewConnectionsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
