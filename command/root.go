package command

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/ApiGateway"
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
			} else if queryName == "instance_start_count_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				EC2.GetInstanceStartCountPanel(cmd, clientAuth, nil)

			} else if queryName == "instance_stop_count_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				instanceStopCount, err := EC2.GetInstanceStopCountPanel(cmd, clientAuth, nil)
				if err != nil {
					return
				}
				fmt.Println(instanceStopCount)

			} else if queryName == "instance_stop_count_panel_test" && (elementType == "EC2" || elementType == "AWS/EC2") {
				EC2.GetInstanceStartCountPanel(cmd, clientAuth, nil)

			} else if queryName == "instance_error_rate_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				EC2.GetInstanceErrorRatePanel(cmd, clientAuth, nil)

			} else if queryName == "custom_alert_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				EC2.GetEc2CustomAlertPanel(cmd, clientAuth)

			} else if queryName == "instance_running_hour_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				EC2.GetInstanceRunningHourPanel(cmd, clientAuth, nil)
			} else if queryName == "hosted_services_overview_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				hostedServicesOverview, err := EC2.GetHostedServicesData(cmd)
				if err != nil {
					log.Fatal(err)
				}

				// Print header
				fmt.Printf("%-15s %-15s %-15s %-10s %-15s %-15s\n",
					"Service Name", "Health Status", "Response Time", "Error Rate", "Availability", "Throughput")

				// Print service overview
				for _, service := range hostedServicesOverview {
					fmt.Printf("%-15s %-15s %-15s %-10s %-15s %-15s\n",
						service.ServiceName, service.HealthStatus, service.ResponseTime, service.ErrorRate,
						service.Availability, service.Throughput)
				}
			} else if queryName == "error_tracking_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				events, err := EC2.ListErrorEvents()
				if err != nil {
					return
				}
				for _, event := range events {
					// Perform further processing on each event
					fmt.Println("Event ID:", event.EventID)
					fmt.Println("Timestamp:", event.Timestamp)
					fmt.Println("Error Code:", event.ErrorCode)
					fmt.Println("Severity:", event.Severity)
					fmt.Println("Description:", event.Description)
					fmt.Println("Source Component:", event.SourceComponent)
					fmt.Println("Action Taken:", event.ActionTaken)
					fmt.Println("Resolution Status:", event.ResolutionStatus)
					fmt.Println("Additional Notes:", event.AdditionalNotes)
					fmt.Println("---------------------------------------")
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
			} else if queryName == "cpu_usage_user_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "cpu_usage_nice_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "cpu_usage_idle_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "mem_usage_free_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "mem_cached_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "mem_usage_total_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "mem_usage_used_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "disk_writes_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "disk_reads_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "disk_available_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "disk_used_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "net_inpackets_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "net_inbytes_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "net_outbytes_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "net_outpackets_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
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
			} else if queryName == "net_throughput_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkThroughputPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network throught metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "instance_status_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {

				instanceInfo, err := EC2.GetInstanceStatus(cmd, clientAuth)
				if err != nil {
					log.Fatalf("Error getting instance status: %v", err)
				}
				for _, info := range instanceInfo {
					fmt.Printf("Instance ID: %s, Instance Type: %s, Availability Zone: %s, State: %s, System Checks Status: %s, Custom Alert: %t\n",
						info.InstanceID, info.InstanceType, info.AvailabilityZone, info.State, info.SystemChecksStatus, info.CustomAlert)
				}
			} else if queryName == "instance_health_check_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				instanceInfo, err := EC2.GetInstanceHealthCheck()
				if err != nil {
					return
				}
				fmt.Printf("%-20s %-15s %-15s %-15s %-20s %-15s %-5s %-25s %-25s\n",
					"Instance ID", "Instance Type", "Availability Zone", "State", "System Checks Status",
					"Instance Checks Status", "Alarm", "System Check Time", "Instance Check Time")

				// Print instance information
				for _, info := range instanceInfo {
					fmt.Printf("%-20s %-15s %-15s %-15s %-20s %-15s %-5t %-25s %-25s\n",
						info.InstanceID, info.InstanceType, info.AvailabilityZone, info.InstanceStatus,
						info.SystemChecks, info.InstanceChecks, info.SystemCheck, info.InstanceCheck)
				}
			} else if queryName == "network_inbound_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkInBoundPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network inbound metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_outbound_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetNetworkOutBoundPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network inbound metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "latency_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetLatencyPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting latency metric data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "alert_and_notification_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, err := EC2.GetAlertsAndNotificationsPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network inbound metric data: ", err)
					return
				}
				// if responseType == "frame" {
				// 	fmt.Println(cloudwatchMetricResp)
				// } else {
				fmt.Println(jsonResp)
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
			} else if queryName == "storage_utilization_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetStorageUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting storage utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
				// } else if queryName == "incident_response_time_panel" && elementType == "EKS" {
				// 	jsonResp, cloudwatchMetricResp, err := EKS.GetIncidentResponseTimeData(cmd, clientAuth, nil)
				// 	if err != nil {
				// 		log.Println("Error getting storage utilization: ", err)
				// 		return
				// 	}
				// 	if responseType == "frame" {
				// 		fmt.Println(cloudwatchMetricResp)
				// 	} else {
				// 		fmt.Println(jsonResp)
				// 	}
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
			} else if queryName == "cpu_graph_utilization_panel" && elementType == "EKS" {
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
			} else if queryName == "memory_graph_utilization_panel" && elementType == "EKS" {
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
			} else if queryName == "disk_io_performance_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNetworkInOutData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting disk_io_performance panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_node_utilization_panel" && elementType == "EKS" {
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
					log.Println("Error getting network_throughput_single_Panel: ", err)
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
			} else if queryName == "node_condition_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeConditionPanel(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting node_condition panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}

			} else if queryName == "data_transfer_rate_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetEksDataTransferRatePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting data_transfer_rate_panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "resource_utilization_patterns_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetResourceUtilizationData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting resource_utilization_panel: ", err)
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
			} else if queryName == "cpu_graph_utilization_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
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
			} else if queryName == "memory_graph_utilization_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
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
			} else if queryName == "Network_utilization_panel" && elementType == "AWS/ECS" {
				jsonResp, cloudwatchMetricResp, err := ECS.GetNetworkUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting Network utilization for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}

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

			} else if queryName == "net_rxinbytes_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSNetworkRxInBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network received in bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_txinbytes_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSNetworkTxInBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network transmitted in bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "volume_read_bytes_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSReadBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting volume read bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "volume_write_bytes_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSWriteBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting volume write bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
				// } else if queryName == "iam_role_and_policies_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				// 	jsonResp, cloudwatchMetricResp, err := ECS.GetECSIAMRolesPanel(cmd, clientAuth)
				// 	if err != nil {
				// 		log.Println("Error getting iam role and policies for ECS: ", err)
				// 		return
				// 	}
				// 	if responseType == "frame" {
				// 		fmt.Println(cloudwatchMetricResp)
				// 	} else {
				// 		fmt.Println(jsonResp)
				// 	}
				// } else if queryName == "active_services_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				// 	jsonResp, rawResp, err := ECS.GetActiveServicesPanelData(cmd, clientAuth)
				// 	if err != nil {
				// 		log.Println("Error getting active services metrics data: ", err)
				// 		return
				// 	}
				// 	if responseType == "frame" {
				// 		fmt.Println("Raw Data:")
				// 		fmt.Println(rawResp)
				// 	} else {
				// 		fmt.Println("JSON Data:")
				// 		fmt.Println(string(jsonResp))
				// 	}
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
			} else if queryName == "latency_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaLatencyData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting latency panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_used_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaMemoryData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory used panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "total_functions_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaTotalFunctionData(clientAuth, nil)
				if err != nil {
					log.Println("Error total functions panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "idle_functions_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp := Lambda.GetLambdaIdleFunctionData(clientAuth, nil)
				if err != nil {
					log.Println("Error getting idle functions panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "throttles_function_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp := Lambda.GetLambdaThrottlesFunctionData(clientAuth)
				if err != nil {
					log.Println("Error getting throttles function panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "trends_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaTrendsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting trends panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "net_received_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaNetReceivedData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda net received  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "request_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaRequestData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda request  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "failure_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaFailureData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda failure  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "error_messages_count_panel" && elementType == "Lambda" {
				Lambda.GetLambdaErrorMessageCountPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda functions  data: ", err)
					return
				}
			} else if queryName == "throttling_trends_panel" && elementType == "Lambda" {
				Lambda.GetThrottlingTrendsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda functions  data: ", err)
					return
				}
			} else if queryName == "function_panel" && elementType == "Lambda" {
				Lambda.GetFunctionPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda functions  data: ", err)
					return
				}
			} else if queryName == "top_failure_function_panel" && elementType == "Lambda" {
				Lambda.GetTotalFailureFunctionsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda functions  data: ", err)
					return
				}
				//} else if queryName == "top_used_function_panel" && elementType == "Lambda" {
				//Lambda.GetTopUsedFunctionsPanel(cmd, clientAuth, nil)
				//if err != nil {
				//log.Println("Error getting lambda functions  data: ", err)
				//return
				//}
			} else if queryName == "success_and_failure_function_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaSuccessFailureData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting success and failure function panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_used_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaCpuData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda cpu used  data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "rest_api_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiGatewayRestAPIData(clientAuth, nil)
				if err != nil {
					log.Println("Error getting rest api data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "http_api_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiGatewayHttpApiData(clientAuth, nil)
				if err != nil {
					log.Println("Error getting http api data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "websocket_api_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiGatewayWebSocketAPIData(clientAuth, nil)
				if err != nil {
					log.Println("Error getting http api data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "4xx_errors_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApi4xxErrorData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda error  data: ", err)
					log.Println("Error getting api 4xx errors data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "5xx_errors_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApi5xxErrorData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting api 4xx errors data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "latency_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiLatencyData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting api latency data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "integration_latency_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiIntegrationLatencyData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting api integration latency data: ", err)
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
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2hostedServicesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.ListErrorsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2DiskAvailableCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkInboundCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkOutboundCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2AlarmandNotificationcmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceStopCmd)
	//AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceStopCmdTest)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkOutBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceStatusCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceErrorRateCmd)
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
	// AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSIncidentResponseTimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeDowntimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeEventLogsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeUptimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSServiceAvailabilityCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSStorageUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSCpuUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSCpuUtilizationGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxCpuReservedCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSMemoryUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSMemoryUtilizationGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxMemoryReservedCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSStorageUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSNetworkRxInBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSNetworkTxInBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSReadBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxECSWriteBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaCpuCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaFailureCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaSuccessFailureCmd)

	AwsxCloudWatchMetricsCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
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
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
