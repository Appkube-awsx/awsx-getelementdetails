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
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"
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
				instanceStartCount, err := EC2.GetInstanceStartCountPanel(cmd, clientAuth, nil)
				if err != nil {
					return
				}
				fmt.Println(instanceStartCount)

			} else if queryName == "instance_stop_count_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				instanceStopCount, err := EC2.GetInstanceStopCountPanel(cmd, clientAuth, nil)
				if err != nil {
					return
				}
				fmt.Println(instanceStopCount)

			} else if queryName == "instance_hours_stopped_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				instanceStopHourCount, err := EC2.GetInstanceStoppedCountPanel(cmd, clientAuth, nil)
				if err != nil {
					return
				}
				fmt.Println(instanceStopHourCount)
			} else if queryName == "instance_running_hour_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				instanceRunningHour, err := EC2.GetInstanceRunningHour(cmd, clientAuth, nil)
				if err != nil {
					return
				}
				fmt.Println(instanceRunningHour)

			} else if queryName == "instance_stop_count_panel_test" && (elementType == "EC2" || elementType == "AWS/EC2") {
				EC2.GetInstanceStartCountPanel(cmd, clientAuth, nil)

			} else if queryName == "error_rate_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				cloudwatchMetricData, err := EC2.GetInstanceErrorRatePanel(cmd, clientAuth, nil)
				fmt.Println(cloudwatchMetricData)
				if err != nil {
					return
				}
				// fmt.Println(cloudwatchMetricData)

			} else if queryName == "custom_alert_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				cloudwatchMetric, _ := EC2.GetEc2CustomAlertPanel(cmd, clientAuth, nil)
				fmt.Println(cloudwatchMetric)

				//} else if queryName == "instance_running_hour_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				//EC2.GetInstanceRunningHourPanel(cmd, clientAuth, nil)
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
			} else if queryName == "instance_status_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {

				instanceStatus, err := EC2.GetInstanceStatus(cmd, clientAuth)
				if err != nil {
					log.Fatalf("Error getting instance status: %v", err)
				}

				// Print instance information
				fmt.Printf("Instance ID: %s, Instance Type: %s, Availability Zone: %s, State: %s, System Checks Status: %s, Custom Alert: %t, Health Percentage: %.2f%%\n",
					instanceStatus.InstanceID, instanceStatus.InstanceType, instanceStatus.AvailabilityZone, instanceStatus.State, instanceStatus.SystemChecksStatus, instanceStatus.CustomAlert, instanceStatus.HealthPercentage)

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
			} else if queryName == "disk_io_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetEC2DiskIOPerformancePanel(cmd, clientAuth, nil)
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
			} else if queryName == "cpu_utilization_graph_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetCpuUtilizationGraphPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "memory_utilization_graph_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err := EC2.GetMemoryUtilizationGraphPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
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
				// } else if queryName == "instance_status_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {

				// 	instanceInfo, err := EC2.GetInstanceStatus(cmd, clientAuth)
				// 	if err != nil {
				// 		log.Fatalf("Error getting instance status: %v", err)
				// 	}
				// 	for _, info := range instanceInfo {
				// 		fmt.Printf("Instance ID: %s, Instance Type: %s, Availability Zone: %s, State: %s, System Checks Status: %s, Custom Alert: %t\n",
				// 			info.InstanceID, info.InstanceType, info.AvailabilityZone, info.State, info.SystemChecksStatus, info.CustomAlert)
				// 	}
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
			} else if queryName == "network_traffic_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				jsonResp, cloudwatchMetricResp, err, _ := EC2.GetNetworkTrafficPanel(cmd, clientAuth, nil)
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
				// } else if queryName == "latency_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				// 	jsonResp, cloudwatchMetricResp := EC2.LatencyPanel(cmd, clientAuth, nil)
				// 	if err != nil {
				// 		log.Println("Error getting latency metric data: ", err)
				// 		return
				// 	}
				// 	if responseType == "frame" {
				// 		fmt.Println(cloudwatchMetricResp)
				// 	} else {
				// 		fmt.Println(jsonResp)
				// 	}
				// 	// } else if queryName == "custom_alert_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
				// 	// 	cloudwatchMetric, _ := EC2.GetEc2CustomAlertPanel(cmd, clientAuth)
				// 	// 	fmt.Println(cloudwatchMetric)

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
				jsonResp, cloudwatchMetricResp, err := EC2.GetStorageUtilizationPanel(cmd, clientAuth, nil)
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
			} else if queryName == "node_stability_index_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeStabilityData(cmd, clientAuth, nil)
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
			} else if queryName == "incident_response_time_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetIncidentResponseTimeData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting storage utilization: ", err)
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
			} else if queryName == "allocatable_memory_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetAllocatableMemData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting allocatable memory panel: ", err)
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
			} else if queryName == "node_recovery_time_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeRecoveryTime(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting node recovery time panel: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "node_failure_panel" && elementType == "EKS" {
				jsonResp, cloudwatchMetricResp, err := EKS.GetNodeFailureData(cmd, clientAuth, nil)
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

				// } else if queryName == "data_transfer_rate_panel" && elementType == "EKS" {
				// 	jsonResp, cloudwatchMetricResp, err := EKS.GetEksDataTransferRatePanel(cmd, clientAuth, nil)
				// 	if err != nil {
				// 		log.Println("Error getting data_transfer_rate_panel: ", err)
				// 		return
				// 	}
				// 	if responseType == "frame" {
				// 		fmt.Println(cloudwatchMetricResp)
				// 	} else {
				// 		fmt.Println(jsonResp)
				// 	}
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
			} else if queryName == "available_memory_over_time_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetAvailableMemoryOverTimeData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting volume write bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "top_events_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, err := ECS.GetECSTopEventsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting top events data: ", err)
					return
				}
				fmt.Println(jsonResp)

			} else if queryName == "registration_events_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, err := ECS.GetRegistrationEventsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting registration events data: ", err)
					return
				}
				fmt.Println(jsonResp)

			} else if queryName == "deregistration_events_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, err := ECS.GetDeRegistrationEventsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting deregistration events data: ", err)
					return
				}
				fmt.Println(jsonResp)
			} else if queryName == "resource_deleted_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, err := ECS.GetECSResourceDeletedEvents(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting resource deleted panel: ", err)
					return
				}
				fmt.Println(jsonResp)
			} else if queryName == "resources_created_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, err := ECS.GetECSResourceCreatedEvents(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting resource created panel: ", err)
					return
				}
				fmt.Println(jsonResp)
			} else if queryName == "failed_tasks_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				failedTask, err := ECS.GetECSFailedTasksEvents(cmd, clientAuth, nil)
				if err != nil {

					return
				}
				fmt.Println(failedTask)
			} else if queryName == "failed_services_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				failedService, err := ECS.GetECSFailedServiceEvents(cmd, clientAuth, nil)
				if err != nil {

					return
				}
				fmt.Println(failedService)
			} else if queryName == "active_services_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				activeService, err := ECS.GetECSActiveServiceEvents(cmd, clientAuth, nil)
				if err != nil {

					return
				}
				fmt.Println(activeService)
			} else if queryName == "active_connection_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				activeConnection, err := ECS.GetECSActiveConnectionEvents(cmd, clientAuth, nil)
				if err != nil {
					//log.Println("Error getting resouce created panel: ", err)
					return
				}
				fmt.Println(activeConnection)
			} else if queryName == "new_connection_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				newConnection, err := ECS.GetECSNewConnectionEvents(cmd, clientAuth, nil)
				if err != nil {

					return
				}
				fmt.Println(newConnection)

			} else if queryName == "active_tasks_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				activeTask, err := ECS.GetECSActiveTaskEvents(cmd, clientAuth, nil)
				if err != nil {

					return
				}
				fmt.Println(activeTask)
			} else if queryName == "resource_updated_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, err := ECS.GetECSResourceUpdatedEvents(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting resource updated panel: ", err)
					return
				}
				fmt.Println(jsonResp)
			} else if queryName == "container_net_received_inbytes_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSContainerNetRxInBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting container net received in bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "container_net_transmit_inbytes_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetECSContainerNetTxInBytesPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting container net transmitted in bytes data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}

			} else if queryName == "container_memory_usage_panel" && (elementType == "AWS/ECS" || elementType == "ECS") {
				jsonResp, cloudwatchMetricResp, err := ECS.GetContainerMemoryUsageData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory usage for ECS: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "service_error_panel" && (elementType == "ECS" || elementType == "AWS/ECS") {
				events, err := ECS.ListServiceErrors()
				if err != nil {
					return
				}
				for _, event := range events {
					// Perform further processing on each event
					fmt.Println("TIME STAMP:", event.Timestamp)
					fmt.Println("SERVICE NAME:", event.ServiceName)
					fmt.Println("TASK ID:", event.TaskID)
					fmt.Println("ERROR TYPE:", event.ErrorType)
					fmt.Println("ERROR DESCRIPTION:", event.ErrorDescription)
					fmt.Println("RESOLUTION TIME:", event.ResolutionTime)
					fmt.Println("IMPACT LEVEL:", event.ImpactLevel)
					fmt.Println("RESOLUTION DETAILS:", event.ResolutionDetails)
					fmt.Println("STATUS:", event.Status)
					fmt.Println("---------------------------------------")
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
			} else if queryName == "error_and_warning_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaErrorAndWarningData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda error and warning  data: ", err)
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
			} else if queryName == "functions_by_region_panel" && elementType == "Lambda" {
				log.Printf("ClientAuth: %+v\n", clientAuth)
				jsonResp, cloudwatchMetricResp, totalFunctions, _ := Lambda.GetLambdaFunctionsByRegion(clientAuth)
				if err != nil {
					log.Println("Error total functions panel data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
					fmt.Println("Function By Region:", totalFunctions)
				} else {
					fmt.Println(jsonResp)
					fmt.Println("Function By Region:", totalFunctions)
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
			} else if queryName == "concurrency_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaConcurrencyData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda concurrency data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "used_and_unused_memory_data_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaUnusedMemoryPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda used and unused memory data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "max_memory_used_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaMaxMemoryData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda max memory used data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "max_memory_used_graph_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaMaxMemoryGraphData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda max memory used graph data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "number_of_calls_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaNumberOfCallsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda max memory used data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cold_start_duration_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaColdStartData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda max memory used graph data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "execution_time_panel" && elementType == "Lambda" {
				jsonResp, cloudwatchMetricResp, err := Lambda.GetLambdaExecutionTimePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda max memory used graph data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "invocation_trend_panel" && elementType == "Lambda" {
				jsonResp, err := Lambda.GetInvocationTrendData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda max memory used graph data: ", err)
					return
				}
				fmt.Println(jsonResp)

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
				jsonResp, err := Lambda.GetErrorMessageCountData(cmd, clientAuth, nil)
				if err != nil {
					return
				}
				fmt.Println(jsonResp)
			} else if queryName == "throttling_trends_panel" && elementType == "Lambda" {
				jsonResp, err := Lambda.GetThrottlingTrendsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting lambda functions  data: ", err)
					return
				}
				fmt.Println(jsonResp)
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
			} else if queryName == "success_and_failed_function_panel" && elementType == "Lambda" {
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
			} else if queryName == "successful_and_failed_events_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				responseType, _ := cmd.PersistentFlags().GetString("responseType")
				jsonResp, uptimeMetricResp, err := ApiGateway.GetApiSuccessFailedData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting API uptime data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(uptimeMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "top_events_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, err := ApiGateway.GetTopEventsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting total api data: ", err)
					return
				}
				fmt.Println(jsonResp)

			} else if queryName == "successful_event_details_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, err := ApiGateway.GetSuccessEventData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting successful api data: ", err)
					return
				}
				fmt.Println(jsonResp)

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
			} else if queryName == "total_api_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetTotalApiData(clientAuth, nil)
				if err != nil {
					log.Println("Error getting total api data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "failed_event_details" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, err := ApiGateway.GetFailedEventData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting total api data: ", err)
					return
				}
				fmt.Println(jsonResp)
			} else if queryName == "error_logs_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, err := ApiGateway.GetErrorLogsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting total api data: ", err)
					return
				}
				fmt.Println(jsonResp)

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
			} else if queryName == "response_time_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiResponseTimePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting api response time data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "uptime_percentage_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				responseType, _ := cmd.PersistentFlags().GetString("responseType")
				jsonResp, uptimeMetricResp, err := ApiGateway.GetApiUptimeData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting API uptime data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(uptimeMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cache_hit_count_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				responseType, _ := cmd.PersistentFlags().GetString("responseType")
				jsonResp, uptimeMetricResp, err := ApiGateway.GetApiCacheHitsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting cache hit data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(uptimeMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cache_miss_count_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				responseType, _ := cmd.PersistentFlags().GetString("responseType")
				jsonResp, uptimeMetricResp, err := ApiGateway.GetApiCacheMissData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting cache miss data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(uptimeMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "downtime_incident_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp := ApiGateway.GetDowntimeIncidentsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting total api data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "uptime_of_deployment_stages" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp := ApiGateway.GetApiUptimedata(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting uptime deployment data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "total_api_calls_panel" && (elementType == "AWS/ApiGateway" || elementType == "ApiGateway") {
				jsonResp, cloudwatchMetricResp, err := ApiGateway.GetApiCallsData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting total api request data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_utilization_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSCpuUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting cpu utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "database_connections_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetDatabaseConnectionsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting database connections: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "index_size_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetIndexSizePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting index size: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "maintenance_schedule_overview_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				events, err := RDS.ListScheduleOverview()
				if err != nil {
					return
				}
				for _, event := range events {
					// Perform further processing on each event
					fmt.Println("MAINTENANCE TYPE:", event.MaintenanceType)
					fmt.Println("DESCRIPTION:", event.Description)
					fmt.Println("START TIME:", event.StartTime)
					fmt.Println("END TIME:", event.EndTime)
					fmt.Println("---------------------------------------")
				}

			} else if queryName == "cpu_credit_usage_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetCPUCreditUsagePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting credit usage: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "storage_utilization_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSStorageUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting credit usage: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_credit_balance_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetCPUCreditBalancePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting credit balance: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_surplus_credit_balance_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetCPUSurplusCreditBalance(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting surplus credit balance: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "cpu_surplus_credits_charged_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetCPUSurplusCreditCharged(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting surplus credits charged: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "write_iops_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSWriteIOPSPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting write iops: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "read_iops_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSReadIOPSPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting read iops: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_utilization_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSNetworkUtilizationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting memory utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_traffic_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSNetworkTrafficPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network utilization: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "instance_health_check_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				instanceInfo, err := RDS.GetDBInstanceHealthCheck()
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
			} else if queryName == "cpu_utilization_graph_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSCPUUtilizationGraphPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network utilization graph: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "alert_and_notification_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, err := RDS.GetAlertsAndNotificationsPanell(cmd, clientAuth)
				if err != nil {
					log.Println("Error getting network inbound metric data: ", err)
					return
				}
				// if responseType == "frame" {
				// 	fmt.Println(cloudwatchMetricResp)
				// } else {
				fmt.Println(jsonResp)
			} else if queryName == "iops_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSIopsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting iops data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "freeable_memory_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSFreeableMemoryPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting freeable memory data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "free_storage_space_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSFreeStorageSpacePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting free storage space data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "disk_queue_depth_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSDiskQueueDepthPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting disk queue depth data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "replication_slot_disk_usage" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSReplicationSlotDiskUsagePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting read iopslogvolume data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_receive_throughput_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSNetworkReceiveThroughputPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network receive throughput data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "network_transmit_throughput_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSNetworkTransmitThroughputPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting network receive throughput data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "database_workload_overview_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSDBLoadPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting database workload overview data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "db_load_non_cpu_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSDBLoadNonCPU(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting database non-cpu load data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "db_load_cpu_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err, _ := RDS.GetRDSDBLoadCPU(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting database cpu load data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "latency_analysis_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSLatencyAnalysisData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting latency analysis data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "transaction_logs_generation_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetTransactionLogsGenerationPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting read iops: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "transaction_logs_disk_usage_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetTransactionLogsDiskUsagePanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting read iops: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}

			} else if queryName == "recent_error_log_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRdsErrorLogsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting recent error logs: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "recent_event_log_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp,err := RDS.GetRecentEventLogsPanel(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting recent events logs: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "uptime_percentage" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, cloudwatchMetricResp, err := RDS.GetRDSUptimeData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting uptime percentage data: ", err)
					return
				}
				if responseType == "frame" {
					fmt.Println(cloudwatchMetricResp)
				} else {
					fmt.Println(jsonResp)
				}
			} else if queryName == "error_analysis_panel" && (elementType == "RDS" || elementType == "AWS/RDS") {
				jsonResp, _ := RDS.GetErrorAnalysisData(cmd, clientAuth, nil)
				if err != nil {
					log.Println("Error getting read iops: ", err)
					return
				}
				for _, entry := range jsonResp {
					fmt.Printf("%+v\n", entry)
				}
				// if responseType == "json" {
				// 	fmt.Println(cloudwatchMetricResp)
				// } else {
				// 	fmt.Println(jsonResp)
				// }
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
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CustomAlertPanelCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemoryUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkUtilizationCmd)
	// AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2StorageUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUsageUserCmd)
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
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2CpuUtilizationGraphsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2MemoryUtilizationGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.ListErrorsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEC2NetworkTrafficCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2DiskAvailableCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkInboundCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkOutboundCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2AlarmandNotificationcmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceStopCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEC2DiskIOPerformanceCmd)
	//AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceStopCmdTest)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2NetworkOutBytesCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceStatusCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2ErrorRatePanelCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EC2.AwsxEc2InstanceHealthCheckCmd)
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
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSNodeStabilityCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(EKS.AwsxEKSIncidentResponseTimeCmd)
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
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxResourceCreatedPanelCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxResourceDeletedPanelCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxResourceUpdatedPanelCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ECS.AwsxEcsServiceErrorCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaCpuCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaFailureCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaSuccessFailureCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(Lambda.AwsxLambdaNumberOfCallsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSCpuUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSNetworkTrafficCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSNetworkUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSFreeStorageSpaceCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSFreeableMemoryCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSDiskQueueDepthCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSNetworkTrafficCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSReplicationSlotDiskUsageCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSNetworkReceiveThroughputCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSNetworkTransmitThroughputCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSDBLoadCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSDBLoadNonCPUCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSDBLoadCPUCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSStorageUtilizationCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSLatencyAnalysisCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSCPUCreditUsageCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSCPUSurplusCreditBalanceCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSSurplusCreditsChargedCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSCpuUtilizationGraphCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSDatabaseConnectionsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSDBLoadCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSIndexSizeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxDBInstanceHealthCheckCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSIopsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSReadIOPSCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSTransactionLogsDiskCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSTransactionLogsGenCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSWriteIOPSCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSErrorAnalysisCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(RDS.AwsxRDSUptimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.ApiResponseTimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.AwsxApiCacheHitsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.AwsxApiCacheMissCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.AwsxApiDowntimeIncidentsCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.AwsxApiUptimeCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.AwsxApiDeploymentCmd)
	AwsxCloudWatchMetricsCmd.AddCommand(ApiGateway.AwsxApiCallsCmd)

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
