
- [awsx-getelementdetails](#awsx-getelementdetails)
- [subcommands and options for EC2](#subcommands-and-options-for-ec2)
- [subcommands and options for EKS](#subcommands-and-options-for-eks)
- [subcommands and options for ECS](#subcommands-and-options-for-ecs)



# awsx-getelementdetails
It implements the awsx plugin getElementDetails 

This subcommand will need to take care for all the cloud elements and for every element, we need to support the composite method like network_utilization_panel. So , we can keep a single repo for the subcommand and keep separate folders for the different element handlers.
# subcommands and options for EC2

| S.No | CLI Spec|  Description                           
|------|----------------|----------------------|
| 1    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="AWS/EC2" --query="ec2-config-data"  | This will get the specific EC2 instance config data |

| 2    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="AWS/EC2" --query="cpu_utilization_panel" --instanceID="i-05e4e6757f13da657" --responseType=json --startTime="" --endTime="" | This will get the specific EC2 instance cpu utilization panel data in hybrid structure |

| 3    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="AWS/EC2" --query="memory_utilization_panel" --instanceID="i-05e4e6757f13da657" --responseType=json --startTime="" --endTime="" | This will get the specific EC2 instance memory utilization panel data in hybrid structure |

| 4   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="AWS/EC2" --query="storage_utilization_panel" --instanceID="i-05e4e6757f13da657" --RootVolumeId="i-05e4e6757f13da657" --EBSVolume1Id="vol-0db5984a7f9d77c4d" --EBSVolume2Id="vol-0e065bd2535df7a54" --responseType=json --startTime="" --endTime=""|
| This will get the specific EC2 instance storage utilization panel data in hybrid structure|

| 5   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="AWS/EC2" --query="network_utilization_panel" --instanceID="i-05e4e6757f13da657" --responseType=json --startTime="" --endTime="" | This will get the specific EC2 instance network utilization panel data in hybrid structure |

# awsx-getelementdetails
It implements the awsx plugin getElementDetails 

This subcommand will need to take care for all the cloud elements and for every element, we need to support the composite method like cpu_utilization,network_utilization_panel. So , we can keep a single repo for the subcommand and keep separate folders for the different element handlers.

# subcommands and options for EKS

| 1   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="cpu_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime="" | This will get the specific EKS cluster cpu utilization panel data in hybrid structure |

| 2   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="memory_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime="" | This will get the specific EKS cluster memory utilization panel data in hybrid structure |

| 3  | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="storage_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime="" | This will get the specific EKS cluster storage utilization panel data in hybrid structure |

| 4  | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="network_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime=""  | This will get the specific EKS cluster network utilization panel data in hybrid structure |

| 5   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="cpu_requests_panel"  --clusterName="myclustTT" --responseType=json   --startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"   | This will get the specific EKS cpu requests panel data in hybrid structure |

| 6   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="allocatable_cpu_panel"  --clusterName="myclustTT" --responseType=json   --startTime="2024-02-01T00:00:00Z" --endTime="2024-02-01T23:59:59Z"    | This will get the specific EKS allocatable cpu requests panel data in hybrid structure |

| 7  | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="cpu_limits_panel"  --clusterName="myclustTT" --responseType=json   --startTime="" --endTime=""    | This will get the specific EKS cluster cpu limits  panel data in hybrid structure |


# awsx-getelementdetails
It implements the awsx plugin getElementDetails 

This subcommand will need to take care for all the cloud elements and for every element, we need to support the composite method like cpu_utilization,memory_utilization,storage-utilization,network_utilization_panel. So , we can keep a single repo for the subcommand and keep separate folders for the different element handlers.

# subcommands and options for EKS

| 1   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="cpu_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime=""  | This will get the specific ECS cluster cpu utilization panel data in hybrid structure |

| 2   | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="memory_utilization_panel"  --clusterName="myclustTT" --responseType=json   --startTime="" --endTime="" | This will get the specific ECS memory utilization panel data in hybrid structure |

| 3  | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights" --query="storage_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime="" | This will get the specific ECS cluster storage utilization panel data in hybrid structure |

| 4  | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType="ContainerInsights  --query="network_utilization_panel" --clusterName="myclustTT" --responseType=json   --startTime="" --endTime=""  | This will get the specific ECS cluster network utilization panel data in hybrid structure |











