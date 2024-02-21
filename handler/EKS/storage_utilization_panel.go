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

type StorageUtilizationResult struct {
	RootVolumeUsage float64 `json:"rootVolumeUsage"`
	EBSVolume1Usage float64 `json:"ebsVolume1Usage"`
	EBSVolume2Usage float64 `json:"ebsVolume2Usage"`
}

var AwsxEKSStorageUtilizationCmd = &cobra.Command{
	Use:   "storage_utilization_panel",
	Short: "get storage utilization metrics data",
	Long:  `command to get storage utilization metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetStorageUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting storage utilization: ", err)
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

func GetStorageUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Parse start time if provided
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

	// Get Root Volume Usage
	rootVolumeUsage, err := GetStorageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "node_filesystem_utilization", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting root volume usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolumeUsage"] = rootVolumeUsage

	// Get EBS Volume 1 Usage
	ebsVolume1Usage, err := GetStorageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "node_filesystem_utilization", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS volume 1 usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBSVolume1Usage"] = ebsVolume1Usage

	// Get EBS Volume 2 Usage
	ebsVolume2Usage, err := GetStorageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "node_filesystem_utilization", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS volume 2 usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBSVolume2Usage"] = ebsVolume2Usage

	// Create JSON output
	jsonOutput := StorageUtilizationResult{
		RootVolumeUsage: *rootVolumeUsage.MetricDataResults[0].Values[0],
		EBSVolume1Usage: *ebsVolume1Usage.MetricDataResults[0].Values[0],
		EBSVolume2Usage: *ebsVolume2Usage.MetricDataResults[0].Values[0],
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetStorageMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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
							// Add dimensions for specific EBS volumes if needed
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
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
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSStorageUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}