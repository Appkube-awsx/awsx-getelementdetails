# awsx-getelementdetails

This tool implements the `awsx` plugin `getElementDetails`.

## Overview

The `awsx-getelementdetails` subcommand supports various cloud elements. For each element, it provides support for composite methods such as `network_utilization_panel`, `memory_utilization_panel`, `storage_utilization_panel`, and `network_utilization_panel`. The codebase is organized with a single repository containing separate folders for different element handlers.

## All Subcommands and Options

### EC2

All supported subcommands and their source code locations are listed in the [AWSX API Specs](https://github.com/Appkube-awsx/awsx-api).

| S.No | Panel Name                | Description                                             | Data Output   | Commands |
|------|---------------------------|---------------------------------------------------------|---------------|----------|
| 1    | cpu_utilization_panel     | Get specific EC2 instance CPU utilization panel data    | Percentage(%) | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="cpu_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""` |
| 2    | memory_utilization_panel  | Get specific EC2 instance memory utilization panel data | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="memory_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""` |
| 3    | storage_utilization_panel | Get specific EC2 instance storage utilization panel data | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="storage_utilization_panel" --elementType="AWS/EC2" --RootVolumeId="i-05e4e6757f13da657" --EBSVolume1Id="vol-0db5984a7f9d77c4d" --EBSVolume2Id="vol-0e065bd2535df7a54" --responseType=json --startTime="" --endTime=""` |
| 4    | network_utilization_panel | Get specific EC2 instance network utilization panel data | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="network_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""` |

### EKS

All supported subcommands and their source code locations are listed in the [AWSX API Specs](https://github.com/Appkube-awsx/awsx-api).

| S.No | Panel Name                | Description                                             | Data Output   | Commands |
|------|---------------------------|---------------------------------------------------------|---------------|----------|
| 1    | cpu_utilization_panel     | Get specific EKS cluster CPU utilization panel data     | Percentage(%) | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 2    | memory_utilization_panel  | Get specific EKS cluster memory utilization panel data  | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="memory_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 3    | storage_utilization_panel | Get specific EKS cluster storage utilization panel data  | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="storage_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 4    | network_utilization_panel | Get specific EKS cluster network utilization panel data  | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="network_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 5    | cpu_requests_panel        | Collect information about specific cloud elements       | Timeseries    | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_requests_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime="--startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"` |
| 6    | allocatable_cpu_panel     | Collect information about specific cloud elements       | Timeseries    | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="allocatable_cpu_panel" --elementType="ContainerInsights" --responseType=json --startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"` |
| 7    | cpu_limits_panel          | Collect information about specific cloud elements       | Timeseries    | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_limits_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |

### ECS

All supported subcommands and their source code locations are listed in the [AWSX API Specs](https://github.com/Appkube-awsx/awsx-api).

| S.No | Panel Name                | Description                                             | Data Output   | Commands |
|------|---------------------------|---------------------------------------------------------|---------------|----------|
| 1    | cpu_utilization_panel     | Get specific ECS cluster CPU utilization panel data     | Percentage(%) | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 2    | memory_utilization_panel  | Get specific ECS cluster memory utilization panel data  | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="memory_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 3    | storage_utilization_panel | Get specific ECS cluster storage utilization panel data  | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="storage_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 4    | network_utilization_panel | Get specific ECS cluster network utilization panel data  | Bytes         | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="network_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
| 5    | cpu_requests_panel        | Collect information about specific cloud elements       | Timeseries    | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_requests_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime="--startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"` |
| 6    | cpu_limits_panel          | Collect information about specific cloud elements       | Timeseries    | `go run awsx-getelementdetails.go --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_limits_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""` |
