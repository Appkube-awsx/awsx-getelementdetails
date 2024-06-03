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

type InstanceBackupStatus struct {
	SuccessfulBackups int `json:"SuccessfulBackups"`
	MissedBackups     int `json:"MissedBackups"`
}

var instanceBackupstatusPanelCmd = &cobra.Command{
	Use:   "backup_status",
	Short: "Gets the count of successful and missed backups",
	Long:  `Command to get the count of successful and missed backups using AWS EC2 snapshots`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running backup status command")
		authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Printf("Error during authentication: %v\n", err)
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, err := GetBackupStatus(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting rest API data: ", err)
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

func GetBackupStatus(cmd *cobra.Command, clientAuth *model.Auth, ec2Client *ec2.EC2) (string, error) {

	if ec2Client == nil {
		ec2Client = awsclient.GetClient(*clientAuth, awsclient.EC2_CLIENT).(*ec2.EC2)
	}
	allSnapshotsInput := &ec2.DescribeSnapshotsInput{}
	allSnapshotsResult, err := ec2Client.DescribeSnapshots(allSnapshotsInput)
	if err != nil {
		return "", fmt.Errorf("failed to describe all snapshots: %v", err)
	}

	// Count completed snapshots
	completedSnapshotsInput := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String("completed")},
			},
		},
	}
	completedSnapshotsResult, err := ec2Client.DescribeSnapshots(completedSnapshotsInput)
	if err != nil {
		return "", fmt.Errorf("failed to describe completed snapshots: %v", err)
	}

	totalSnapshotsCount := len(allSnapshotsResult.Snapshots)
	completedSnapshotsCount := len(completedSnapshotsResult.Snapshots)
	missedBackupsCount := totalSnapshotsCount - completedSnapshotsCount

	data := &InstanceBackupStatus{
		SuccessfulBackups: completedSnapshotsCount,
		MissedBackups:     missedBackupsCount,
	}

	//strData, nil := json.Marshal(data)
	jsonResp, _ := json.Marshal(data)
	return string(jsonResp), nil
}

func init() {
	comman_function.InitAwsCmdFlags(instanceBackupstatusPanelCmd)
}
