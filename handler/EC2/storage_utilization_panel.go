package EC2

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/Appkube-awsx/awsx-common/authenticate"
// 	"github.com/Appkube-awsx/awsx-common/awsclient"
// 	"github.com/Appkube-awsx/awsx-common/cmdb"
// 	"github.com/Appkube-awsx/awsx-common/config"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/spf13/cobra"
// )

// type volumeUsage struct {
// 	Value float64 `json:"value"`
// 	Unit  string  `json:"unit"`
// }

// type volumeMetrics struct {
// 	RootVolumeUsage volumeUsage `json:"RootVolumeUsage"`
// 	EBSVolume1Usage volumeUsage `json:"EBSVolume1Usage"`
// 	EBSVolume2Usage volumeUsage `json:"EBSVolume2Usage"`
// }

// var AwsxEc2StorageUtilizationCmd = &cobra.Command{
// 	Use:   "storage_utilization_panel",
// 	Short: "get storage utilization metrics data",
// 	Long:  `command to get storage utilization metrics data`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("running from child command")
// 		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
// 		if err != nil {
// 			log.Printf("Error during authentication: %v\n", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return
// 			}
// 			return
// 		}
// 		if authFlag {
// 			responseType, _ := cmd.PersistentFlags().GetString("responseType")
// 			jsonResp, cloudwatchMetricResp, err := GetVolumeMetricPanel(cmd, clientAuth, nil)
// 			if err != nil {
// 				log.Println("Error getting storage utilization: ", err)
// 				return
// 			}
// 			if responseType == "frame" {
// 				fmt.Println(cloudwatchMetricResp)
// 			} else {
// 				// default case. it prints json
// 				fmt.Println(jsonResp)
// 			}
// 		}

// 	},
// }

// func GetVolumeMetricPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
// 	elementId, _ := cmd.PersistentFlags().GetString("elementId")
// 	elementType, _ := cmd.PersistentFlags().GetString("elementType")
// 	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
// 	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
// 	RootVolumeId, _ := cmd.PersistentFlags().GetString("RootVolumeId")
// 	EBSVolume1Id, _ := cmd.PersistentFlags().GetString("EBSVolume1Id")
// 	EBSVolume2Id, _ := cmd.PersistentFlags().GetString("EBSVolume2Id")

// 	if elementId != "" {
// 		log.Println("getting cloud-element data from cmdb")
// 		apiUrl := cmdbApiUrl
// 		if cmdbApiUrl == "" {
// 			log.Println("using default cmdb url")
// 			apiUrl = config.CmdbUrl
// 		}
// 		log.Println("cmdb url: " + apiUrl)
// 		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 		instanceId = cmdbData.InstanceId

// 	}

// 	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
// 	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

// 	var startTime, endTime *time.Time

// 	// Parse start time if provided
// 	if startTimeStr != "" {
// 		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing start time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", nil, err
// 			}
// 			return "", nil, err
// 		}
// 		startTime = &parsedStartTime
// 	} else {
// 		defaultStartTime := time.Now().Add(-5 * time.Minute)
// 		startTime = &defaultStartTime
// 	}

// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing end time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", nil, err
// 			}
// 			return "", nil, err
// 		}
// 		endTime = &parsedEndTime
// 	} else {
// 		defaultEndTime := time.Now()
// 		endTime = &defaultEndTime
// 	}
// 	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

// 	// Get metrics for root volume
// 	rootVolumeMetrics, err := GetMetricData(clientAuth, instanceId, RootVolumeId, elementType, startTime, endTime, "VolumeStalledIOCheck", "SampleCount", cloudWatchClient)
// 	if err != nil {
// 		log.Println("Error in getting metrics for root volume: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["RootVolume"] = rootVolumeMetrics

// 	// Get metrics for EBS1 volume
// 	ebsVolume1Metrics, err := GetMetricData(clientAuth, instanceId, EBSVolume1Id, elementType, startTime, endTime, "VolumeStalledIOCheck", "SampleCount", cloudWatchClient)
// 	if err != nil {
// 		log.Println("Error in getting metrics for EBS1 volume: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["EBSVolume1"] = ebsVolume1Metrics

// 	// Get metrics for EBS2 volume
// 	ebsVolume2Metrics, err := GetMetricData(clientAuth, instanceId, EBSVolume2Id, elementType, startTime, endTime, "VolumeStalledIOCheck", "SampleCount", cloudWatchClient)
// 	if err != nil {
// 		log.Println("Error in getting metrics for EBS2 volume: ", err)
// 		return "", nil, err
// 	}
// 	cloudwatchMetricData["EBSVolume2"] = ebsVolume2Metrics

// 	// JSON output for volume metrics
// 	var volumeMetricsOutput volumeMetrics

// 	if len(rootVolumeMetrics.MetricDataResults) >= 3 &&
// 		len(rootVolumeMetrics.MetricDataResults[0].Values) >= 1 &&
// 		len(rootVolumeMetrics.MetricDataResults[1].Values) >= 1 &&
// 		len(rootVolumeMetrics.MetricDataResults[2].Values) >= 1 {
// 		volumeMetricsOutput = volumeMetrics{
// 			RootVolumeUsage: volumeUsage{
// 				Value: *rootVolumeMetrics.MetricDataResults[0].Values[0],
// 				Unit:  "GB",
// 			},
// 			EBSVolume1Usage: volumeUsage{
// 				Value: *rootVolumeMetrics.MetricDataResults[1].Values[0],
// 				Unit:  "GB",
// 			},
// 			EBSVolume2Usage: volumeUsage{
// 				Value: *rootVolumeMetrics.MetricDataResults[2].Values[0],
// 				Unit:  "GB",
// 			},
// 		}
// 	} else {
// 		log.Println("Error: Not enough data in MetricDataResults.")
// 	}

// 	jsonString, err := json.Marshal(volumeMetricsOutput)
// 	if err != nil {
// 		log.Println("Error in marshalling volume metrics json in string: ", err)
// 		return "", nil, err
// 	}

// 	return string(jsonString), cloudwatchMetricData, nil
// }

// func GetMetricData(clientAuth *model.Auth, instanceID, elementType string, volumeID string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
// 	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)
// 	var metricDataQueries []*cloudwatch.MetricDataQuery

// 	for i, metricName := range metrics {
// 		query := &cloudwatch.MetricDataQuery{
// 			Id: aws.String(fmt.Sprintf("m%d", i+1)),
// 			MetricStat: &cloudwatch.MetricStat{
// 				Metric: &cloudwatch.Metric{
// 					Dimensions: []*cloudwatch.Dimension{
// 						{
// 							Name:  aws.String("InstanceId"),
// 							Value: aws.String(instanceID),
// 						},
// 						{
// 							Name:  aws.String("VolumeId"),
// 							Value: aws.String(volumeID),
// 						},
// 					},
// 					MetricName: aws.String(metricName),
// 					Namespace:  aws.String("AWS/" + elementType),
// 				},
// 				Period: aws.Int64(300),
// 				Stat:   aws.String("SampleCount"), // You can customize this if needed
// 			},
// 		}
// 		metricDataQueries = append(metricDataQueries, query)
// 	}

// 	input := &cloudwatch.GetMetricDataInput{
// 		EndTime:           endTime,
// 		StartTime:         startTime,
// 		MetricDataQueries: metricDataQueries,
// 	}

// 	if cloudWatchClient == nil {
// 		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
// 	}

// 	result, err := cloudWatchClient.GetMetricData(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
// func init() {
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("query", "", "query")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("endTime", "", "end time")
// 	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
// }
