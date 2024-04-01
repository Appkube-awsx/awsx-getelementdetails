package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

type StorageUtlizationResult struct {
	RootVolumeUtilization float64 `json:"RootVolumeUsage"`
	EBS1VolumeUtilization float64 `json:"EBSVolume1Usage"`
	EBS2VolumeUtilization float64 `json:"EBSVolume2Usage"`
}

var AwsxRDSStorageUtilizationCmd = &cobra.Command{
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
			jsonResp, cloudwatchMetricResp, err := GetRDSStorageUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting storage utilization: ", err)
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

func GetRDSStorageUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rootVolumeUsage, err := GetRDSStorageUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "FreeStorageSpace", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting Root Volume Utilization: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolumeUtilization"] = rootVolumeUsage

	// Get EBS1 Volume Utilization
	ebs1VolumeUsage, err := GetRDSStorageUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "FreeStorageSpace", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS1 Volume Utilization: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBS1VolumeUtilization"] = ebs1VolumeUsage

	// Get EBS2 Volume Utilization
	ebs2VolumeUsage, err := GetRDSStorageUtilizationMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Average", "FreeStorageSpace", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS2 Volume Utilization: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBS2VolumeUtilization"] = ebs2VolumeUsage

	// Calculate average of all three volumes
	averageStorageResult := StorageUtlizationResult{
		RootVolumeUtilization: round(calculateAverage(rootVolumeUsage)/1000000000, 2),
		EBS1VolumeUtilization: round(calculateAverage(ebs1VolumeUsage)/2000000000, 2),
		EBS2VolumeUtilization: round(calculateAverage(ebs2VolumeUsage)/2000000000, 2),
	}

	jsonString, err := json.Marshal(averageStorageResult)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetRDSStorageUtilizationMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "AWS/RDS"

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{},
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

func round(val float64, places int) float64 {
	precision := math.Pow(10, float64(places))
	return math.Round(val*precision) / precision
}

func init() {
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxRDSStorageUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
