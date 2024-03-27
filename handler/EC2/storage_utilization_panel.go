package EC2

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

var AwsxEc2StorageUtilizationCmd = &cobra.Command{
	Use:   "Storage_utilization_panel",
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
				log.Println("Error getting cpu utilization: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetStorageUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
		defaultStartTime := time.Now().Add(-15 * time.Minute)
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
	rootVolumeUsage, err := GetStorageUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "disk_used_percent", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting Root Volume Utilization: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolumeUtilization"] = rootVolumeUsage

	// Get EBS1 Volume Utilization
	ebs1VolumeUsage, err := GetStorageUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "disk_used_percent", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS1 Volume Utilization: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBS1VolumeUtilization"] = ebs1VolumeUsage

	// Get EBS2 Volume Utilization
	ebs2VolumeUsage, err := GetStorageUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "disk_used_percent", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS2 Volume Utilization: ", err)
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


func GetStorageUtilizationMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "CWAgent"
	// if elementType == "EC2" {
	// 	elmType = "AWS/" + elementType
	// }
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
								Name:  aws.String("InstanceId"),
								Value: aws.String(instanceID),
							},
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
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2StorageUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

// package EC2

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

// type StorageResult struct {
// 	RootVolumeUtilization float64 `json:"RootVolumeUtilization"`
// 	EBS1VolumeUtilization float64 `json:"EBS1VolumeUtilization"`
// 	EBS2VolumeUtilization float64 `json:"EBS2VolumeUtilization"`
// }

// const (
// 	bytesToGigabytes = 1024 * 1024 * 1024
// )

// var StorageUtilizationCmd = &cobra.Command{
// 	Use:   "storage_utilization_panel",
// 	Short: "get storage utilization metrics data",
// 	Long:  `command to get storage utilization metrics data`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("running storage utilization panel")
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
// 			jsonResp, err := GetStorageUtilizationPanel(cmd, clientAuth, nil)
// 			if err != nil {
// 				log.Println("Error getting storage utilization: ", err)
// 				return
// 			}
// 			fmt.Println(jsonResp)
// 		}
// 	},
// }

// func GetStorageUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, error) {
// 	// Initialize CloudWatch client if not provided
// 	if cloudWatchClient == nil {
// 		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
// 	}
// 	elementId, _ := cmd.PersistentFlags().GetString("elementId")
// 	// elementType, _ := cmd.PersistentFlags().GetString("elementType")
// 	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
// 	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

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
// 			return "", nil
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
// 				return "", err
// 			}
// 			return "", err
// 		}
// 		startTime = &parsedStartTime
// 	}

// 	// Parse end time if provided
// 	if endTimeStr != "" {
// 		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			log.Printf("Error parsing end time: %v", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return "", err
// 			}
// 			return "", err
// 		}
// 		endTime = &parsedEndTime
// 	}

// 	// If start time is not provided, use last 15 minutes
// 	if startTime == nil {
// 		defaultStartTime := time.Now().Add(-15 * time.Minute)
// 		startTime = &defaultStartTime
// 	}

// 	// If end time is not provided, use current time
// 	if endTime == nil {
// 		defaultEndTime := time.Now()
// 		endTime = &defaultEndTime
// 	}

// 	// If start time is after end time, return null
// 	if startTime.After(*endTime) {
// 		log.Println("Start time is after end time")
// 		return "null", nil
// 	}

// 	// Get metrics for root volume utilization
// 	rootVolumeUtilization, err := GetStorageUtilizationMetricData(clientAuth, instanceId, "disk_used", startTime, endTime, cloudWatchClient)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Get metrics for EBS volume 1 utilization
// 	ebs1VolumeUtilization, err := GetStorageUtilizationMetricData(clientAuth, instanceId, "VolumeBytesUsed", startTime, endTime, cloudWatchClient)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Get metrics for EBS volume 2 utilization
// 	ebs2VolumeUtilization, err := GetStorageUtilizationMetricData(clientAuth, instanceId, "VolumeBytesUsed", startTime, endTime, cloudWatchClient)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Calculate percentages for utilization
// 	rootVolumeUtilizationPercent := (rootVolumeUtilization / bytesToGigabytes) * 100
// 	ebs1VolumeUtilizationPercent := (ebs1VolumeUtilization / bytesToGigabytes) * 100
// 	ebs2VolumeUtilizationPercent := (ebs2VolumeUtilization / bytesToGigabytes) * 100

// 	// Create JSON output
// 	jsonOutput := StorageResult{
// 		RootVolumeUtilization: rootVolumeUtilizationPercent,
// 		EBS1VolumeUtilization: ebs1VolumeUtilizationPercent,
// 		EBS2VolumeUtilization: ebs2VolumeUtilizationPercent,
// 	}

// 	jsonString, err := json.Marshal(jsonOutput)
// 	if err != nil {
// 		log.Println("Error in marshalling json in string: ", err)
// 		return "", err
// 	}

// 	return string(jsonString), nil
// }

// func GetStorageUtilizationMetricData(clientAuth *model.Auth, instanceId string, metricName string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (float64, error) {
// 	// Define input parameters for the metric query
// 	input := &cloudwatch.GetMetricDataInput{
// 		StartTime: startTime,
// 		EndTime:   endTime,
// 		MetricDataQueries: []*cloudwatch.MetricDataQuery{
// 			{
// 				Id: aws.String("m1"),
// 				MetricStat: &cloudwatch.MetricStat{
// 					Metric: &cloudwatch.Metric{
// 						MetricName: aws.String(metricName),
// 						Namespace:  aws.String("CWAgent"), // Adjust namespace if necessary
// 					},
// 					Period: aws.Int64(300),
// 					Stat:   aws.String("Average"),
// 				},
// 			},
// 		},
// 	}

// 	// Make the API call to CloudWatch to get the metric data
// 	result, err := cloudWatchClient.GetMetricData(input)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Extract the metric value
// 	if len(result.MetricDataResults) > 0 && len(result.MetricDataResults[0].Values) > 0 {
// 		return *result.MetricDataResults[0].Values[0], nil
// 	}

// 	return 0, fmt.Errorf("no metric data available for %s", metricName)
// }

// func init() {
// 	// Initialize flags or command options if needed
// }
