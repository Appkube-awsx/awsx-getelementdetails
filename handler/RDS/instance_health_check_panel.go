package RDS

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/spf13/cobra"
)

var AwsxDBInstanceHealthCheckCmd = &cobra.Command{

	Use:   "dbinstance_health_check_panel",
	Short: "get dbinstance health check metrics data",
	Long:  `command to get dbinstance status metrics data`,

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
			check, err := GetDBInstanceHealthCheck()
			if err != nil {
				log.Printf("Error getting instance status: %v", err)
			}
			fmt.Println(check)
		}
	},
}

type  InstanceHealthCheckData struct {
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

func GetDBInstanceHealthCheck() ([]InstanceHealthCheckData, error) {
	instanceData := []InstanceHealthCheckData{
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
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("query", "", "query")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxDBInstanceHealthCheckCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
