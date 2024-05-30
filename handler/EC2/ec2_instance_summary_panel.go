package EC2

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type InstanceSummary struct {
	InstanceName string
	InstanceID   string
	InstanceType string
	AvailabilityZone string
	Status       string
}

var AwsxEC2InstanceSummaryCmd = &cobra.Command{
	Use:   "ec2_instance_summary_panel",
	Short: "Retrieve EC2 instance summary data",
	Long:  `Command to retrieve EC2 instance summary data`,

	Run: func(cmd *cobra.Command, args []string) {
		authFlag, clientAuth, err := handleAuths(cmd)
		if err != nil {
			log.Println("Error during authentication:", err)
			return
		}

		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			instanceSummaries, printResp, err := GetEC2InstanceSummaryPanel(clientAuth)
			if err != nil {
				log.Println("Error getting EC2 instance summary:", err)
				return
			}

			if responseType == "frame" {
				fmt.Println(instanceSummaries)
			} else {
				fmt.Println(printResp)
			}
		}
	},
}

func GetEC2InstanceSummaryPanel(clientAuth *model.Auth) ([]InstanceSummary, string, error) {
	instanceSummaries, err := GetEC2InstanceSummary(clientAuth)
	if err != nil {
		return nil, "", err
	}
	formattedTable := printTables(instanceSummaries)
	return instanceSummaries, formattedTable, nil
}

func GetEC2InstanceSummary(clientAuth *model.Auth) ([]InstanceSummary, error) {
	// Use existing AWS client
	svc := awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)

	// Describe EC2 instances
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Fatalf("Error describing EC2 instances: %v", err)
	}

	var instanceSummaries []InstanceSummary
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// Extract instance name from tags
			instanceName := ""
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceName = *tag.Value
					break
				}
			}

			instanceSummary := InstanceSummary{
				InstanceName: instanceName,
				InstanceID:   *instance.InstanceId,
				InstanceType: *instance.InstanceType,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				Status:       *instance.State.Name,
			}

			instanceSummaries = append(instanceSummaries, instanceSummary)
		}
	}

	return instanceSummaries, nil
}

func handleAuths(cmd *cobra.Command) (bool, *model.Auth, error) {
	authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
	if err != nil {
		return false, nil, err
	}
	return authFlag, clientAuth, nil
}

func printTables(instanceSummaries []InstanceSummary) string {
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"Instance Name", "Instance ID", "Instance Type", "Availability Zone", "Status"})

	if instanceSummaries == nil {
		buffer.WriteString("No instance summaries found.")
		return buffer.String()
	}

	for _, summary := range instanceSummaries {
		table.Append([]string{
			summary.InstanceName,
			summary.InstanceID,
			summary.InstanceType,
			summary.AvailabilityZone,
			summary.Status,
		})
	}

	table.Render()

	return buffer.String()
}

func init() {
	AwsxEC2InstanceSummaryCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
