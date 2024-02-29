
## AWS cloudwatch-metric-cli Documentation

### Overview

cli to monitors AWS resources using cloudwatch metric queries. It is written in Go and customizable with various parameters.

### Prerequisites
- Go installed.

### Command Details
```
go run awsx-getelementdetails.go  --vaultUrl=<afreenXXXXXXX1309> --elementId=9321 --query="cpu_utilization_panel" --elementType="EC2" --responseType=json --startTime=2023-12-01T00:00:00Z --endTime=2023-12-02T23:59:59Z

```
### Command Parameter:
- --crossAccountRoleArn: AWS IAM role ARN for cross-account access.
- -cloudWatchQueries: JSON array of CloudWatch queries.
       mandatory paramters of cloudWatchQueries
            1. elementId
            2. elementType
            3. query
            4. responseType
            5. startTime and End Time
    
### Logic to get GLOBAL_AWS_SECRETS (access/secret key) in cli: 
        Since we are only passing crossAccountRoleArn, we need GLOBAL_AWS_SECRETS (access/secret key) from vault. It can be retrieved by two ways explaind below: 
            1. make vault call with static key (GLOBAL_AWS_SECRETS)
            2. If vault is not available, get the GLOBAL_AWS_SECRETS from environment variable
            3. If GLOBAL_AWS_SECRETS not found in environment variable, program should exit with error - clien connection could not be established. access/secret key not found

# Proposed changes
        1. awsx-common 
        GLOBAL_AWS_SECRETS logic described in above para should be implemented in awsx-common layer. awsx-common is responsible to do authentication and provide aws connection based on aws element types (e.g. cloudwatch-metric etc..)
        2. Appkube-cloud-datasource
        Once cli changes are done and validated the above command, we need to make following changes in Appkube-cloud-datasource
            2.1 crossRoleArn for the elementId 
	        2.2 Full and final query 
                    NOTE: since appkube-cloud-datasource is able to make a json with all the required query params, we don't need any tranformation in this json in api layer. So pass this query json to cli as it is. cli will parse this json to make cloudwatch-query-input

# Integration with awsx-metric api
    http://<server>:port/awsx-metrics


