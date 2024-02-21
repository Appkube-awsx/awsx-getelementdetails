package EKS

import (
	"encoding/json"
	"fmt"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"

	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type MemoryResult struct {
	CurrentUsage float64 `json:"currentUsage"`
	AverageUsage float64 `json:"averageUsage"`
	MaxUsage     float64 `json:"maxUsage"`
}

var AwsxEKSMemoryUtilizationCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GeteksMemoryUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting memory utilization: ", err)
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

func GeteksMemoryUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
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
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
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
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	currentUsage, err := GeteksContainerMetricData(clientAuth, instanceId, elementType, startTime, endTime, "SampleCount", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CurrentUsage"] = currentUsage
	averageUsage, err := GeteksContainerMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["AverageUsage"] = averageUsage
	maxUsage, err := GeteksContainerMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Maximum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting maximum: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["MaxUsage"] = maxUsage
	jsonOutput := map[string]float64{
		"CurrentUsage": *currentUsage.MetricDataResults[0].Values[0],
		"AverageUsage": *averageUsage.MetricDataResults[0].Values[0],
		"MaxUsage":     *maxUsage.MetricDataResults[0].Values[0],
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil

}

func GeteksContainerMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_memory_utilization"),
						Namespace:  aws.String(elmType),
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
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func init() {
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSMemoryUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}