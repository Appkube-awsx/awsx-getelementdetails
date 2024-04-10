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

type HealthyHostCountData struct {
	HealthyHostCount []struct {
		Timestamp time.Time
		Value     float64
	} `json:"HealthyHostCount"`
}

var AwsxNLBHealthyHostCountCmd = &cobra.Command{
	Use:   "nlb_healthy_host_count_panel",
	Short: "Get NLB healthy host count metrics data",
	Long:  `Command to get NLB healthy host count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNLBHealthyHostCountPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB healthy host count: ", err)
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

func GetNLBHealthyHostCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := GetNLBHealthyHostCountMetricData(clientAuth, nlbArn, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB healthy host count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["HealthyHostCount"] = rawData

	result := processNLBHealthyHostCountRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNLBHealthyHostCountMetricData(clientAuth *model.Auth, nlbArn string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
							{
								Name:  aws.String("TargetGroup"),
								Value: aws.String("targetgroup/k8s-istiosys-istioing-30129717de/b5e55c2955f8e65f"),
							},
						},
						MetricName: aws.String("HealthyHostCount"),
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

func processNLBHealthyHostCountRawData(result *cloudwatch.GetMetricDataOutput) HealthyHostCountData {
	var rawData HealthyHostCountData
	rawData.HealthyHostCount = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.HealthyHostCount[i].Timestamp = *timestamp
		rawData.HealthyHostCount[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("nlbArn", "", "NLB ARN")
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNLBHealthyHostCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
