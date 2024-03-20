package EC2

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/spf13/cobra"
)

var AwsxEc2InstanceHealthCheckCmd = &cobra.Command{

	Use:   "instance_health_check_panel",
	Short: "get instance health check metrics data",
	Long:  `command to get instance status metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")

		var authFlag, _, err = authenticate.AuthenticateCommand(cmd)

		if err != nil {
			log.Printf("Error during authentication: %v\n", err)

			err := cmd.Help()

			if err != nil {
				return
			}

			return
		}
		if authFlag {
			check, err := GetInstanceHealthCheck()
			if err != nil {
				log.Printf("Error getting instance status: %v", err)
			}
			fmt.Println(check)
		}
	},
}

type InstanceDummyData struct {
	InstanceID           string
	InstanceType         string
	AvailabilityZone     string
	InstanceStatus       string
	CpuUtilization       string
	DiskSpaceUtilization string
	SystemChecks         string
	InstanceChecks       string
	Alarm                string
	SystemCheck          string
	InstanceCheck        string
}

func GetInstanceHealthCheck() ([]InstanceDummyData, error) {
	instanceData := []InstanceDummyData{
		{
			InstanceID:           "i-1234567890abcdef0",
			InstanceType:         "t2.micro",
			AvailabilityZone:     "us-east-1a",
			InstanceStatus:       "running",
			CpuUtilization:       "10%",
			DiskSpaceUtilization: "50%",
			SystemChecks:         "ok",
			InstanceChecks:       "ok",
			Alarm:                "none",
			SystemCheck:          time.Now().Add(-1 * time.Minute).Format("06-01-02"), // Format as yy-mm-dd
			InstanceCheck:        time.Now().Add(-2 * time.Minute).Format("06-01-02"), // Format as yy-mm-dd
		},
		{
			InstanceID:           "i-0987654321fedcba0",
			InstanceType:         "t2.medium",
			AvailabilityZone:     "us-west-2b",
			InstanceStatus:       "stopped",
			CpuUtilization:       "0%",
			DiskSpaceUtilization: "75%",
			SystemChecks:         "ok",
			InstanceChecks:       "warning",
			Alarm:                "none",
			SystemCheck:          time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			InstanceCheck:        time.Now().Add(-4 * time.Minute).Format(time.RFC3339),
		},
		{
			InstanceID:           "i-0987654321fedcba0",
			InstanceType:         "t2.medium",
			AvailabilityZone:     "us-west-2b",
			InstanceStatus:       "stopped",
			CpuUtilization:       "50%",
			DiskSpaceUtilization: "85%",
			SystemChecks:         "ok",
			InstanceChecks:       "warning",
			Alarm:                "none",
			SystemCheck:          time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			InstanceCheck:        time.Now().Add(-4 * time.Minute).Format(time.RFC3339),
		},
		{
			InstanceID:           "i-0987654321fedcba0",
			InstanceType:         "t2.medium",
			AvailabilityZone:     "us-west-2b",
			InstanceStatus:       "stopped",
			CpuUtilization:       "40%",
			DiskSpaceUtilization: "75%",
			SystemChecks:         "ok",
			InstanceChecks:       "warning",
			Alarm:                "none",
			SystemCheck:          time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			InstanceCheck:        time.Now().Add(-4 * time.Minute).Format(time.RFC3339),
		},
		{
			InstanceID:           "i-0987654321fedcba0",
			InstanceType:         "t2.medium",
			AvailabilityZone:     "us-west-2b",
			InstanceStatus:       "stopped",
			CpuUtilization:       "30%",
			DiskSpaceUtilization: "45%",
			SystemChecks:         "ok",
			InstanceChecks:       "warning",
			Alarm:                "none",
			SystemCheck:          time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			InstanceCheck:        time.Now().Add(-4 * time.Minute).Format(time.RFC3339),
		},
	}
	return instanceData, nil
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
