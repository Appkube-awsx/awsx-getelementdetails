package NLB

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type SecurityGroupInfo struct {
	GroupID       string
	InboundRules  string
	OutboundRules string
}

var AwsxSecurityGroupCmd = &cobra.Command{
	Use:   "security_group_panel",
	Short: "Retrieve security group configurations",
	Long:  `Command to retrieve security group configurations`,

	Run: func(cmd *cobra.Command, args []string) {
		authFlag, clientAuth, err := handleAuth(cmd)
		if err != nil {
			log.Println("Error during authentication:", err)
			return
		}

		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			securityGroups, printresp, err := GetSecurityGroupConfigurations(clientAuth)
			if err != nil {
				log.Println("Error getting security group configurations:", err)
				return
			}

			if responseType == "frame" {
				fmt.Println(securityGroups)
			} else {
				fmt.Println(printresp)
			}
		}
	},
}

func GetSecurityGroupConfigurations(clientAuth *model.Auth) ([]SecurityGroupInfo, string, error) {
	// Retrieve security group configurations
	securityGroups, err := DescribeSecurityGroups(clientAuth)
	if err != nil {
		return nil, "", err
	}
	formattedTable := printSecurityGroups(securityGroups)
	return securityGroups, formattedTable, nil
}

func DescribeSecurityGroups(clientAuth *model.Auth) ([]SecurityGroupInfo, error) {
	// Use existing AWS client
	svc := awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)

	// Describe security groups
	describeSGInput := &ec2.DescribeSecurityGroupsInput{}
	describeSGOutput, err := svc.DescribeSecurityGroups(describeSGInput)
	if err != nil {
		log.Fatalf("Error describing security groups: %v", err)
	}

	// Retrieve security group configurations
	var securityGroups []SecurityGroupInfo
	for _, sg := range describeSGOutput.SecurityGroups {
		inboundRules := formatRules(sg.IpPermissions)
		outboundRules := formatRules(sg.IpPermissionsEgress)

		securityGroup := SecurityGroupInfo{
			GroupID:       *sg.GroupId,
			InboundRules:  formatRuleList(inboundRules),
			OutboundRules: formatRuleList(outboundRules),
		}
		securityGroups = append(securityGroups, securityGroup)
	}

	return securityGroups, nil
}

func formatRules(rules []*ec2.IpPermission) []string {
	var formattedRules []string
	for _, rule := range rules {
		fromPort := "traffic-port"
		toPort := "31036"
		protocol := "HTTP"

		if rule.FromPort != nil {
			fromPort = fmt.Sprintf("%d", *rule.FromPort)
		}
		if rule.ToPort != nil {
			toPort = fmt.Sprintf("%d", *rule.ToPort)
		}
		if rule.IpProtocol != nil {
			protocol = *rule.IpProtocol
		}

		formattedRule := fmt.Sprintf("FromPort: %s, ToPort: %s, Protocol: %s", fromPort, toPort, protocol)
		formattedRules = append(formattedRules, formattedRule)
	}
	return formattedRules
}

func formatRuleList(rules []string) string {
	return strings.Join(rules, "\n")
}

func handleAuth(cmd *cobra.Command) (bool, *model.Auth, error) {
	authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
	if err != nil {
		return false, nil, err
	}
	return authFlag, clientAuth, nil
}

func printSecurityGroups(securityGroups []SecurityGroupInfo) string {
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"Security Group ID", "Inbound Rules", "Outbound Rules"})

	if securityGroups == nil {
		buffer.WriteString("No security groups found.")
		return buffer.String()
	}

	for _, sg := range securityGroups {
		table.Append([]string{
			sg.GroupID,
			sg.InboundRules,
			sg.OutboundRules,
		})
	}

	table.Render()

	return buffer.String()
}

func init() {
	AwsxSecurityGroupCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
