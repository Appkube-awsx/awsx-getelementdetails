package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var AwsxEc2InstanceStatusCmd = &cobra.Command{

	Use: "instance_status_panel",

	Short: "get instance status metrics data",

	Long: `command to get instance status metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("running from child command")

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
			instanceInfo, err := GetInstanceStatus(cmd, clientAuth)
			if err != nil {
				log.Fatalf("Error getting instance status: %v", err)
			}

			// Print or utilize the instance information
			for _, info := range instanceInfo {
				fmt.Printf("Instance ID: %s, Instance Type: %s, Availability Zone: %s, State: %s, System Checks Status: %s, Custom Alert: %t\n",
					info.InstanceID, info.InstanceType, info.AvailabilityZone, info.State, info.SystemChecksStatus, info.CustomAlert)
			}
		}

	},
}

type InstanceInfo struct {
	InstanceID         string
	InstanceType       string
	AvailabilityZone   string
	State              string
	SystemChecksStatus string
	CustomAlert        bool
}

func GetInstanceStatus(cmd *cobra.Command, clientauth *model.Auth) ([]InstanceInfo, error) {
	// Initialize EC2 client
	ec2Client := awsclient.GetClient(*clientauth, awsclient.EC2_CLIENT).(*ec2.EC2)

	// Initialize CloudWatch client
	cloudWatchClient := awsclient.GetClient(*clientauth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	log.Println("Getting AWS EC2 instance list")

	// Retrieve instance status
	resp, err := ec2Client.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}

	// Define a slice to hold instance information
	var instances []InstanceInfo

	// Populate instance information slice
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceID := aws.StringValue(instance.InstanceId)
			instanceType := aws.StringValue(instance.InstanceType)
			availabilityZone := aws.StringValue(instance.Placement.AvailabilityZone)
			state := aws.StringValue(instance.State.Name)
			systemChecksStatus := getSystemChecksStatus(ec2Client, instanceID)
			hasCustomAlert, err := checkForCustomAlert(cloudWatchClient, instanceID)
			if err != nil {
				log.Printf("Error checking custom alert for instance %s: %v", instanceID, err)
				continue // Skip to the next instance
			}

			// Append instance information to the slice
			instances = append(instances, InstanceInfo{
				InstanceID:         instanceID,
				InstanceType:       instanceType,
				AvailabilityZone:   availabilityZone,
				State:              state,
				SystemChecksStatus: systemChecksStatus,
				CustomAlert:        hasCustomAlert,
			})
		}
	}

	return instances, nil
}

// getSystemChecksStatus retrieves the status of system checks for the instance (passed or failed).
func getSystemChecksStatus(ec2Client *ec2.EC2, instanceID string) string {
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	resp, err := ec2Client.DescribeInstanceStatus(params)
	if err != nil {
		log.Println("Error retrieving system checks status:", err)
		return "Unknown"
	}
	if len(resp.InstanceStatuses) == 0 {
		return "Unknown"
	}
	for _, status := range resp.InstanceStatuses {
		if aws.StringValue(status.InstanceStatus.Status) != "ok" {
			return "Failed"
		}
	}
	return "Passed"
}

// checkForCustomAlert checks if the instance has custom alerts.
func checkForCustomAlert(cloudWatchClient *cloudwatch.CloudWatch, instanceID string) (bool, error) {
	// Specify the filters to retrieve alarms associated with the given instance
	filters := []*cloudwatch.DimensionFilter{
		{
			Name:  aws.String("InstanceId"),
			Value: aws.String(instanceID),
		},
	}
	fmt.Println(filters)
	// Retrieve CloudWatch alarms using DescribeAlarms API
	resp, err := cloudWatchClient.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
		StateValue:      aws.String("ALARM"), // Optionally filter by alarm state
		AlarmNamePrefix: aws.String(instanceID),
	})
	if err != nil {
		return false, err
	}

	// If there are any alarms associated with the instance, return true
	if len(resp.MetricAlarms) > 0 {
		return true, nil
	}

	// Otherwise, return false
	return false, nil
}

func init() {
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxEc2InstanceStatusCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
