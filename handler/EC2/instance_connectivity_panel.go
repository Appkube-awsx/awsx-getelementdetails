package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

type Instance struct {
	InstanceName     string
	InstanceType     string
	AvailabilityZone string
	Status           string
	DisconnectReason string
}

var AwsxNetworkConnectivityCmd = &cobra.Command{
	Use:   "instance_connectivity_panel",
	Short: "gets connectivity of instances",
	Long:  `Command to get connectivity of instances`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
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
			jsonResp, cloudwatchMetricResp, err := GetConnectivityData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting connectivity of instances : ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetConnectivityData(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (string, []Instance, error) {

	instances_details, err := GetConnectivityDetails(cmd, clientAuth, ec2Client)
	if err != nil {
		log.Println("Error in connectivity of instances: ", err)
		return "", nil, err
	}

	log.Println("Data: ", instances_details)

	jsonString, err := json.MarshalIndent(instances_details, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal summary to JSON: %v", err)
	}

	return string(jsonString), instances_details, nil
}

func GetConnectivityDetails(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) ([]Instance, error) {
	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}

	input := &ec2.DescribeInstancesInput{}

	result, err := ec2Client.DescribeInstances(input)
	if err != nil {
		log.Fatalf("Failed to describe instances: %v", err)
	}

	var instances []Instance

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var instanceName string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceName = *tag.Value
					break
				}
			}

			var status string
			var disconnectReason string

			if *instance.State.Name == "running" {
				status = "available"
			} else if *instance.State.Name == "stopped" {
				status = "Unavailable"
				if instance.StateReason != nil {
					disconnectReason = *instance.StateReason.Message
				}
			}

			instanceData := Instance{
				InstanceName:     instanceName,
				InstanceType:     *instance.InstanceType,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				Status:           status,
				DisconnectReason: disconnectReason,
			}
			instances = append(instances, instanceData)
		}
	}

	return instances, nil

}

func init() {
	comman_function.InitAwsCmdFlags(AwsxNetworkConnectivityCmd)

}
