# awsx-getelementdetails

This tool implements the `awsx` plugin `getElementDetails`.


## All Subcommands and Options

# awsx-getelementdetails

This tool implements the `awsx` plugin `getElementDetails`.

## Overview

The `awsx-getelementdetails` subcommand supports various cloud elements. For each element, it provides support for composite methods such as `network_utilization_panel`, `memory_utilization_panel`, `storage_utilization_panel`, and `network_utilization_panel`. The codebase is organized with a single repository containing separate folders for different element handlers.

## All Subcommands and Options

| S.No | Sub-command | Description |Panels Name | Specs Links |
|------|-------------|-------------|-------------|-------------|
| 1    | EC2         | This will provide all details about ec2 panels.Collect Information about specific cloud elements - Run Queries|1. cpu_utilization_panel 2. memory_utilzation_panel 3. storage_utilization_panel 4. network_utilization_panel 5. cpu_usage_user_panel 6. cpu_usage_idle_panel 7. cpu_usgae_sys_panel 8. cpu_usage_nice_panel 9. mem_usage_total_panel 10. mem_usage_free_panel 11. mem_usage_used_panel 12. mem_physicalram_panel 13. disk_read_panel 14.disk_write_panel 15. disk_used_panel 16. disk_available_panel 17.net_inpackets 18. net_outpackets 19. net_inbytes 20. net_outbytes | [EC2 Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/EC2/ec2-api-spec.md) |
| 2    | ECS         | This will provide all details about ECS panels.Collect Information about specific cloud elements - Run Queries |1. cpu_utilization_panel 2. memory_utilzation_panel 3. storage_utilization_panel 4. network_utilization_panel 5. cpu_utilization_graph_panel 6. cpu_reservation_panel 7. cpu_sys_panel 8. cpu_nice_panel 9. memory_utilization_panel 10. memory_reservation_panel 11. container_memory_usage_panel 12. available_memory_overtime_panel 13. volume_read_bytes_panel 14. volume_write_bytes_panel 15. i/o_bytes_panel 16. disk_available 17. net_inbytes_panel 17. net_outbytes_panel 18. container_net_received_bytes_panel 19.container_net_transmitinbytes_panel 20. net_rxinbytes_panel 21. net_txinbytes_panel| [ECS Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/ECS/ecs-api-spec.md) |
| 3    | EKS         | This will provide all details about EKS panels.Collect Information about specific cloud elements - Run Queries |1. cpu_utilization_panel 2. memory_utilzation_panel 3. storage_utilization_panel 4. network_utilization_panel 5. cpu_request_panel 6. allocatable_cpu_panel 7. cpu_limits_panel 8. cpu_utilization_graph_panel 9. memeory_request_panel 10. memory_limits_panel 11.allocatable_memory_panel 12.memory_utilization_graph_panel 13. disk_utilization_panel 14. network_in_out_panel 15. cpu_utilization_pod_panel 16. memory_usage_panel 17.network_throughput 18. node_capacity 19.node_condition 20. disk_performance 21. node_events_logs 22. alerts_and_warnings_panel   | [EKS Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/EKS/eks-api-spec.md)|
| 4    | Lambda        | This will provide all details about Lambda panels.Collect Information about specific cloud elements - Run Queries |1. cost_panel 2. total_function_panel 3. idle_function_panel 4. error_rate_panel 5. throttles_fun_panel 6. total_function_cost_panel 7. top_error_products_panel 8. top_used_function_panel 9. function_panel 10. error_panel 11. throttles_panel 12.latency_panel 13. trends_panel 14. failure_function_panel 15. cpu_used_panel 16.net_receieved_panel 17.request_panel 18. memory_used_panel  19. top_failure_function-panel| [Lambda Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/Lambda/lambda-api-spec.md) |

