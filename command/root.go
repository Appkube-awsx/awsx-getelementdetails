package command

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/ECS"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/Lambda"
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
			} else if queryName == "memory_utilization_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemoryUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_utilization_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkUtilizationPanel(cmd, clientAuth, nil)
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
			} else if queryName == "cpu_usage_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCPUUsageUserPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting CPU Usage User: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_usage_sys_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCPUUsageSysPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting CPU Usage Sys metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_usage_nice" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCPUUsageNicePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting CPU Usage Nice Metric Data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_usage_idle" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCPUUsageIdlePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting CPU Usage Idle Metric Data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "mem_usage_free" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemUsageFreePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory usage free metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "mem_cached" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemCachePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory cached metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			}else if queryName == "mem_usage_total" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemUsageTotal(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory usage total metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "mem_usage_used" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemUsageUsed(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory used usage metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_write" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskWritePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting Disk Write Metric Data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_read" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskReadPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting in Disk read metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_availble" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskAvailablePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting in Disk available metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_used" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetDiskUsedPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting in Disk used metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_inpackets" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkInPacketsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network Input packets metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_inbytes" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkInBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network In Bytes metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_Outbytes" && (elementType == "EC2" || elementType == "AWS/EC2")  {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkOutBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network Out bytes metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_outpackets" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkOutPacketsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network Out packets: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "storage_utilization_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkOutPacketsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting storage utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_panel" && (elementType == "AWS/EKS" || elementType == "EKS") {
				jsonResp, cloudwatchMetricResp, err := EKS.GetEKScpuUtilizationPanel(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPURequestData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GeteksMemoryUtilizationPanel(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "storag_utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetStorageUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetDiskUtilizationData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "allocatable_cpu_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetAllocatableCPUData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPULimitsData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPUUtilizationData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryRequestData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryLimitsData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryUtilizationGraphData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization graph panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_in_out_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkInOutData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting Network_in_out_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_node_graph_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetCPUUtilizationNodeData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting cpu utilization node graph panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_usage_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetMemoryUsageData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkThroughputPanel(cmd, clientAuth, nil)
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
				nodeCapacityPanel, err := EKS.GetNodeCapacityPanel(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeUptimePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting node_uptime_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_throughput_single_panel" && elementType == "EKS" {
				cloudwatchMetricResp, jsonResp, err := EKS.GetNetworkThroughputSinglePanel(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeDowntimePanel(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkAvailabilityData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetServiceAvailabilityData(cmd, clientAuth, nil)
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
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeEventLogsSinglePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting node_event_logs_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECScpuUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting cpu utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_utilization_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSMemoryUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_graph_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetCPUUtilizationGraphData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting cpu utilization graph for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_utilization_graph_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetMemoryUtilizationGraphData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization graph for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
				//} else if queryName == "Network_utilization_panel" && elementType == "AWS/ECS" {
				//jsonResp, cloudwatchMetricResp, err := ECS.GetNetworkUtilizationPanel(cmd, clientAuth,nil)
				//if err != nil {
				//log.Println("Error getting Network utilization for ECS: ", err)
				//return
				//}
				//if responseType == "frame" {
				//fmt.Println(cloudwatchMetricResp)
				//} else {
				//fmt.Println(jsonResp)
				//}

			} else if queryName == "storage_utilization_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetStorageUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting storage utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_reservation_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetCPUReservationData(cmd, clientAuth, nil)
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

			} else if queryName == "memory_reservation_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetMemoryReservationData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "error_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaErrorData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda error  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "throttles_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaThrottleData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda throttles  data: ", err)
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
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemoryUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkUtilizationCmd)
	// AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2StorageUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUsageUserCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUsageIdleCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuSysTimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUsageNiceCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemCachedCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemoryUsageTotalCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemoryUsageUsedCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemoryUsageFreeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkInBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkOutBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkOutPacketsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkInPacketsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2DiskReadCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2DiskWriteCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2DiskUsedCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2DiskAvailableCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSAllocatableCpuCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSCpuLimitsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSCpuRequestsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSCpuUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSCpuUtilizationGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSCpuUtilizationNodeGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSDiskUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSMemoryLimitsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSMemoryRequestsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSMemoryUsageCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSMemoryUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSMemoryUtilizationGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNetworkAvailabilityCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNetworkInOutCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNetworkThroughputCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNetworkThroughputSingleCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNetworkUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeCapacityCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeDowntimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeEventLogsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeUptimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSServiceAvailabilityCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSStorageUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSCpuUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSMemoryUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSStorageUtilizationCmd)
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
