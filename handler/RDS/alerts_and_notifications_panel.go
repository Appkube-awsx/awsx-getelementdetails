package RDS

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"

	// "github.com/olekukonko/tablewriter"
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
				//printTable(notifications)
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
	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}

	// Retrieve CloudWatch alarms
	alarms, err := comman_function.GetCloudWatchAlarms(clientAuth, startTime, endTime)
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

// func printTable(notifications []AlarmNotification) {
// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetHeader([]string{"Timestamp", "Alert", "Description"})

// 	for _, notification := range notifications {
// 		table.Append([]string{
// 			notification.Timestamp.Format(time.RFC3339),
// 			notification.Alert,
// 			notification.Description,
// 		})
// 	}

// 	table.Render()
// }

func init() {
	comman_function.InitAwsCmdFlags(RdsAlarmandNotificationcmd)
}
