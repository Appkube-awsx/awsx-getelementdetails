package EC2

import (
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

type InstanceInfo struct {
	InstanceID         string
	InstanceType       string
	AvailabilityZone   string
	State              string
	SystemChecksStatus string
	CustomAlert        bool
	HealthPercentage   float64
}

func GetInstanceStatus(cmd *cobra.Command, clientauth *model.Auth) (InstanceInfo, error) {
    cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	// elementType, _ := cmd.PersistentFlags().GetString("elementType")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, _ := cmdb.GetCloudElementData(apiUrl, elementId)
		// if err != nil {
		// 	return ,err
		// }
		instanceId = cmdbData.InstanceId
	}
    instanceID := "i-078bafb47ad7de492"// Initialize EC2 client
	ec2Client := awsclient.GetClient(*clientauth, awsclient.EC2_CLIENT).(*ec2.EC2)

	// Initialize CloudWatch client
	cloudWatchClient := awsclient.GetClient(*clientauth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	log.Printf("Getting AWS EC2 instance status for instance ID: %s\n", instanceId)

	// Retrieve instance information
	resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceId)},
	})
	if err != nil {
		return InstanceInfo{}, err
	}

	// Check if instance exists
	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return InstanceInfo{}, fmt.Errorf("instance with ID %s not found", instanceId)
	}

	instance := resp.Reservations[0].Instances[0]

	// Get system checks status
	systemChecksStatus := getSystemChecksStatus(ec2Client, instanceId)

	// Check for custom alert
	hasCustomAlert, err := checkForCustomAlert(cloudWatchClient, instanceId)
	if err != nil {
		return InstanceInfo{}, err
	}

	// Calculate health percentage
	passedCount, failedCount := 0, 0
	switch systemChecksStatus {
	case "Passed":
		passedCount = 1
	case "Failed":
		failedCount = 1
	}
	totalInstances := passedCount + failedCount
	var healthPercentage float64
	if totalInstances > 0 {
		healthPercentage = float64(passedCount) / float64(totalInstances) * 100
	}

	instanceInfo := InstanceInfo{
		InstanceID:         instanceID,
		InstanceType:       aws.StringValue(instance.InstanceType),
		AvailabilityZone:   aws.StringValue(instance.Placement.AvailabilityZone),
		State:              aws.StringValue(instance.State.Name),
		SystemChecksStatus: systemChecksStatus,
		CustomAlert:        hasCustomAlert,
		HealthPercentage:   healthPercentage,
	}

	return instanceInfo, nil
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
		return "Stopped/Terminated"
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
	// Retrieve CloudWatch alarms using DescribeAlarms API
	resp, err := cloudWatchClient.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
		StateValue:      aws.String("ALARM"), // Optionally filter by alarm state
		AlarmNamePrefix: aws.String(instanceID),
	})
	if err != nil {
		return false, err
	}

	// If there are any alarms associated with the instance, return true
	return len(resp.MetricAlarms) > 0, nil
}

var AwsxEc2InstanceStatusCmd = &cobra.Command{
	Use:   "instance_status_panel",
	Short: "get instance status metrics data",
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
			// Example instance ID


			// Call GetInstanceStatus
			instanceStatus, err := GetInstanceStatus(cmd, clientAuth)
			if err != nil {
				log.Fatalf("Error getting instance status: %v", err)
			}

			// Print instance information
			fmt.Printf("Instance ID: %s, Instance Type: %s, Availability Zone: %s, State: %s, System Checks Status: %s, Custom Alert: %t, Health Percentage: %.2f%%\n",
				instanceStatus.InstanceID, instanceStatus.InstanceType, instanceStatus.AvailabilityZone, instanceStatus.State, instanceStatus.SystemChecksStatus, instanceStatus.CustomAlert, instanceStatus.HealthPercentage)
		}
	},
}





// package EC2

// import (
// 	"fmt"
// 	"log"

// 	"github.com/Appkube-awsx/awsx-common/authenticate"
// 	"github.com/Appkube-awsx/awsx-common/awsclient"
// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// 	"github.com/spf13/cobra"
// )

// type InstanceInfo struct {
// 	InstanceID         string
// 	InstanceType       string
// 	AvailabilityZone   string
// 	State              string
// 	SystemChecksStatus string
// 	CustomAlert        bool
// 	HealthPercentage   float64
// }

// var instanceStatusData []InstanceInfo

// func GetInstanceStatus(cmd *cobra.Command, clientauth *model.Auth) ([]InstanceInfo, error) {
// 	// Initialize EC2 client
// 	ec2Client := awsclient.GetClient(*clientauth, awsclient.EC2_CLIENT).(*ec2.EC2)

