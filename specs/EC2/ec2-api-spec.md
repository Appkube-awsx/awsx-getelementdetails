- [awsx-getelementdetails](#awsx-getelementdetails)
- [ui-analysys-and listing-methods](#ui-analysys-and-listing-methods)
  - [cpu\_utilization\_panel](#cpu_utilization_panel)
  - [memory\_utilization\_panel](#memory_utiization_panel)
  - [storage\_utilization\_panel](#storage_utiization_panel)
  - [network\_utilization\_panel](#network_utiization_panel)
 
- [list of subcommands and options for EC2](#list-of-subcommands-and-options-for-ec2)
 
# awsx-getelementdetails
It implements the awsx plugin getElementDetails
 
# ui-analysys-and listing-methods
![Alt text](image.png)
1. cpu_utilization_panel
2. memory_utilization_panel
3. storage_utilization_panel
4. network_utilization_panel

![Alt text](image-1.png)

5. hosted_business_services

![Alt text](image-2.png)

6. ec2_config_data

![Alt text](image-3.png)

7. ec2_sla_data

![Alt text](image-4.png)

8. ec2_cost_and_spike


# ui-analysys-and listing-methods
![Alt text](image.png)
1. cpu_utilization_panel

## cpu_utiization_panel

**called from subcommand**
 
awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EC2 --query="cpu_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="cpu_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EC2, elementId="1234" , query=cpu_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**

1. CPU utilization
{
	CurrentUsage:25%,
	AverageUsage:30%,
	MaxUsage:40%
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU utilization panel - Write a custom metric for cpu utilization, where we shall write a program for current, avg and max.

 **Pseudo Code:**  
 
 

# ui-analysys-and listing-methods
![Alt text](image.png)
2. memory_utilization_panel 

## memory_utiization_panel

**called from subcommand**
 
awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EC2 --query="memory_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="memory_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EC2, elementId="1234" , query=memory_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**

2.  Memory utilization
{
    CurrentUsage:25GB,
    AverageUsage:30GB,
	MaxUsage:40GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- MemoryUtilization - Write a custom metric for memory utilization

 **Pseudo Code:** 

 
 
 # ui-analysys-and listing-methods
![Alt text](image.png)
3. storage_utilization_panel 

## storage_utiization_panel

**called from subcommand**
 
awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EC2 --query="storage_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="storage_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EC2, elementId="1234" , query=storage_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**

3.  Storage utilization
{
    RootVolumeUsage:25GB,
    EBSVolume1Usage:30GB,
	EBSVolume2Usage:40GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- Storage Utilization panel - Write a custom metric for storage utilization, where we shall write a program for root volume usage and ebs disks usage.

 **Pseudo Code:**  
 
 

 # ui-analysys-and listing-methods
![Alt text](image.png)
4. network_utilization_panel 

## network_utiization_panel


**called from subcommand**
 
awsx-getelementdetails --vaultURL=vault.synectiks.net --elementId="1234" --elementType=EC2 --query="network_utilization_panel" --timeRange={}


**called from maincommand**

awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="network_utilization_panel" --timeRange={}

**Called from API**

/awsx-api/getQueryOutput? elementType=EC2, elementId="1234" , query=network_utilization_panel, --timeRange={}


**Desired Output in json / graph format:**

4.Network utilization
{
    Inbound traffic:500Mbps,
    Outbound traffic:200Mbps,
	Data Transferred:10GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- Network utilization panel - Write a custom metric for Network utilization, where we shall write a program for root volume usage and ebs disks usage.

 **Pseudo Code:**
 

 
# list of subcommands and options for EC2
 
| S.No | CLI Spec|  Description                          
|------|----------------|----------------------|
| 1    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="cpu_utilization_panel"  | This will get the specific EC2 instance cpu utilization panel data in hybrid structure |
| 2    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="memory_utilization_panel" | This will get the specific EC2 instance memory utilization panel data in hybrid structure|
| 3    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="storage_utilization_panel"  | This will get the specific EC2 instance storage utilization panel data in hybrid structure |
| 4    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="network_utilization_panel"  | This will get th1e specific EC2 instance network utilization panel data in hybrid structure |



##cpu performance graph panel gettimg data by CWagent



## Acknowledgements

 - [Awesome Readme Templates](https://awesomeopensource.com/project/elangosundar/awesome-README-templates)
 - [Awesome README](https://github.com/matiassingers/awesome-readme)
 - [How to write a Good readme](https://bulldogjob.com/news/449-how-to-write-a-good-readme-for-your-github-project)


## API Reference

#### Get all items

```http
  GET /api/items
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `api_key` | `string` | **Required**. Your API key |

#### Get item

```http
  GET /api/items/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### add(num1, num2)

Takes two numbers and returns the sum.


## Appendix

Any additional information goes here

