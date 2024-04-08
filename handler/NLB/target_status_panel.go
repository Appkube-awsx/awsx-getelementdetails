package NLB

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type TargetStatuss struct {
	TargetID     string
	TargetHealth string
	Reason       string // Reason field added
	Region       string // Region field added
}

var AwsxNLBTargetStatussCmd = &cobra.Command{
	Use:   "new_target_status_panel",
	Short: "Retrieve target status metrics data",
	Long:  `Command to retrieve target status metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		authFlag, clientAuth, err := handleAuths(cmd)
		if err != nil {
			log.Println("Error during authentication:", err)
			return
		}

		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			targetStatuses, printresp, err := GetTargetStatussPanel(clientAuth)
			if err != nil {
				log.Println("Error getting target status:", err)
				return
			}

			if responseType == "frame" {
				fmt.Println(targetStatuses)
			} else {
				fmt.Println(printresp)
			}
		}
	},
}

func GetTargetStatussPanel(clientAuth *model.Auth) ([]TargetStatuss, string, error) {
	// Retrieve target status from AWS NLB target groups
	targetStatuses, err := GetNLBTargetStatus(clientAuth)
	if err != nil {
		return nil, "", err
	}
	formattedTable := printTables(targetStatuses)
	return targetStatuses, formattedTable, nil
}

func GetNLBTargetStatus(clientAuth *model.Auth) ([]TargetStatuss, error) {
	// Use existing AWS client
	svc := awsclient.GetClient(*clientAuth, awsclient.ELBV2_CLIENT).(*elbv2.ELBV2)

	// Describe NLB target groups
	targetGroupsOutput, err := svc.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		log.Fatalf("Error describing target groups: %v", err)
	}

	// Retrieve target status for each target group
	var targetStatuses []TargetStatuss
	for _, tg := range targetGroupsOutput.TargetGroups {
		targetHealthOutput, err := svc.DescribeTargetHealth(&elbv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err != nil {
			log.Printf("Error describing target health for target group %s: %v", *tg.TargetGroupName, err)
			continue
		}

		// Extract relevant information from the response and construct TargetStatus objects
		for _, healthDescription := range targetHealthOutput.TargetHealthDescriptions {
			reason := ""
			if healthDescription.TargetHealth.Reason != nil {
				reason = *healthDescription.TargetHealth.Reason
			}

			// Get the region from the ARN of the target group
			region := getRegionFromARN(*tg.TargetGroupArn)

			targetStatus := TargetStatuss{
				TargetID:     *healthDescription.Target.Id,
				TargetHealth: *healthDescription.TargetHealth.State,
				Reason:       reason, // Set Reason field
				Region:       region, // Set Region field
			}

			// Set custom reason for healthy targets
			if *healthDescription.TargetHealth.State == "healthy" {
				targetStatus.Reason = "Target.HealthChecks"
			}

			targetStatuses = append(targetStatuses, targetStatus)
		}
	}

	return targetStatuses, nil
}

func handleAuths(cmd *cobra.Command) (bool, *model.Auth, error) {
	authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
	if err != nil {
		return false, nil, err
	}
	return authFlag, clientAuth, nil
}

func printTables(targetStatuses []TargetStatuss) string {
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"Target ID", "Target Health", "Description"}) // Update header

	if targetStatuses == nil {
		buffer.WriteString("No target statuses found.")
		return buffer.String()
	}

	for _, status := range targetStatuses {
		// Combine Reason and Region into the Description field
		description := fmt.Sprintf("%s, %s", status.Reason, status.Region)
		table.Append([]string{
			status.TargetID,
			status.TargetHealth,
			description,
		})
	}

	table.Render()

	return buffer.String()
}

func init() {
	AwsxNLBTargetStatussCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}

func getRegionFromARN(arn string) string {
	// ARN format: arn:aws:elasticloadbalancing:<region>:<account-id>:targetgroup/<target-group-name>/<target-group-id>
	// Splitting the ARN string by ":" and extracting the region from the third element
	arnParts := strings.Split(arn, ":")
	if len(arnParts) >= 4 {
		return arnParts[3]
	}
	return ""
}
