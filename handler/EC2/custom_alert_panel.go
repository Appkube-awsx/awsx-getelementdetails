package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	"log"
)

var AwsxEc2CustomAlertPanelCmd = &cobra.Command{
	Use:   "custom_alert_panel",
	Short: "get custom alerts for EC2 security group changes",
	Long:  `command to get custom alerts for EC2 security group changes`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from custom alert panel")

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
			cloudwatchMetric, err := GetEc2CustomAlertPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting custom alerts: ", err)
				return
			}
			fmt.Println(cloudwatchMetric)
		}
	},
}

func GetEc2CustomAlertPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) ([]*cloudwatchlogs.GetQueryResultsOutput, error) {
	logGroupName, _ := cmd.PersistentFlags().GetString("logGroupName")
	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}
	logGroupName, err = commanFunction.GetCmdbLogsData(cmd)
	if err != nil {
		return nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	results, err := commanFunction.FiltercloudWatchLogs(clientAuth, startTime, endTime, logGroupName, "")
	if err != nil {
		log.Println("Error in getting custom alert data: ", err)
		return nil, err
	}

	return results, nil
}

func init() {
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CustomAlertPanelCmd.PersistentFlags().String("secretKey", "", "aws secret key")
}
