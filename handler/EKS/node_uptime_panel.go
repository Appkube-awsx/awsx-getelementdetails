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

type NodeUptimeDataPoint struct {
	Timestamp  time.Time `json:"Timestamp"`
	NodeUptime float64   `json:"NodeUptime"`
}

var AwsxEKSNodeUptimeCmd = &cobra.Command{
	Use:   "node_uptime_panel",
	Short: "get node uptime metrics data",
	Long:  `command to get node uptime metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNodeUptimePanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Node uptime data: ", err)
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

func GetNodeUptimePanel(cmd *cobra.Command, clientAuth *model.Auth,cloudWatchClient *cloudwatch.CloudWatch) (string, []NodeUptimeDataPoint, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	// elementType, _ := cmd.PersistentFlags().GetString("elementType")
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

	// Get node metrics
	nodeMetrics, err := GetNodeMetrics(clientAuth, instanceId, startTime, endTime,cloudWatchClient)
	if err != nil {
		log.Println("Error in getting node metrics: ", err)
		return "", nil, err
	}

	// Calculate node uptime data points
	var uptimeData []NodeUptimeDataPoint
	for i := 0; i < len(nodeMetrics.MetricDataResults[0].Values); i++ {
		uptime := 0.0
		if *nodeMetrics.MetricDataResults[0].Values[i] > 0 {
			uptime = 1.0
		}
		dataPoint := NodeUptimeDataPoint{
			Timestamp:  *nodeMetrics.MetricDataResults[0].Timestamps[i],
			NodeUptime: uptime,
		}
		uptimeData = append(uptimeData, dataPoint)
	}

	jsonString, err := json.Marshal(uptimeData)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), uptimeData, nil
}

func GetNodeMetrics(clientAuth *model.Auth, instanceId string, startTime, endTime *time.Time,cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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

func init() {
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeUptimeCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}