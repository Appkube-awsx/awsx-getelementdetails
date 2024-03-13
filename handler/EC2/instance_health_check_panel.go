package EC2

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var AwsxEc2InstanceHealthCheckCmd = &cobra.Command{

	Use:   "instance_health_check_panel",
	Short: "get instance health check metrics data",
	Long:  `command to get instance status metrics data`,

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
			GetInstanceHealthCheck(clientAuth)
		}
	},
}

// GetInstanceStatus retrieves EC2 instance information including instance ID, instance type,
// availability zone, system check status, and custom alerts.
func GetInstanceHealthCheck(clientauth *model.Auth) error {
	// Initialize EC2 client
	ec2Client := awsclient.GetClient(*clientauth, awsclient.EC2_CLIENT).(*ec2.EC2)

	// Initialize CloudWatch client
	cloudWatchClient := awsclient.GetClient(*clientauth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	log.Println("Getting AWS EC2 instance list")

	// Retrieve instance status
	resp, err := ec2Client.DescribeInstances(nil)
	if err != nil {
		return err
	}

	// Print header
	fmt.Printf("%-20s %-15s %-15s %-15s %-20s %-15s %-5s %-25s %-25s\n", "Instance ID", "Instance Type", "Availability Zone", "State", "System Checks Status", "Instance Checks Status", "Alarm", "System Check Time", "Instance Check Time")

	// Print instance information
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceID := aws.StringValue(instance.InstanceId)
			instanceType := aws.StringValue(instance.InstanceType)
			availabilityZone := aws.StringValue(instance.Placement.AvailabilityZone)
			state := aws.StringValue(instance.State.Name)
			systemChecksStatus := getSystemChecksStatus(ec2Client, instanceID)
			instanceChecksStatus := getInstanceChecksStatus(ec2Client, instanceID)
			alarmStatus, systemCheckTime, instanceCheckTime := getAlarmAndCheckStatus(cloudWatchClient, instanceID)

			// Print instance details
			fmt.Printf("%-20s %-15s %-15s %-15s %-20s %-15s %-5t %-25s %-25s\n",
				instanceID, instanceType, availabilityZone, state, systemChecksStatus,
				instanceChecksStatus, alarmStatus, systemCheckTime, instanceCheckTime)
		}
	}

	return nil
}

// getInstanceChecksStatus retrieves the status of instance checks for the instance (passed or failed).
func getInstanceChecksStatus(ec2Client *ec2.EC2, instanceID string) string {
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	resp, err := ec2Client.DescribeInstanceStatus(params)
	if err != nil {
		log.Println("Error retrieving instance checks status:", err)
		return "Unknown"
	}
	if len(resp.InstanceStatuses) == 0 {
		return "Unknown"
	}
	for _, status := range resp.InstanceStatuses {
		if aws.StringValue(status.InstanceState.Name) != "running" {
			return "Failed"
		}
	}
	return "Passed"
}

// getAlarmAndCheckStatus retrieves the status of alarms and the time of the last system and instance checks.
func getAlarmAndCheckStatus(cloudWatchClient *cloudwatch.CloudWatch, instanceID string) (bool, string, string) {
	// Retrieve CloudWatch alarms using DescribeAlarms API
	resp, err := cloudWatchClient.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
		StateValue:      aws.String("ALARM"), // Optionally filter by alarm state
		AlarmNamePrefix: aws.String(instanceID),
	})
	if err != nil {
		log.Printf("Error retrieving alarms for instance %s: %v", instanceID, err)
		return false, "Unknown", "Unknown"
	}

	// If there are any alarms associated with the instance, return true and the time of the last system and instance checks
	if len(resp.MetricAlarms) > 0 {
		systemCheckTime, instanceCheckTime := getLastCheckTimes(cloudWatchClient, instanceID)
		return true, systemCheckTime, instanceCheckTime
	}

	// Otherwise, return false and "Unknown" for check times
	return false, "Unknown", "Unknown"
}

// getLastCheckTimes retrieves the time of the last system and instance checks.
func getLastCheckTimes(cloudWatchClient *cloudwatch.CloudWatch, instanceID string) (string, string) {
	// Retrieve CloudWatch metric data for system and instance checks
	systemCheckData, err := cloudWatchClient.GetMetricData(&cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("system-checks"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String("StatusCheckFailed_System"),
					},
					Period: aws.Int64(300), // 5-minute period for system checks
					Stat:   aws.String("Maximum"),
				},
			},
		},
		StartTime: aws.Time(time.Now().Add(-time.Hour)), // Start time: 1 hour ago
		EndTime:   aws.Time(time.Now()),                 // End time: Now
		ScanBy:    aws.String("TimestampDescending"),
	})
	if err != nil {
		log.Printf("Error retrieving system check data for instance %s: %v", instanceID, err)
		return "Unknown", "Unknown"
	}

	instanceCheckData, err := cloudWatchClient.GetMetricData(&cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("instance-checks"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String("StatusCheckFailed"),
					},
					Period: aws.Int64(300), // 5-minute period for instance checks
					Stat:   aws.String("Maximum"),
				},
			},
		},
		StartTime: aws.Time(time.Now().Add(-time.Hour)), // Start time: 1 hour ago
		EndTime:   aws.Time(time.Now()),                 // End time: Now
		ScanBy:    aws.String("TimestampDescending"),
	})
	if err != nil {
		log.Printf("Error retrieving instance check data for instance %s: %v", instanceID, err)
		return "Unknown", "Unknown"
	}

	// Extract the timestamps of the last system and instance checks
	var systemCheckTime, instanceCheckTime string
	if len(systemCheckData.MetricDataResults) > 0 && len(systemCheckData.MetricDataResults[0].Timestamps) > 0 {
		systemCheckTime = systemCheckData.MetricDataResults[0].Timestamps[0].Format(time.RFC3339)
	}
	if len(instanceCheckData.MetricDataResults) > 0 && len(instanceCheckData.MetricDataResults[0].Timestamps) > 0 {
		instanceCheckTime = instanceCheckData.MetricDataResults[0].Timestamps[0].Format(time.RFC3339)
	}

	return systemCheckTime, instanceCheckTime
}

func init() {
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxEc2InstanceHealthCheckCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
