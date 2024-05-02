package EKS

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
	"strconv"
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
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Root Volume Usage
	rootVolumeUsage, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_utilization", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting root volume usage: ", err)
		return "", nil, err
	}
	rootVolumeUsageValue := *rootVolumeUsage.MetricDataResults[0].Values[0]
	rootVolumeUsageStr := strconv.FormatFloat(rootVolumeUsageValue, 'f', 2, 64)

	// Get EBS Volume 1 Usage
	ebsVolume1Usage, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_inodes", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS volume 1 usage: ", err)
		return "", nil, err
	}
	ebsVolume1Percentage := (*ebsVolume1Usage.MetricDataResults[0].Values[0] / 10000000.0) // Replace 100.0 with the total space for EBS Volume 1
	ebsVolume1PercentageStr := strconv.FormatFloat(ebsVolume1Percentage, 'f', 2, 64)

	// Get EBS Volume 2 Usage
	ebsVolume2Usage, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_inodes", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS volume 2 usage: ", err)
		return "", nil, err
	}
	ebsVolume2Percentage := (*ebsVolume2Usage.MetricDataResults[0].Values[0] / 10999999.0) // Replace 200.0 with the total space for EBS Volume 2
	ebsVolume2PercentageStr := strconv.FormatFloat(ebsVolume2Percentage, 'f', 2, 64)

	// Convert formatted strings back to float64
	rootVolumeUsageFloat, err := strconv.ParseFloat(rootVolumeUsageStr, 64)
	if err != nil {
		log.Println("Error converting string to float64: ", err)
		return "", nil, err
	}
	ebsVolume1PercentageFloat, err := strconv.ParseFloat(ebsVolume1PercentageStr, 64)
	if err != nil {
		log.Println("Error converting string to float64: ", err)
		return "", nil, err
	}
	ebsVolume2PercentageFloat, err := strconv.ParseFloat(ebsVolume2PercentageStr, 64)
	if err != nil {
		log.Println("Error converting string to float64: ", err)
		return "", nil, err
	}

	// Create JSON output
	jsonOutput := StorageUtilizationResult{
		RootVolumeUsage: rootVolumeUsageFloat,
		EBSVolume1Usage: ebsVolume1PercentageFloat,
		EBSVolume2Usage: ebsVolume2PercentageFloat,
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}
	return string(jsonString), cloudwatchMetricData, nil
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
