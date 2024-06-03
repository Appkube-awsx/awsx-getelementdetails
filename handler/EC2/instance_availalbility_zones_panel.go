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

type ZoneInfo struct {
	AvailabilityZone string  `json:"availability_zone"`
	InstanceCount    int     `json:"instance_count"`
	Percentage       float64 `json:"percentage"`
}

type Summary struct {
	TotalInstances int        `json:"total_instances"`
	Zones          []ZoneInfo `json:"zones"`
}

var AwsxInstanceAvailabilityZoneCmd = &cobra.Command{
	Use:   "instance_availalbility_zones_panel",
	Short: "gets total insatances and availability zones percentage",
	Long:  `Command to get total insatances and availability zones percentage`,

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
			jsonResp, cloudwatchMetricResp, err := GetInstanceAvailabilityZonesData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting instance availalbility zones data : ", err)
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

func GetInstanceAvailabilityZonesData(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (string, Summary, error) {

	instances_details, err := GetEc2InstanceDetails(cmd, clientAuth, ec2Client)
	if err != nil {
		log.Println("Error in getting instance availalbility zones: ", err)
		return "", Summary{}, err
	}

	log.Println("Data: ", instances_details)

	jsonString, err := json.MarshalIndent(instances_details, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal summary to JSON: %v", err)
	}

	return string(jsonString), instances_details, nil
}

func GetEc2InstanceDetails(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (Summary, error) {
	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}

	input := &ec2.DescribeInstancesInput{}

	result, err := ec2Client.DescribeInstances(input)
	if err != nil {
		log.Fatalf("Failed to describe instances: %v", err)
	}

	azCount := make(map[string]int)
	instanceCount := 0

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceCount++
			az := *instance.Placement.AvailabilityZone
			azCount[az]++
		}
	}

	var zones []ZoneInfo
	for az, count := range azCount {
		percentage := (float64(count) / float64(instanceCount)) * 100
		zones = append(zones, ZoneInfo{
			AvailabilityZone: az,
			InstanceCount:    count,
			Percentage:       percentage,
		})
	}

	summary := Summary{
		TotalInstances: instanceCount,
		Zones:          zones,
	}

	return summary, nil

}

func init() {
	comman_function.InitAwsCmdFlags(AwsxInstanceAvailabilityZoneCmd)
}
