package NLB

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

type UnhealthyHostCountData struct {
	UnhealthyHostCount []struct {
		Timestamp time.Time
		Value     float64
	} `json:"UnhealthyHostCount"`
}

var AwsxNLBUnhealthyHostCountCmd = &cobra.Command{
	Use:   "nlb_unhealthy_host_count_panel",
	Short: "Get NLB unhealthy host count metrics data",
	Long:  `Command to get NLB unhealthy host count metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNLBUnhealthyHostCountPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB unhealthy host count: ", err)
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

func GetNLBUnhealthyHostCountPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	//nlbArn, _ := cmd.PersistentFlags().GetString("nlbArn")
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	//elementType, _ := cmd.PersistentFlags().GetString("elementType")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")

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
	rawData, err := GetNLBUnhealthyHostCountMetricData(clientAuth, instanceId, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB unhealthy host count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["UnhealthyHostCount"] = rawData

	result := processNLBUnhealthyHostCountRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNLBUnhealthyHostCountMetricData(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for NLB %s from %v to %v", instanceId, startTime, endTime)

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
								Value: aws.String(instanceId),
							},
							{
								Name:  aws.String("TargetGroup"),
								Value: aws.String("targetgroup/k8s-istiosys-istioing-30129717de/b5e55c2955f8e65f"),
							},
						},
						MetricName: aws.String("UnHealthyHostCount"),
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

func processNLBUnhealthyHostCountRawData(result *cloudwatch.GetMetricDataOutput) UnhealthyHostCountData {
	var rawData UnhealthyHostCountData
	rawData.UnhealthyHostCount = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.UnhealthyHostCount[i].Timestamp = *timestamp
		rawData.UnhealthyHostCount[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxNLBUnhealthyHostCountCmd.PersistentFlags().String("instanceId", "", "instanceId")
	AwsxNLBUnhealthyHostCountCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBUnhealthyHostCountCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNLBUnhealthyHostCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
