# awsx-getelementdetails

This tool implements the `awsx` plugin `getElementDetails`.

## Overview

The `awsx-getelementdetails` subcommand supports various cloud elements. For each element, it provides support for composite methods such as `network_utilization_panel`, `memory_utilization_panel`, `storage_utilization_panel`, and `network_utilization_panel`. The codebase is organized with a single repository containing separate folders for different element handlers.

## All Subcommands and Options

# awsx-getelementdetails

This tool implements the `awsx` plugin `getElementDetails`.

## Overview

The `awsx-getelementdetails` subcommand supports various cloud elements. For each element, it provides support for composite methods such as `network_utilization_panel`, `memory_utilization_panel`, `storage_utilization_panel`, and `network_utilization_panel`. The codebase is organized with a single repository containing separate folders for different element handlers.

## All Subcommands and Options

| S.No | Sub-command | Description |Panels Name | Specs Links |
|------|-------------|-------------|-------------|-------------|
| 1    | EC2         | This will provide all details about ec2 panels.Collect Information about specific cloud elements - Run Queries|1. cpu_utilization_panel 2. memory_utilzation_panel 3. storage_utilization_panel 4. network_utilization_panel | [EC2 Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/EC2/ec2-api-spec.md) |
| 2    | ECS         | This will provide all details about ECS panels.Collect Information about specific cloud elements - Run Queries |1. cpu_utilization_panel 2. memory_utilzation_panel 3. storage_utilization_panel 4. network_utilization_panel | [ECS Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/ECS/ecs-api-spec.md) |
| 3    | EKS         | This will provide all details about EKS panels.Collect Information about specific cloud elements - Run Queries |1. cpu_utilization_panel 2. memory_utilzation_panel 3. storage_utilization_panel 4. network_utilization_panel 5. cpu_request_panel 6. allocatable_cpu_panel 7. cpu_limits_panel | [EKS Specs](https://github.com/Appkube-awsx/awsx-getelementdetails/blob/main/specs/EKS/eks-api-spec.md)
