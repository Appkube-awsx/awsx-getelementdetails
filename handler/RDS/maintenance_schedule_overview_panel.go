package RDS

import (
	// "encoding/json"
	// "fmt"
	// "log"
	"github.com/Appkube-awsx/awsx-getelementdetails/comman-function"

	"github.com/spf13/cobra"
)

type ScheduleOverview struct {
	MaintenanceType string `json:"MAINTENANCE TYPE"`
	Description     string `json:"DESCRIPTION"`
	StartTime       string `json:"START TIME"`
	EndTime         string `json:"END TIME"`
}

var ListScheduleOverviewCmd = &cobra.Command{
	Use:   "ListScheduleOverview",
	Short: "List schedule overview",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ListScheduleOverview()
		if err != nil {
			return
		}
	},
}

func ListScheduleOverview() ([]ScheduleOverview, error) {
	scheduleoverview := []ScheduleOverview{
		{
			MaintenanceType: "Patch Update",
			Description:     "Applying security patches and update",
			StartTime:       "2023-09-03 01:00 AM",
			EndTime:         "2023-09-05 10:00 AM",
		},
		{
			MaintenanceType: "Database Backup",
			Description:     "Full backup of the database for disaster recovery",
			StartTime:       "2023-09-05 01:00 AM",
			EndTime:         "2023-09-10 10:00 AM",
		},
		{
			MaintenanceType: "Network Changes",
			Description:     "Network configuration changes for optimization",
			StartTime:       "2023-09-10 01:00 AM",
			EndTime:         "2023-09-15 10:00 AM",
		},
	}

	// // Convert error events to JSON and print them
	// jsonData, err := json.MarshalIndent(scheduleoverview, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshalling JSON: %v", err)
	// }
	// fmt.Println(string(jsonData))
	return scheduleoverview, nil
}

func init() {
	comman_function.InitAwsCmdFlags(ListScheduleOverviewCmd)
}
