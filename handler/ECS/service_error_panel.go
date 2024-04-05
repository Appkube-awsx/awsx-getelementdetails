package ECS

import (
	"github.com/spf13/cobra"
)

type ServiceError struct {
	Timestamp         string `json:"TIMESTAMP"`
	ServiceName       string `json:"SERVICE NAME"`
	TaskID            string `json:"TASK ID"`
	ErrorType         string `json:"ERROR TYPE"`
	ErrorDescription  string `json:"ERROR DESCRIPTION"`
	ResolutionTime    string `json:"RESOLUTION TIMESTAMP"`
	ImpactLevel       string `json:"IMPACT LEVEL"`
	ResolutionDetails string `json:"RESOLUTION DETAILS"`
	Status            string `json:"STATUS"`
}

var AwsxEcsServiceErrorCmd = &cobra.Command{
	Use:   "AwsxEcsServiceError",
	Short: "List AWS ECS service errors",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ListServiceErrors()
		if err != nil {
			return
		}
	},
}

func ListServiceErrors() ([]ServiceError, error) {
	serviceErrors := []ServiceError{
		{
			Timestamp:         "2023-09-03 08:15:00",
			ServiceName:       "ApiUserService",
			TaskID:            "task-123",
			ErrorType:         "4xx",
			ErrorDescription:  "Client Error",
			ResolutionTime:    "2023-09-03 08:30:00",
			ImpactLevel:       "High",
			ResolutionDetails: "Increased timeout settings",
			Status:            "investigating",
		},
		{
			Timestamp:         "2023-09-04 10:30:00",
			ServiceName:       "ImageProcessingService",
			TaskID:            "task-456",
			ErrorType:         "5xx",
			ErrorDescription:  "Server Error",
			ResolutionTime:    "2023-09-04 11:00:00",
			ImpactLevel:       "Medium",
			ResolutionDetails: "Increased memory allocation",
			Status:            "resolved",
		},
		{
			Timestamp:         "2023-09-05 12:45:00",
			ServiceName:       "NotificationService",
			TaskID:            "task-789",
			ErrorType:         "2xx",
			ErrorDescription:  "Success",
			ResolutionTime:    "2023-09-05 13:00:00",
			ImpactLevel:       "Low",
			ResolutionDetails: "Corrected configuration parameters",
			Status:            "ongoing",
		},
		{
			Timestamp:         "2023-09-06 14:00:00",
			ServiceName:       "BackendService",
			TaskID:            "task-101",
			ErrorType:         "3xx",
			ErrorDescription:  "Redirection",
			ResolutionTime:    "2023-09-06 14:30:00",
			ImpactLevel:       "High",
			ResolutionDetails: "Restarted database service",
			Status:            "resolved",
		},
		{
			Timestamp:         "2023-09-07 16:00:00",
			ServiceName:       "NotificationService",
			TaskID:            "task-102",
			ErrorType:         "4xx",
			ErrorDescription:  "Client Error",
			ResolutionTime:    "2023-09-07 16:30:00",
			ImpactLevel:       "Medium",
			ResolutionDetails: "Updated API endpoint",
			Status:            "resolved",
		},
		{
			Timestamp:         "2023-09-08 18:00:00",
			ServiceName:       "ApiUserService",
			TaskID:            "task-103",
			ErrorType:         "5xx",
			ErrorDescription:  "Server Error",
			ResolutionTime:    "2023-09-08 18:30:00",
			ImpactLevel:       "High",
			ResolutionDetails: "Increased server capacity",
			Status:            "ongoing",
		},
	}

	return serviceErrors, nil
}

	
func init() {
	AwsxEcsServiceErrorCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEcsServiceErrorCmd.PersistentFlags().String("endTime", "", "end time")
}

