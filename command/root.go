package command

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/ECS"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"
	"github.com/spf13/cobra"
)

var AwsxCloudWatchMetricsCmd = &cobra.Command{
	Use:   "getAwsCloudWatchMetrics",
	Short: "getAwsCloudWatchMetrics command gets cloudwatch metrics data",
	Long:  `getAwsCloudWatchMetrics command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

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
			queryName, _ := cmd.PersistentFlags().GetString("query")
			elementType, _ := cmd.PersistentFlags().GetString("elementType")
			// cloudWatchQuery, _ := cmd.PersistentFlags().GetString("cloudWatchQuery")
			responseType, _ := cmd.PersistentFlags().GetString("responseType")

			if queryName == "cpu_utilization_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCpuUtilizationPanel(cmd, clientAuth, nil)
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
			} else if queryName == "memory_utilization_panel" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemoryUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_utilization_panel" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					// default case. it prints json
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_read" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskReadPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting CPU User Usage: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_usage_panel" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCPUUsageUserPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting CPU User Usage: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_write" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskWritePanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting Disk Write Data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_used" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskUsagePanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting in used disk data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_inpackets" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkInPacketsPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network In packets: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_inbytes" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkInBytesPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network In packets: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_Outbytes" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkOutBytesPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network In packets: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_outpackets" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkOutPacketsPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network Out packets: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "storage_utilization_panel" && elementType == "EC2" {
				jsonResp, cloudwatchMetricResp, err := EC2.GetVolumeMetricsPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting storage utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetEKScpuUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting CPU utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_requests_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPURequestData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting CPU requests : ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GeteksMemoryUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "allocatable_cpu_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetAllocatableCPUData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting allocatable cpu panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_limits_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPULimitsData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting cpu limits: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_graph_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPUUtilizationData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting cpu utilization graph panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_requests_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryRequestData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory request panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_limits_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryLimitsData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory limits panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_utilization_graph_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryUtilizationGraphData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory utilization graph panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "Network_in_out_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkInOutData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting Network_in_out_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "CPU_Utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPU_UtilizationData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting CPU_Utilization_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_usage_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryUsageData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory_usage_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_throughput_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkThroughputPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network_throughput_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "node_capacity_panel" && elementType == "EKS" {
				nodeCapacityPanel, err := EKS.GetNodeCapacityPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting node_capacity_panel: ", err)
					return
				}

				jsonResp := nodeCapacityPanel.JsonData
				cloudwatchMetricResp := nodeCapacityPanel.RawData

				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "node_uptime_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeUptimePanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting node_uptime_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_throughput_Panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkThroughputSinglePanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network_throughput_Panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "node_downtime_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeDowntimePanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting node_downtime_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_availability_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkAvailabilityData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network_availability_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "service_availability_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetServiceAvailabilityData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting service_availability_panel: ", err)
					return
				}
				if responseType == "frame" {
					for _, dataPoint := range cloudwatchMetricResp {
						fmt.Printf("%v %f\n", dataPoint.Timestamp, dataPoint.Availability)
					}
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "node_event_logs_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeEventLogsSinglePanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting node_event_logs_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "Cpu_utilization_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetContainerPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting cpu utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "Memory_utilization_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetecsMemoryUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "Network_utilization_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetNetworkUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting Network utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}

			} else if queryName == "storage_utilization_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetStorageUtilizationPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting storage utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_reservation_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetCPUReservationData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting cpu reservation data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					// default case. it prints json
					fmt.Println(jsonResp)
				}

			} else if queryName == "memory_reservation_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetMemoryReservationData(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting memory  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else {
				fmt.Println("query not found")
			}
		}
	},
}

func Execute() {
	if err := AwsxCloudWatchMetricsCmd.Execute(); err != nil {
		log.Printf("error executing command: %v\n", err)
	}
}

func init() {
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUtilizationCmd)
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("RootVolumeId", "", "root volume id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("EBSVolume1Id", "", "ebs volume 1 id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("EBSVolume2Id", "", "ebs volume 2 id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("query", "", "query")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
