- [awsx-getelementdetails](#awsx-getelementdetails)
- [ui-analysys-and listing-methods](#ui-analysys-and-listing-methods)
  - [cpu\_utilization\_panel](#cpu_utilization_panel)
  - [memory\_utilization\_panel](#memory_utiization_panel)
  - [storage\_utilization\_panel](#storage_utiization_panel)
  - [network\_utilization\_panel](#network_utiization_panel)
  - [cpu\_utilization\_graph\_panel](#cpu_utilizaion_graph_panel)
  - [cpu\_reservation\_panel](#cpu_reservation_panel)   
  - [cpu\_usage\_sys\_panel](#cpu_usage_sys_panel)
  - [cpu\_usage\_nice\_panel](#cpu_usage_nice_panel)
  - [memory\_utilization\_graph\_panel](#memory_utilization_graph_panel)
  - [memory\_reservation\_panel](#memory_reservation_panel)
  - [container\_memory\_usage\_panel](#container_memory_usage_panel)
  - [memory\_overtime\_panel](#memory_overtime_panel)
  - [volume\_readBytes\_panel](#volume_readbytes_panel)
  - [volume\_writeBytes\_panel](#volume_writebytes_panel)
  - [I/O\_bytes\_panel](#input_output_bytes_panel)
  - [disk\_available\_panel](#disk_available_panel) 
  - [net\_inBytes\_panel](#net_inbytes_panel)
  - [net\_outBytes\_panel](#net_outbytes_panel)
  - [net\_ReceiveInBytes\_panel](#net_receiveinbytes_panel)
  - [net\_TransmitInbytes\_panel](#net_transmitinbytes_panel)
  - [net\_RxInBytes\_panel](#net_rxinbytes_panel) 
  - [net\_TxInBytes\_panel](#net_txinbytes_panel)
  
 
- [list of subcommands and options for ECS](#list-of-subcommands-and-options-for-ecs)

list of subcommands and options for EC2
 
# awsx-getelementdetails
It implements the awsx plugin getElementDetails
 
# ui-analysys-and listing-methods
![Alt text](ecs_screen1.png)
1. cpu_utilization_panel
2. memory_utilization_panel
3. storage_utilization_panel
4. network_utilization_panel
5. cpu_utilization_panel
6. cpu_usage_idle_panel
7. cpu_reservation_panel
8. cpu_usage_nice_panel
9.  memory_utilization_panel
10. memory_reservation_panel
11. memory_usage_panel
12. memory_overtime_panel
13. volume_readBytes_panel
14. volume_writeBytes_panel
15. I/O_Bytes_panel
16. disk_available_panel
17. net_inBytes_panel
18. net_outBytes_panel
19. net_ReceiveInBytes_panel
20. net_transmitInBytes_panel
21. net_RxInBytes_panel
22. net_TxInBytes_panel

_
# ui-analysys-and listing-methods

1. cpu_utilization_panel
![Alt text](ecs_screen1.png)

## cpu_utiization_panel


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="cpu_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="cpu_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=cpu_utilization_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



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

2. memory_utilization_panel
![Alt text](ecs_screen1.png)

## memory_utiization_panel


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="memory_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="memory_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=memory_utilization_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
2.  Memory utilization
{
    CurrentUsage:15GB,
    AverageUsage:25GB,
	MaxUsage:50GB
}


**Algorithm/ Pseudo Code**

**Algorithm:** 
- MemoryUtilization - Write a custom metric for memory utilization, where we shall write a program for current, avg and max.

 **Pseudo Code:** 

 
 
 # ui-analysys-and listing-methods

3. storage_utilization_panel 
![Alt text](ecs_screen1.png)

## storage_utiization_panel

**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="storage_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="storage_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=storage_utilization_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



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

4. network_utilization_panel 
![Alt text](ecs_screen1.png)

## network_utiization_panel

**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="storage_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="storage_utilization_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=storage_utilization_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


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
 
 
 # ui-analysys-and listing-methods

## cpu_utilizaion_graph_panel

5. cpu_utilization_graph_panel
![Alt text](ecs_screen2.png)


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="cpu_utilization_graph_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="cpu_utilization_graph_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=cpu_utilization_graph_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
5. CPU utilizaion  graph panel

	-CPUUtilization



**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU utilization graph  -Fire a cloudwatch query for cpu_utilization_graph_panel, using metric CPUUtilization.

 **Pseudo Code:** 
 
 # ui-analysys-and listing-methods


## cpu_reservation_panel


![Alt text](ecs_screen2.png)
6. cpu_reservation_panel


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="cpu_reservation_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="cpu_reservation_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=cpu_reservation_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
6. CPU reservation panel

	-CPUReservation


**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU reservation  -Fire a cloudwatch query for cpu_reservation_panel, using metric CPUReservation.

 **Pseudo Code:** 
 
 # ui-analysys-and listing-methods

## cpu_usage_sys_panel

7. cpu_usage_system_panel
![Alt text](ecs_screen2.png)

 
**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="cpu_usage_system_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="cpu_usage_system_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=cpu_usage_system_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
7. CPU usage system panel

	-cpu_usage_system


**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU usage system  -Fire a cloudwatch query for cpu_usage_system_panel, using metric cpu_usage_system.

 **Pseudo Code:** 
 
# ui-analysys-and listing-methods

## cpu_usage_nice_panel

8. cpu_usage_nice_panel
![Alt text](ecs_screen2.png)


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="cpu_usage_nice_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="cpu_usage_nice_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=cpu_usage_nice_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
8. CPU usage nice panel

	-cpu_usage_nice



**Algorithm/ Pseudo Code**

**Algorithm:** 
- CPU usage nice  -Fire a cloudwatch query for cpu_usage_nice_panel, using metric cpu_usage_nice.

 **Pseudo Code:** 
 
 
 # ui-analysys-and listing-methods

##  memory_utilization_graph_panel

9. memory_utilization_graph_panel
![Alt text](ecs_screen3.png)


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="memory_utilization_graph_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="memory_utilization_graph_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=memory_utilization_graph_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



 
**Desired Output in json / graph format:**
9. memory utilization graph  panel

	-MemoryUtilizaion_graph_panel



**Algorithm/ Pseudo Code**

**Algorithm:** 
- Memory utilization graph panel  -Fire a cloudwatch query for memory_utilization_graph_panel, using metric MemoryUtilization_graph_panel.

 **Pseudo Code:** 
 
 
 # ui-analysys-and listing-methods

##  memory_reservation_panel

10. memory_reservation_panel
![Alt text](ecs_screen3.png)


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="memory_reservation_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="memory_reservation_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=memory_reservation_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
10. memory reservation panel

	-MemoryReservation_panel



**Algorithm/ Pseudo Code**

**Algorithm:** 
- Memory reservation panel  -Fire a cloudwatch query for memory_reservation_panel, using metric memory_resevation_panel.

 **Pseudo Code:** 
 
 # ui-analysys-and listing-methods

##  container_memory_usage_panel

11. container_memory_usage_panel
![Alt text](ecs_screen3.png)


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="container_memory_usage_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="container_memory_usage_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=container_memory_usage_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
11. container_memory_usage_panel

	-mem_used_panel



**Algorithm/ Pseudo Code**

**Algorithm:** 
- container Memory used panel  -Fire a cloudwatch query for memory_usage_panel, using metric container memory_usage_panel.

 **Pseudo Code:** 
 
 
 # ui-analysys-and listing-methods

##  available_memory_overtime_panel

12. available_available_memory_overtime_panel
![Alt text](ecs_screen3.png)



**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="available_memory_overtime_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="available_memory_overtime_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=available_memory_overtime_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
12. available memory overtime panel

	-available memory_overtime_panel



**Algorithm/ Pseudo Code**

**Algorithm:** 
- Available Memory_overtime panel  -Fire a cloudwatch query for available memory_overtime_panel, using metric memory_overtime_panel.

 **Pseudo Code:**  

 
# ui-analysys-and listing-methods

##  volume_readBytes_panel

13. volume_readBytes_panel
![Alt text](ecs_screen4.png)

**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="volume_readBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="volume_readBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=volume_readBytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
13. volume_readBytes panel

	-volume_readBytes_panel



**Algorithm/ Pseudo Code**

**Algorithm:** 
- volume readBytes panel  -Fire a cloudwatch query for volume_readBytes_panel, using metric volume_readBytes_panel.

 **Pseudo Code:**  
 

 # ui-analysys-and listing-methods

 ##  volume_writebytes_panel

14. volume_writeBytes_panel
![Alt text](ecs_screen4.png)




**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="volume_writeBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="volume_writeBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=volume_writeBytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
14. volume_writeBytes panel

	-volume_writebytes_panel



**Algorithm/ Pseudo Code**

**Algorithm:** 
- volume writeBytes panel  -Fire a cloudwatch query for volume_writeBytes_panel, using metric volume_writeBytes_panel.

 **Pseudo Code:**  
 
 
 # ui-analysys-and listing-methods

##  input_output_bytes_panel

15. I/O_Bytes_panel
![Alt text](ecs_screen4.png)


##  input_output_bytes_panel


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="input_output_bytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="input_output_bytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=input_output_bytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z



**Desired Output in json / graph format:**
15. In/Out Bytes panel

	-in/out bytes_panel
	


**Algorithm/ Pseudo Code**

**Algorithm:** 
- in/Out bytes panel  -Fire a cloudwatch query for disk_used_panel, using metric InBytes, OutBytes.

 **Pseudo Code:**  
 
 # ui-analysys-and listing-methods

##  disk_available_panel

16. disk_available_panel
![Alt text](ecs_screen4.png)



**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="disk_available_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="disk_available_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=disk_available_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?zone=us-east-1&externalId=<afreenxxxx1309>&crossAccountRoleArn=<afreenxxxx1309>&elementType=AWS/EC2&instanceID=i-05e4e6757f13da657&query=disk_available_panel


**Desired Output in json / graph format:**
16. disk_available panel

	-disk_available_panel
	  

**Algorithm/ Pseudo Code**

**Algorithm:** 
- disk available panel  -Fire a cloudwatch query for disk_available_panel, using metric disk_available_panel.

 **Pseudo Code:**  
 
 
# ui-analysys-and listing-methods

##  net_inBytes_panel

17. net\_inBytes\_panel
![Alt text](ecs_screen5.png)


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="net_inBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="net_inBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=net_inBytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z




**Desired Output in json / graph format:**
17. network_inBytes_panel

	-network_in_panel
	

**Algorithm/ Pseudo Code**

**Algorithm:** 
- network_inBytes panel  -Fire a cloudwatch query for network_inBytes_panel, using metric NetworkBytesIn.

 **Pseudo Code:**  
 
 
 
 # ui-analysys-and listing-methods

##  net_outBytes_panel

18. net_outBytes_panel
![Alt text](ecs_screen5.png)



**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="net_outBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="net_outBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=net_outBytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
18. network_outBytes_panel

	-NetworkBytesOut
	

**Algorithm/ Pseudo Code**

**Algorithm:** 
- network_outBytes panel  -Fire a cloudwatch query for network_outBytes_panel, using metric NetworkBytesOut.

 **Pseudo Code:**  
 
 
 # ui-analysys-and listing-methods

19. net\_ReceiveInBytes\_panel
![Alt text](ecs_screen5.png)


##  net\_RecieveInBytes\_panel


**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="net_recieveInBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="net_recieveInBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=net_recieveInBytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
19. network_ReceiveInBytes_panel

	-network_ReceiveInBytes_panel
	

**Algorithm/ Pseudo Code**

**Algorithm:** 
- network_ReceiveInBytes panel  -Fire a cloudwatch query for network_ReceiveInBytes_panel, using metric network_ReceiveInBytes_panel.

 **Pseudo Code:**  
 
 # ui-analysys-and listing-methods

##  net\_transmitInBytes\_panel

20. net\_transmitInBytes\_panel
![Alt text](ecs_screen5.png)



**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="net_transmitInBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="net_transmitInBytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=net_transmitInBytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
20. network_transmitInBytes_panel

	-network_transmitInBytes_panel
	

**Algorithm/ Pseudo Code**

**Algorithm:** 
- network_transmitInBytes panel  -Fire a cloudwatch query for network_transmitInBytes_panel, using metric NetworkBytesIn_panel.

 **Pseudo Code:**  
 
 
 # ui-analysys-and listing-methods

##  net\_RxInBytes\_panel

21. net\_RxInBytes\_panel
![Alt text](ecs_screen5.png)



**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="net_rxinbytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="net_rxinbytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=net_rxinbytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
21. network_transmitInBytes_panel

	-network_transmitInBytes_panel
	

**Algorithm/ Pseudo Code**

**Algorithm:** 
- network_RxInBytes panel  -Fire a cloudwatch query for network_RxInBytes_panel, using metric NetworkBytesIn_panel.

 **Pseudo Code:**  
 
 # ui-analysys-and listing-methods

##  net\_TxInBytes\_panel

22. net\_TxInBytes\_panel
![Alt text](ecs_screen5.png)





**called from subcommand**

go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=900001 --query="net_txinbytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**called from maincommand**

awsx --vaultUrl=<afreenXXXXXXX1309> --elementId=90001  --query="net_txinbytes_panel" --elementType="ECS" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z


**Called from API**

http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenXXXX>&elementId=900001&elementType=ECS&query=net_txinbytes_panel&responseType=json&startTime=2023-12-01T00:00:00Z&endTime=2023-12-02T23:59:59Z


**Desired Output in json / graph format:**
22. network_TxInBytes_panel

	-network_TxInBytes_panel
	

**Algorithm/ Pseudo Code**

**Algorithm:** 
- network_TxInBytes panel  -Fire a cloudwatch query for network_TxInBytes_panel, using metric NetworkBytesIn_panel.

 **Pseudo Code:**  


 
# list of subcommands and options for ECS
 
| S.No | CLI Spec|  Description                          
|------|----------------|----------------------|
| 1    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="900001" --elementType=EC2 --query="cpu_utilization_panel"  | This will get the specific EC2 instance cpu utilization panel data in hybrid structure |
| 2    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="900001" --elementType=EC2 --query="memory_utilization_panel" | This will get the specific EC2 instance memory utilization panel data in hybrid structure|
| 3    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="900001" --elementType=EC2 --query="storage_utilization_panel"  | This will get the specific EC2 instance storage utilization panel data in hybrid structure |
| 4    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="900001" --elementType=EC2 --query="network_utilization_panel"  | This will get th1e specific EC2 instance network utilization panel data in hybrid structure |




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

