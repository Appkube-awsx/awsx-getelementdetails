package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// AwsxEC2ActiveInstanceCountCmd defines the command for retrieving the count of active EC2 instances
var AwsxEC2ActiveInstanceCountCmd = &cobra.Command{
	Use:   "active_instances_count",
	Short: "Retrieve count of active EC2 instances",
	Long:  `Command to retrieve count of active (running) EC2 instances`,

	Run: func(cmd *cobra.Command, args []string) {
		authFlag, clientAuth, err := HandleAuths(cmd)
		if err != nil {
			log.Println("Error during authentication:", err)
			return
		}

		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			activeinstance, err := GetEC2ActiveInstanceCount(clientAuth)
			if err != nil {
				log.Println("Error getting EC2 instance summary:", err)
				return
			}

			if responseType == "frame" {
				fmt.Println(activeinstance)
			} else {
				printResp := fmt.Sprintf("Active Instance Count: %d", activeinstance)
				fmt.Println(printResp)
			}
		}
	},
}

// GetEC2ActiveInstanceCount retrieves the count of active (running) EC2 instances
func GetEC2ActiveInstanceCount(clientAuth *model.Auth) (int, error) {
	svc := awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)

	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Fatalf("Error describing EC2 instances: %v", err)
	}

	activeInstanceCount := 0
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if *instance.State.Name == "running" {
				activeInstanceCount++
			}
		}
	}

	return activeInstanceCount, nil
}

// handleAuths handles the authentication process
func HandleAuths(cmd *cobra.Command) (bool, *model.Auth, error) {
	authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
	if err != nil {
		return false, nil, err
	}
	return authFlag, clientAuth, nil
}

func init() {
	AwsxEC2ActiveInstanceCountCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
