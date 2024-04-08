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

type ActiveConnectionsData struct {
	ActiveConnections []struct {
		Timestamp time.Time
		Value     float64
	} `json:"ActiveConnections"`
}

var AwsxNLBActiveConnectionsCmd = &cobra.Command{
	Use:   "nlb_active_connections_panel",
	Short: "Get NLB active connections metrics data",
	Long:  `Command to get NLB active connections metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNLBActiveConnectionsPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting NLB active connections: ", err)
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

func GetNLBActiveConnectionsPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

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
	rawData, err := GetNLBActiveConnectionsMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NLB active connections data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["ActiveConnections"] = rawData

	result := processNLBActiveConnectionsRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNLBActiveConnectionsMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for NLB %s from %v to %v %v", instanceId, elementType, startTime, endTime)

	elmType := "AWS/NetworkELB"

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
						MetricName: aws.String("ActiveFlowCount"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"),
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

func processNLBActiveConnectionsRawData(result *cloudwatch.GetMetricDataOutput) ActiveConnectionsData {
	var rawData ActiveConnectionsData
	rawData.ActiveConnections = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.ActiveConnections[i].Timestamp = *timestamp
		rawData.ActiveConnections[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {

	AwsxNLBActiveConnectionsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxNLBActiveConnectionsCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxNLBActiveConnectionsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
