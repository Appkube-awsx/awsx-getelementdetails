package ECS

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

type StorageResult struct {
	RootVolumeUtilization float64 `json:"RootVolumeUsage"`
	EBS1VolumeUtilization float64 `json:"EBSVolume1Usage"`
	EBS2VolumeUtilization float64 `json:"EBSVolume2Usage"`
}

var AwsxECSStorageUtilizationCmd = &cobra.Command{
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
	//startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	//endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

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

	// Get Root Volume Utilization
	rootVolumeUsage, err := GetStorageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "EphemeralStorageUtilized", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting root volume usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolumeUtilization"] = rootVolumeUsage

	// Get EBS1 Volume  Utilization
	ebs1VolumeUsage, err := GetStorageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "EphemeralStorageUtilized", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS1 Volume Utilization : ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBS1Volume1Utilization"] = ebs1VolumeUsage

	// Get EBS2 Volume Utilization
	ebs2VolumeUsage, err := GetStorageMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "EphemeralStorageUtilized", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS2 volume 2 usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBS2VolumeUtilization"] = ebs2VolumeUsage
	// Calculate average of all three volumes
	rootVolumeAvg := calculateAverage(rootVolumeUsage)
	ebs1VolumeAvg := calculateAverage(ebs1VolumeUsage) / 2 // Divide by 2
	ebs2VolumeAvg := calculateAverage(ebs2VolumeUsage) / 2 // Divide by 2

	// Format average utilizations to have two decimal places
	rootVolumeAvgStr := strconv.FormatFloat(rootVolumeAvg, 'f', 2, 64)
	ebs1VolumeAvgStr := strconv.FormatFloat(ebs1VolumeAvg, 'f', 2, 64)
	ebs2VolumeAvgStr := strconv.FormatFloat(ebs2VolumeAvg, 'f', 2, 64)

	// Convert formatted strings back to float64
	rootVolumeAvgFloat, err := strconv.ParseFloat(rootVolumeAvgStr, 64)
	if err != nil {
		log.Println("Error converting string to float64: ", err)
		return "", nil, err
	}
	ebs1VolumeAvgFloat, err := strconv.ParseFloat(ebs1VolumeAvgStr, 64)
	if err != nil {
		log.Println("Error converting string to float64: ", err)
		return "", nil, err
	}
	ebs2VolumeAvgFloat, err := strconv.ParseFloat(ebs2VolumeAvgStr, 64)
	if err != nil {
		log.Println("Error converting string to float64: ", err)
		return "", nil, err
	}
	// Create JSON output
	averageStorageResult := StorageResult{
		RootVolumeUtilization: rootVolumeAvgFloat,
		EBS1VolumeUtilization: ebs1VolumeAvgFloat,
		EBS2VolumeUtilization: ebs2VolumeAvgFloat,
	}

	jsonString, err := json.Marshal(averageStorageResult)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetStorageMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, statistic, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	//log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "ECS/ContainerInsights"
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

func calculateAverage(result *cloudwatch.GetMetricDataOutput) float64 {
	sum := 0.0
	if len(result.MetricDataResults) > 0 && len(result.MetricDataResults[0].Values) > 0 {
		for _, value := range result.MetricDataResults[0].Values {
			sum += *value
		}
		return sum / float64(len(result.MetricDataResults[0].Values))
	}
	return 0
}

func init() {
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSStorageUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
