package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// InstanceCounts struct to hold the counts of running and stopped instances
type InstanceCounts struct {
	RunningInstances int
	StoppedInstances int
}

var AwsxEc2InstanceCountCmd = &cobra.Command{
	Use:   "instance_count_panel",
	Short: "Get instance count metrics data",
	Long:  `Command to get instance count metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running instance count panel command")

		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			if err := cmd.Help(); err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, err := GetInstanceCountPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting instance count data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(jsonResp)
			} else {
				fmt.Println(jsonResp)
			}

		}
	},
}

func GetInstanceCountPanel(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (string, error) {
	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}

	instanceCounts := &InstanceCounts{}
	filters := []*ec2.Filter{
		{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String(ec2.InstanceStateNameRunning)},
		},
	}

	runningParams := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	runningResp, err := ec2Client.DescribeInstances(runningParams)
	if err != nil {
		return "", fmt.Errorf("failed to describe running instances: %v", err)
	}

	instanceCounts.RunningInstances = len(runningResp.Reservations)

	filters[0].Values = []*string{aws.String(ec2.InstanceStateNameStopped)}
	stoppedParams := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	stoppedResp, err := ec2Client.DescribeInstances(stoppedParams)
	if err != nil {
		return "", fmt.Errorf("failed to describe stopped instances: %v", err)
	}

	instanceCounts.StoppedInstances = len(stoppedResp.Reservations)
	jsonResp, _ := json.Marshal(instanceCounts)
	return string(jsonResp), nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2InstanceCountCmd)
}
