
- [awsx-getelementdetails](#awsx-getelementdetails)
- [all subcommands and options for EC2 getElementDetails](#all-subcommands-and-options-for-EC2)
- [all subcommands and options for EKS getElementDetails](#all-subcommands-and-options-for-EKS)
- [all subcommands and options for ECS getElementDetails](#all-subcommands-and-options-for-ECS)




# awsx-getelementdetails
It implements the awsx plugin getElementDetails 

This subcommand will need to take care for all the cloud elements and for every element, we need to support the composite method like network_utilization_panel,memory_utilization_panel,storage_utilization_panel and network_utilization_panel. So , we can keep a single repo for the subcommand and keep separate folders for the different element handlers.

# all subcommands and options for EC2
  <!-- getElementDetails references --> 
All the supported subcommands and there source code locations are mentiioned in 

    https://github.com/AppkubeCloud/appkube-api-specs/blob/main/awsx-api.md

| S.No | Panel Name          | Description                                           | Data output                                 | Commands                                                                                                                                                                            | 
|------|-----------------------|-------------------------------------------------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1    | cpu_utilization_panel | This will get the specific EC2 instance cpu utilization panel data in hybrid structure  | Percentage(%)                                    | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="cpu_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""
 |    |
| 2    | memory_utilization_panel  | This will get the specific EC2 instance memory utilization panel data in hybrid structure | Bytes                                   | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="memory_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""
 |
| 3    | storage_utilization_panel  | This will get the specific EC2 instance storage utilization panel data in hybrid structure | Bytes                                      |go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657" --query="storage_utilization_panel" --elementType="AWS/EC2" --RootVolumeId="i-05e4e6757f13da657" --EBSVolume1Id="vol-0db5984a7f9d77c4d" --EBSVolume2Id="vol-0e065bd2535df7a54"  --responseType=json --startTime="" --endTime=""
 |
| 4    | network_utilization_panel   |This will get the specific EC2 instance network utilization panel data in hybrid structure | Bytes                                 | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --instanceID="i-05e4e6757f13da657"  --query="network_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""
 |
 



# all subcommands and options for EKS
<!-- getElementDetails references --> 

All the supported subcommands and there source code locations are mentiioned in 

    https://github.com/AppkubeCloud/appkube-api-specs/blob/main/awsx-api.md

| S.No | Panel Name          | Description                                           | Data output                                 | Commands                                                                                                                                                                            | 
|------|-----------------------|-------------------------------------------------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1    | cpu_utilization_panel | This will get the specific EKS cluster cpu utilization panel data in hybrid structure  | Percentage(%)                                    | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |    |
| 2    | memory_utilization_panel  | This will get the specific EKS cluster memory utilization panel data in hybrid structure | Bytes                                   | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="memory_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |
| 3    | storage_utilization_panel  | This will get the specific EKS cluster storage utilization panel data in hybrid structure | Bytes                                      |go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="storage_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |
| 4    | network_utilization_panel   |This will get the specific EKS cluster network utilization panel data in hybrid structure | Bytes                                 | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="network_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |
 | 5    | cpu_requests_panel    | Collect Information about specific cloud elements -Run Queries | Timeseries | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_requests_panel " --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""--startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"
 |
 | 6    | allocatable_cpu_panel   | Collect Information about specific cloud elements -Run Queries | Timeseries                                 | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="allocatable_cpu_panel" --elementType="ContainerInsights" --responseType=json --startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"
 |
 | 7    | cpu_limits_panel   | Collect Information about specific cloud elements -Run Queries | Timeseries                                 | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_limits_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |



# all subcommands and options for ECS
<!-- getElementDetails references --> 

All the supported subcommands and there source code locations are mentiioned in 

    https://github.com/AppkubeCloud/appkube-api-specs/blob/main/awsx-api.md

| S.No | Panel Name          | Description                                           | Data output                                 | Commands                                                                                                                                                                            | 
|------|-----------------------|-------------------------------------------------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1    | cpu_utilization_panel | This will get the specific EKS cluster cpu utilization panel data in hybrid structure  | Percentage(%)                                    | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |    |
| 2    | memory_utilization_panel  | This will get the specific EKS cluster memory utilization panel data in hybrid structure  | Bytes                                   | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="memory_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |
| 3    | storage_utilization_panel  | This will get the specific EKS cluster storage utilization panel data in hybrid structure | Bytes                                      |go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="storage_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |
| 4    | network_utilization_panel   | This will get the specific EKS cluster network utilization panel data in hybrid structure | Bytes                                 | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX> --crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="network_utilization_panel" --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |
 | 5    | cpu_requests_panel    | Collect Information about specific cloud elements -Run Queries | Timeseries | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_requests_panel " --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""--startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"
 |
 | 6    | cpu_limits_panel    | Collect Information about specific cloud elements -Run Queries | Timeseries                                 | go run awsx-getelementdetails.go  --zone=us-east-1 --externalId=<afreen1309XXX>--crossAccountRoleArn=<afreen1309XXX> --clusterName="myclustTT" --query="cpu_limits_panel " --elementType="ContainerInsights" --responseType=json --startTime="" --endTime=""
 |


 