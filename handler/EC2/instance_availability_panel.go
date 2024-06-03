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

type InstanceDetails struct {
	InstanceID       string
	AvailabilityZone string
	InstanceName     string
	State            string
	StateReason      string
}

var AwsxInstanceAvailabilityCmd = &cobra.Command{
	Use:   "instance_availalbility_panel",
	Short: "gets total insatances and its State",
	Long:  `Command to get total insatances and its State`,

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
			jsonResp, cloudwatchMetricResp, err := InstanceAvailability(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting instance availalbility data : ", err)
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

func InstanceAvailability(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (string, []InstanceDetails, error) {

	instances_details, err := GetInstanceAvailabilityDetails(clientAuth, ec2Client)
	if err != nil {
		log.Println("Error in getting instance availalbility: ", err)
		return "", nil, err
	}

	log.Println("Data: ", instances_details)

	jsonString, err := json.MarshalIndent(instances_details, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal summary to JSON: %v", err)
	}

	return string(jsonString), instances_details, nil
}

func GetInstanceAvailabilityDetails(clientAuth *model.Auth, ec2Client *ec2.EC2) ([]InstanceDetails, error) {
	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}

	input := &ec2.DescribeInstancesInput{}

	result, err := ec2Client.DescribeInstances(input)
	if err != nil {
		log.Fatalf("Failed to describe instances: %v", err)
	}

	var instances []InstanceDetails

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// Get the instance name tag
			var instanceName string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceName = *tag.Value
					break
				}
			}

			// Get the state reason message if available
			var stateReason string
			if instance.StateReason != nil {
				stateReason = *instance.StateReason.Message
			}

			// Create an InstanceInfo object and append it to the slice.
			instances = append(instances, InstanceDetails{
				InstanceID:       *instance.InstanceId,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				InstanceName:     instanceName,
				State:            *instance.State.Name,
				StateReason:      stateReason,
			})
		}
	}
	return instances, nil

}

func init() {
	comman_function.InitAwsCmdFlags(AwsxInstanceAvailabilityCmd)
}
