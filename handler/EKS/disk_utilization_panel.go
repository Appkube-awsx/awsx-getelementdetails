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

type DiskUtilizationResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxEKSDiskUtilizationCmd = &cobra.Command{
	Use:   "cpu_utilization_panel",
	Short: "get cpu utilization metrics data",
	Long:  `command to get cpu utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetEKScpuUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetDiskUtilizationData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	rawSizeData, err := GetDiskSizeMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw size data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RawSizeData"] = rawSizeData

	rawUtilizationData, err := GetDiskUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw utilization data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RawUtilizationData"] = rawUtilizationData

	result := processDiskUtilizationData(rawSizeData, rawUtilizationData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetDiskSizeMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_filesystem_size"), 
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

func GetDiskUtilizationMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
						MetricName: aws.String("node_filesystem_utilization"), 
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

func processDiskUtilizationData(sizeResult, utilizationResult *cloudwatch.GetMetricDataOutput) DiskUtilizationResult {
	var rawData DiskUtilizationResult
	rawData.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(sizeResult.MetricDataResults[0].Timestamps))

	for i, timestamp := range sizeResult.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		size := *sizeResult.MetricDataResults[0].Values[i]
		utilization := *utilizationResult.MetricDataResults[0].Values[i]
		totalDiskSpaceUsed := size * utilization / 100
		rawData.RawData[i].Value = totalDiskSpaceUsed
	}

	return rawData
}

func init() {
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSDiskUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}