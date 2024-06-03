package EKS

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

type StorageUtilizationResult struct {
	RootVolumeUtilization float64 `json:"RootVolumeUsage"`
	EBS1VolumeUtilization float64 `json:"EBSVolume1Usage"`
	EBS2VolumeUtilization float64 `json:"EBSVolume2Usage"`
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

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}
	instanceId, err = comman_function.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Root Volume Usage
	rootVolumeUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_utilization", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting root volume usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["RootVolumeUtilization"] = rootVolumeUsage

	// Get EBS Volume 1 Usage
	ebs1VolumeUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_inodes", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS1 volume  usage: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["EBS1VolumeUtilization"] = ebs1VolumeUsage

	// Get EBS Volume 2 Usage
	ebs2VolumeUsage, err := comman_function.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_inodes", startTime, endTime, "Average", "ClusterName", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting EBS2 volume  usage: ", err)
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
	jsonOutput := StorageUtilizationResult{
		RootVolumeUtilization: rootVolumeAvgFloat,
		EBS1VolumeUtilization: ebs1VolumeAvgFloat,
		EBS2VolumeUtilization: ebs2VolumeAvgFloat,
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}
	return string(jsonString), cloudwatchMetricData, nil
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
	comman_function.InitAwsCmdFlags(AwsxEKSStorageUtilizationCmd)
	
}
