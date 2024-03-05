package EC2

import (
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/olekukonko/tablewriter"

	// "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	// "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type AlarmNotification struct {
	Timestamp   time.Time
	Alert       string
	Description string
}

var AwsxEc2AlarmandNotificationcmd = &cobra.Command{
	Use:   "alerts_and_notifications_panel",
	Short: "Retrieve recent alerts and notifications related to EC2 instance availability",
	Long:  `Command to retrieve recent alerts and notifications related to EC2 instance availability`,

	Run: func(cmd *cobra.Command, args []string) {
		authFlag, clientAuth, err := handleAuth(cmd)
		if err != nil {
			log.Println("Error during authentication:", err)
			return
		}

		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			notifications, err := GetAlertsAndNotificationsPanel(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting alerts and notifications:", err)
				return
			}

			if responseType == "frame" {
				fmt.Println(notifications)
			} else {
				printTable(notifications)
			}
		}
	},
}

func handleAuth(cmd *cobra.Command) (bool, *model.Auth, error) {
	authFlag, clientAuth, err := authenticate.AuthenticateCommand(cmd)
	if err != nil {
		return false, nil, err
	}
	return authFlag, clientAuth, nil
}

func GetAlertsAndNotificationsPanel(cmd *cobra.Command, clientAuth *model.Auth) ([]AlarmNotification, error) {
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return nil, err
		}
	} else {
		log.Println("Start time not provided. Please provide a start time.")
		return nil, errors.New("start time not provided")
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return nil, err
		}
	} else {
		endTime = time.Now()
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Retrieve CloudWatch alarms
	alarms, err := GetCloudWatchAlarms(clientAuth, &startTime, &endTime)
	if err != nil {
		log.Println("Error getting CloudWatch alarms:", err)
		return nil, err
	}

	// Convert CloudWatch alarms to AlarmNotification struct
	notifications := make([]AlarmNotification, len(alarms))
	for i, alarm := range alarms {
		notifications[i] = AlarmNotification{
			Timestamp:   *alarm.StateUpdatedTimestamp,
			Alert:       *alarm.StateReason,
			Description: *alarm.AlarmDescription,
		}
	}

	return notifications, nil
}

func GetCloudWatchAlarms(clientAuth *model.Auth, startTime, endTime *time.Time) ([]*cloudwatch.MetricAlarm, error) {
	svc := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	// Call DescribeAlarms to get all alarms
	resp, err := svc.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{})
	if err != nil {
		log.Println("Error describing alarms:", err)
		return nil, err
	}

	// Filter alarms based on their last updated time within the specified range
	filteredAlarms := make([]*cloudwatch.MetricAlarm, 0)
	for _, alarm := range resp.MetricAlarms {
		stateUpdatedTime := alarm.StateUpdatedTimestamp
		if stateUpdatedTime != nil && stateUpdatedTime.After(*startTime) && stateUpdatedTime.Before(*endTime) {
			filteredAlarms = append(filteredAlarms, alarm)
		}
	}

	return filteredAlarms, nil
}

func printTable(notifications []AlarmNotification) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Alert", "Description"})

	for _, notification := range notifications {
		table.Append([]string{
			notification.Timestamp.Format(time.RFC3339),
			notification.Alert,
			notification.Description,
		})
	}

	table.Render()
}

func init() {
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("ServiceName", "", "Service Name")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("clusterName", "", "cluster name")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxEc2InstanceRunningHourCmd.PersistentFlags().String("logGroupName", "", "log group name")
}