// 	// Initialize CloudWatch client
// 	cloudWatchClient := awsclient.GetClient(*clientauth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

// 	log.Println("Getting AWS EC2 instance list")

// 	// Retrieve instance status
// 	resp, err := ec2Client.DescribeInstances(nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Populate instance information slice
// 	for _, reservation := range resp.Reservations {
// 		for _, instance := range reservation.Instances {
// 			instanceID := aws.StringValue(instance.InstanceId)
// 			instanceType := aws.StringValue(instance.InstanceType)
// 			availabilityZone := aws.StringValue(instance.Placement.AvailabilityZone)
// 			state := aws.StringValue(instance.State.Name)
// 			systemChecksStatus := getSystemChecksStatus(ec2Client, instanceID)
// 			hasCustomAlert, err := checkForCustomAlert(cloudWatchClient, instanceID)
// 			if err != nil {
// 				// log.Printf("Error checking custom alert for instance %s: %v", instanceID, err)
// 				continue // Skip to the next instance
// 			}

// 			// Calculate health percentage
// 			passedCount, failedCount := 0, 0
// 			switch systemChecksStatus {
// 			case "Passed":
// 				passedCount = 1
// 			case "Failed":
// 				failedCount = 1
// 			}
// 			totalInstances := passedCount + failedCount
// 			var healthPercentage float64
// 			if totalInstances > 0 {
// 				healthPercentage = float64(passedCount) / float64(totalInstances) * 100
// 			}

// 			// Append instance information to the slice
// 			instanceStatusData = append(instanceStatusData, InstanceInfo{
// 				InstanceID:         instanceID,
// 				InstanceType:       instanceType,
// 				AvailabilityZone:   availabilityZone,
// 				State:              state,
// 				SystemChecksStatus: systemChecksStatus,
// 				CustomAlert:        hasCustomAlert,
// 				HealthPercentage:   healthPercentage,
// 			})
// 		}
// 	}

// 	return instanceStatusData, nil
// }

// // getSystemChecksStatus retrieves the status of system checks for the instance (passed or failed).
// func getSystemChecksStatus(ec2Client *ec2.EC2, instanceID string) string {
// 	params := &ec2.DescribeInstanceStatusInput{
// 		InstanceIds: []*string{aws.String(instanceID)},
// 	}
// 	resp, err := ec2Client.DescribeInstanceStatus(params)
// 	if err != nil {
// 		log.Println("Error retrieving system checks status:", err)
// 		return "Unknown"
// 	}
// 	if len(resp.InstanceStatuses) == 0 {
// 		return "Unknown"
// 	}
// 	for _, status := range resp.InstanceStatuses {
// 		if aws.StringValue(status.InstanceStatus.Status) != "ok" {
// 			return "Failed"
// 		}
// 	}
// 	return "Passed"
// }

// // checkForCustomAlert checks if the instance has custom alerts.
// func checkForCustomAlert(cloudWatchClient *cloudwatch.CloudWatch, instanceID string) (bool, error) {
// 	// Retrieve CloudWatch alarms using DescribeAlarms API
// 	resp, err := cloudWatchClient.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
// 		StateValue:      aws.String("ALARM"), // Optionally filter by alarm state
// 		AlarmNamePrefix: aws.String(instanceID),
// 	})
// 	if err != nil {
// 		return false, err
// 	}

// 	// If there are any alarms associated with the instance, return true
// 	return len(resp.MetricAlarms) > 0, nil
// }

// var AwsxEc2InstanceStatusCmd = &cobra.Command{
// 	Use:   "instance_status_panel",
// 	Short: "get instance status metrics data",
// 	Long:  `command to get instance status metrics data`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("running from child command")

// 		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
// 		if err != nil {
// 			log.Printf("Error during authentication: %v\n", err)
// 			err := cmd.Help()
// 			if err != nil {
// 				return
// 			}
// 			return
// 		}

// 		if authFlag {
// 			// Call GetInstanceStatus
// 			instanceStatusData, err := GetInstanceStatus(cmd, clientAuth)
// 			if err != nil {
// 				log.Fatalf("Error getting instance status: %v", err)
// 			}

// 			// Print or utilize the instance information
// 			for _, info := range instanceStatusData {
// 				fmt.Printf("Instance ID: %s, Instance Type: %s, Availability Zone: %s, State: %s, System Checks Status: %s, Custom Alert: %t, Health Percentage: %.2f%%\n",
// 					info.InstanceID, info.InstanceType, info.AvailabilityZone, info.State, info.SystemChecksStatus, info.CustomAlert, info.HealthPercentage)
// 			}
// 		}
// 	},
// }

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
