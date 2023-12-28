
- [awsx-getelementdetails](#awsx-getelementdetails)
- [subcommands and options for EC2](#subcommands-and-options-for-ec2)

# awsx-getelementdetails
It implements the awsx plugin getElementDetails 

This subcommand will need to take care for all the cloud elements and for every element, we need to support the composite method like network_utilization_panel. So , we can keep a single repo for the subcommand and keep separate folders for the different element handlers.
# subcommands and options for EC2

| S.No | CLI Spec|  Description                           
|------|----------------|----------------------|
| 1    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="ec2-config-data"  | This will get the specific EC2 instance config data |
| 2    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="cpu_utilization_panel"  | This will get the specific EC2 instance cpu utilization panel data in hybrid structure |
| 3    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="storage_utilization_panel" | This will get the specific EC2 instance storage utilization panel data in hybrid structure|
| 4    | awsx --vaultURL=vault.synectiks.net getElementDetails --elementId="1234" --elementType=EC2 --query="network_utilization_panel"  | This will get the specific EC2 instance network utilization panel data in hybrid structure |

