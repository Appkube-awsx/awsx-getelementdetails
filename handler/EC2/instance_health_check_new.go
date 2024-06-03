package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// HealthCheck struct to hold instance health data
type HealthCheck struct {
	HealthyInstancesCount   int `json:"HealthyInstancesCount"`
	UnhealthyInstancesCount int `json:"UnhealthyInstancesCount"`
}

var AwsxEc2InstanceHealthCheckNewCmd = &cobra.Command{
	Use:   "instance_health_check_new",
	Short: "Get EC2 instance health check data",
	Long:  `Command to get EC2 instance health check data including counts of healthy and unhealthy instances`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running instance health check command...")
		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			healthCheck, err := GetInstanceHealthCheckNew(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network traffic data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println("This cli is for only json output")
			} else {
				jsonResp, err := json.Marshal(&healthCheck)
				if err != nil {
					log.Println("Error marshalling network traffic data: ", err)
					return
				}
				fmt.Println(string(jsonResp))
			}
		}
	},
}

func GetInstanceHealthCheckNew(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (string, error) {
	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}

	allInstanceStatuses := []*ec2.InstanceStatus{}
	nextToken := ""

	for {
		params := &ec2.DescribeInstanceStatusInput{}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}

		resp, err := ec2Client.DescribeInstanceStatus(params)
		if err != nil {
			return "", fmt.Errorf("failed to describe instance status: %v", err)
		}
		//fmt.Println(resp)
		allInstanceStatuses = append(allInstanceStatuses, resp.InstanceStatuses...)
		if resp.NextToken == nil {
			break
		}
		nextToken = *resp.NextToken
	}

	// Calculate counts of healthy and unhealthy instances
	healthCheck := &HealthCheck{}
	for _, status := range allInstanceStatuses {
		if *status.InstanceStatus.Status == "ok" && *status.SystemStatus.Status == "ok" {
			healthCheck.HealthyInstancesCount++
		} else {
			healthCheck.UnhealthyInstancesCount++
		}
	}
	jsonResp, _ := json.Marshal(healthCheck)
	return string(jsonResp), nil
}

func init() {
	AwsxEc2InstanceHealthCheckNewCmd.PersistentFlags().String("region", "", "AWS region to use")
}
