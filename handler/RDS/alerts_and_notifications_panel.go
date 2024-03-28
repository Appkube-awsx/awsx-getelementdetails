package RDS

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/olekukonko/tablewriter"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type AlarmNotification struct {
	Timestamp   time.Time
	Alert       string
	Description string
}

var RdsAlarmandNotificationcmd = &cobra.Command{
	Use:   "rds_alerts_and_notifications_panel",
	Short: "Retrieve recent alerts and notifications related to RDS instance availability",
	Long:  `Command to retrieve recent alerts and notifications related to RDS instance availability`,

	Run: func(cmd *cobra.Command, args []string) {
		authFlag, clientAuth, err := handleAuth(cmd)
		if err != nil {
			log.Println("Error during authentication:", err)
			return
		}

		if authFlag {
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			notifications, err := GetAlertsAndNotificationsPanell(cmd, clientAuth)
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

func GetAlertsAndNotificationsPanell(cmd *cobra.Command, clientAuth *model.Auth) ([]AlarmNotification, error) {
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
	//RdsAlarmandNotificationcmd.PersistentFlags().String("instanceId", "", "RDS instance ID")
	//RdsAlarmandNotificationcmd.PersistentFlags().String("startTime", "", "Start time for filtering alerts")
	//RdsAlarmandNotificationcmd.PersistentFlags().String("endTime", "", "End time for filtering alerts")
	//RdsAlarmandNotificationcmd.PersistentFlags().String("responseType", "", "Response type. json/frame")
	RdsAlarmandNotificationcmd.PersistentFlags().String("rootvolumeId", "", "root volume id")
	RdsAlarmandNotificationcmd.PersistentFlags().String("ebsvolume1Id", "", "ebs volume 1 id")
	RdsAlarmandNotificationcmd.PersistentFlags().String("ebsvolume2Id", "", "ebs volume 2 id")
	RdsAlarmandNotificationcmd.PersistentFlags().String("elementId", "", "element id")
	RdsAlarmandNotificationcmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	RdsAlarmandNotificationcmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	RdsAlarmandNotificationcmd.PersistentFlags().String("vaultToken", "", "vault token")
	RdsAlarmandNotificationcmd.PersistentFlags().String("accountId", "", "aws account number")
	RdsAlarmandNotificationcmd.PersistentFlags().String("zone", "", "aws region")
	RdsAlarmandNotificationcmd.PersistentFlags().String("accessKey", "", "aws access key")
	RdsAlarmandNotificationcmd.PersistentFlags().String("secretKey", "", "aws secret key")
	RdsAlarmandNotificationcmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	RdsAlarmandNotificationcmd.PersistentFlags().String("externalId", "", "aws external id")
	RdsAlarmandNotificationcmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	RdsAlarmandNotificationcmd.PersistentFlags().String("ServiceName", "", "Service Name")
	RdsAlarmandNotificationcmd.PersistentFlags().String("elementType", "", "element type")
	RdsAlarmandNotificationcmd.PersistentFlags().String("instanceId", "", "instance id")
	RdsAlarmandNotificationcmd.PersistentFlags().String("clusterName", "", "cluster name")
	RdsAlarmandNotificationcmd.PersistentFlags().String("query", "", "query")
	RdsAlarmandNotificationcmd.PersistentFlags().String("startTime", "", "start time")
	RdsAlarmandNotificationcmd.PersistentFlags().String("endTime", "", "endcl time")
	RdsAlarmandNotificationcmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	RdsAlarmandNotificationcmd.PersistentFlags().String("logGroupName", "", "log group name")
}
