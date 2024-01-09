- [awsx-getelementdetails](#awsx-getelementdetails)
- [ui-analysys-and listing-methods](#ui-analysys-and-listing-methods)
  - [cpu\_utilization\_panel](#cpu_utiization_panel)
  - [memory\_utilization\_panel](#memory_utiization_panel)
  - [storage\_utilization\_panel](#storage_utiization_panel)
  - [network\_utilization\_panel](#network_utiization_panel)
  - [cpu\_requests\_panel](#cpu_requests_panel)
  - [allocatable\cpu\_panel](#allocatable_cpu_panel)
  - [cpu\limits\_panel](#cpu_limits_panel)
  - [cpu\_utilization\_graph\_panel](#cpu_utilization_graph_panel)

- [list of subcommands and options for EC2](#list-of-subcommands-and-options-for-eks)

# awsx-getelementdetails
It implements the awsx plugin getElementDetails 

# ui-analysys-and listing-methods
![Alt text](eks-screen-1.png)
1. cpu_utilization_panel 

## cpu_utiization_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="cpu_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=cpu_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**
- CPU utilization
{
	CurrentUsage:25%,
	AverageUsage:30%,
	MaxUsage:40%
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU utilization panel - Write a custom metric for cpu utilization, where we shall write a program for current, avg and max.

 **Pseudo Code:**   




![Alt text](eks-screen-1.png)
2. memory_utilization_panel 

## memory_utiization_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="memory_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=memory_utilization_panel, --timeRange={}


**Desired Output in json format:**
- Memory utilization
{
    CurrentUsage:25GB,
    AverageUsage:30GB,
	MaxUsage:40GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- Memory Utilization panel - Write a custom metric for memory utilization, where we shall write a program for current, avg and max.

**Pseudo Code:**  



# ui-analysys-and listing-methods
![Alt text](eks-screen-1.png)
3. storage_utilization_panel 

## storage_utiization_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="storage_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="storage_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=storage_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**
- Storage utilization
{
    RootVolumeUsage:25GB,
    EBSVolume1Usage:30GB,
	EBSVolume2Usage:40GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- Storage Utilization panel - Write a custom metric for storage utilization, where we shall write a program for root volume usage and ebs disks usage.
    Pseudo Code -

 **Pseudo Code:**

 - [ui-analysys-and listing-methods](#ui-analysys-and-listing-methods)
  - [storage\_utilization\_panel](#network_utilization_panel)
- [list of subcommands and options for EC2](#list-of-subcommands-and-options-for-eks)

# ui-analysys-and listing-methods
![Alt text](eks-screen-1.png)
4. network_utilization_panel 

## network_utiization_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="network_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="network_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=network_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**
- Network utilization
{
    Inbound traffic:500Mbps,
    Outbound traffic:200Mbps,
	Data Transferred:10GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- Network utilization panel - Write a custom metric for Network utilization, where we shall write a program for root volume usage and ebs disks usage.

 **Pseudo Code:**

 - [ui-analysys-and listing-methods](#ui-analysys-and-listing-methods)
  - [cpu\_requests\_panel](#cpu_requests_panel)
- [list of subcommands and options for EC2](#list-of-subcommands-and-options-for-eks)

# ui-analysys-and listing-methods
![Alt text](eks-screen-2.png)


5. cpu_requests_panel 

## cpu_requests_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="cpu_requests_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_requests_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=cpu_requests_panel, --timeRange={}


**Desired Output in  graph format:**
- CPU Requests 


**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU requests panel - Fire a cloudwatch query for CPU requests, using metric namespace as CPU_Requests. 

 **Pseudo Code:**


6. allocatable_cpu_panel 

## allocatable_cpu_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="allocatable_cpu_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="allocatable_cpu_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=allocatable_cpu_panel, --timeRange={}


**Desired Output in  graph format:**
- allocatable_cpu 


**Algorithm/ Pseudo Code**

**Algorithm:** 
- allocatable cpu panel - Fire a cloudwatch query for allocatable cpu, using metric namespace as allocatable_cpu_panel. 

 **Pseudo Code:**
# list of subcommands and options for EKS


7. cpu_limits_panel 

## cpu_limits_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="cpu_limits_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_limits_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=cpu_limits_panel, --timeRange={}


**Desired Output in  graph format:**
- cpu_limits_panel


**Algorithm/ Pseudo Code**

**Algorithm:** 
- cpu_limits_panel - Fire a cloudwatch query for allocatable cpu, using metric namespace as allocatable_cpu_panel. 

 **Pseudo Code:**


8. cpu_utilization_graph_panel 

## cpu_utilization_graph_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="cpu_utilization_graph_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_utilization_graph_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=cpu_utilization_graph_panel, --timeRange={}


**Desired Output in  graph format:**
- cpu_utilization_graph_panel


**Algorithm/ Pseudo Code**

**Algorithm:** 
- cpu_utilization_graph_panel - Fire a cloudwatch query for cpu_utilization_graph_panel, using metric namespace as cpu_utilization_panel. Note - The service name shall be EKS.

 **Pseudo Code:**


 # ui-analysys-and listing-methods
![Alt text](eks-screen-3.png)
9. memory_requests_panel 

## memory_requests_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="memory_requests_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_requests_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=memory_requests_panel, --timeRange={}


**Desired Output in  graph format:**
- memory_requests_panel


**Algorithm/ Pseudo Code**

**Algorithm:** 
- memory_requests_panel - Write a cloudwatch query for memory_requests_panel, where we shall retrieve in graph format.

 **Pseudo Code:**

# ui-analysys-and listing-methods
![Alt text](eks-screen-3.png)

10. memory_limits_panel 

## memory_limits_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="memory_limits_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_limits_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=memory_limits_panel, --timeRange={}


**Desired Output in  graph format:**
- Memory Limits 


**Algorithm/ Pseudo Code**

**Algorithm:** 
- Memory Limits panel - Fire a cloudwatch query for Memory Limits, using metric namespace as memory_limits. 

 **Pseudo Code:**

# ui-analysys-and listing-methods
![Alt text](eks-screen-3.png)
11. allocatable_memory_panel 

## allocatable_memory_panel

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="allocatable_memory_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="allocatable_memory_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=allocatable_memory_panel, --timeRange={}


**Desired Output in  graph format:**
- allocatable_memory_panel 


**Algorithm/ Pseudo Code**

**Algorithm:** 
- allocatable memory panel - Fire a cloudwatch query for allocatable memory, using metric namespace as allocatable_memory. 

 **Pseudo Code:**
# list of subcommands and options for EKS

# ui-analysys-and listing-methods
![Alt text](eks-screen-3.png)
12. memory_utilization_panel 

## memory_utilization_panel 

**called from subcommand**

awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EKS --query="memory_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=memory_utilization_panel, --timeRange={}


**Desired Output in  graph format:**
- memory_utilization_panel 


**Algorithm/ Pseudo Code**

**Algorithm:** 
- memory_utilization_panel - Fire a cloudwatch query for memory_utilization_panel, using metric namespace as memory_utilization. NOTE - The service should be EKS only. 

 **Pseudo Code:**


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_utilization_graph_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EKS, elementId="1234" , query=cpu_utilization_graph_panel, --timeRange={}


**Desired Output in  graph format:**
- cpu_utilization_graph_panel


**Algorithm/ Pseudo Code**

**Algorithm:** 
- cpu_utilization_graph_panel - Fire a cloudwatch query for cpu_utilization_graph_panel, using metric namespace as cpu_utilization_panel. Note - The service name shall be EKS.

 **Pseudo Code:**


| S.No | CLI Spec|  Description                           
|------|----------------|----------------------|
| 1    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_utilization_panel"  | This will get the specific EKS Cluster cpu utilization panel data in hybrid structure |
| 2    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="storage_utilization_panel" | This will get the specific EKS Cluster storage utilization panel data in hybrid structure|
| 3    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="network_utilization_panel"  | This will get the specific EKS Cluster network utilization panel data in hybrid structure |
| 4    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="network_utilization_panel"  | This will get the specific EKS Cluster network utilization panel data in hybrid structure |
| 5    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="CPUrequests"  | This will get the specific EKS Cluster cpu requests to a pod panel data in hybrid structure |
| 6    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="AllocatableCPU"  | This will get the specific EKS Cluster network utilization panel data in hybrid structure |
| 7    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="CPU_limits"  | This will get the specific EKS Cluster cpu limits in a pod, data in hybrid structure |
| 8    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="cpu_utilization_panel"  | This will get the specific EKS Cluster cpu utilization panel data in hybrid structure |
| 9    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_request_panel"  | This will get the specific EKS Cluster memory request panel data in hybrid structure |
| 10    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_limits"  | This will get the specific EKS Cluster memory limits panel over a pod data in hybrid structure |
| 11    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="network_utilization_panel"  | This will get the specific EKS Cluster network utilization panel data in hybrid structure |
| 12    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="memory_utilization_panel"  | This will get the specific EKS memory network utilization panel data in hybrid structure |
| 13    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="disk_utilization_panel"  | This will get the specific EKS Cluster disk utilization(ebs) panel data in hybrid structure |
| 14    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EKS --query="network_in_out_panel"  | This will get the specific EKS Cluster network in & out panel data in hybrid structure |