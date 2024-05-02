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
)

type NodeCapacityMetrics struct {
	CPUUsage     float64 `json:"Cpu_Usage"`
	MemoryUsage  float64 `json:"Memory_Usage"`
	StorageAvail float64 `json:"Storage_Avail"`
}

type NodeCapacityPanel struct {
	RawData  map[string]*cloudwatch.GetMetricDataOutput `json:"raw_data"`
	JsonData string                                     `json:"json_data"`
}

var AwsxEKSNodeCapacityCmd = &cobra.Command{
	Use:   "node_capacity_panel",
	Short: "get node capacity metrics data",
	Long:  `command to get node capacity metrics data`,

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
			nodeCapacityPanel, err := GetNodeCapacityPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting node capacity data: ", err)
				return
			}

			jsonResp := nodeCapacityPanel.JsonData
			cloudwatchMetricResp := nodeCapacityPanel.RawData

			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNodeCapacityPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (*NodeCapacityPanel, error) {
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cpuUsageRawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_cpu_utilization", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		return nil, err
	}

	memoryUsageRawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_memory_utilization", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		return nil, err
	}

	storageAvailRawData, err := commanFunction.GetMetricData(clientAuth, instanceId, "ContainerInsights", "node_filesystem_utilization", startTime, endTime, "Sum", "ClusterName", cloudWatchClient)
	if err != nil {
		return nil, err
	}

	totalCPU := 100.0 // Assuming 100% CPU
	totalMemory := 100.0
	totalStorage := 100.0

	nodeCapacity := NodeCapacityMetrics{
		CPUUsage:     calculateComman(cpuUsageRawData, totalCPU),
		MemoryUsage:  calculateComman(memoryUsageRawData, totalMemory),
		StorageAvail: calculateComman(storageAvailRawData, totalStorage),
	}

	jsonData, err := json.Marshal(nodeCapacity)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return nil, err
	}

	nodeCapacityPanel := &NodeCapacityPanel{
		RawData:  map[string]*cloudwatch.GetMetricDataOutput{"Cpu": cpuUsageRawData, "Memory": memoryUsageRawData, "Storage": storageAvailRawData},
		JsonData: string(jsonData),
	}

	return nodeCapacityPanel, nil
}

func calculateComman(data *cloudwatch.GetMetricDataOutput, commanData float64) float64 {
	var sum float64
	for _, result := range data.MetricDataResults {
		for _, value := range result.Values {
			sum += *value
		}
	}
	return (sum / float64(len(data.MetricDataResults))) / commanData
}

func init() {
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeCapacityCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
