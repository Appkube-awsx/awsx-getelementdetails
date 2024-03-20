package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// ErrorEvent represents an error event.
type ErrorEvent struct {
	EventID          string `json:"event_id"`
	Timestamp        string `json:"timestamp"`
	ErrorCode        string `json:"error_code"`
	Severity         string `json:"severity"`
	Description      string `json:"description"`
	SourceComponent  string `json:"source_component"`
	ActionTaken      string `json:"action_taken"`
	ResolutionStatus string `json:"resolution_status"`
	AdditionalNotes  string `json:"additional_notes"`
}

// ListErrorsCmd represents the command to list error events.
var ListErrorsCmd = &cobra.Command{
	Use:   "listErrors",
	Short: "List error events",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ListErrorEvents()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println(string(data))
	},
}

// ListErrorEvents retrieves a list of error events.
func ListErrorEvents() ([]byte, error) {
	// Dummy error events (replace with actual data retrieval)
	errorEvents := []ErrorEvent{
		{
			EventID:          "EVT-001",
			Timestamp:        "2023-07-15 10:45 AM",
			ErrorCode:        "ERR-101",
			Severity:         "critical",
			Description:      "Database Connection timeout",
			SourceComponent:  "Database server",
			ActionTaken:      "investigating",
			ResolutionStatus: "pending",
			AdditionalNotes:  "AWS support ticket",
		},
		{
			EventID:          "EVT-002",
			Timestamp:        "2023-07-15 10:45 AM",
			ErrorCode:        "ERR-204",
			Severity:         "Major",
			Description:      "Application server crash",
			SourceComponent:  "App Server",
			ActionTaken:      "restarted",
			ResolutionStatus: "resolved",
			AdditionalNotes:  "Updated application to the latest version",
		},
		{
			EventID:          "EVT-003",
			Timestamp:        "2023-07-15 10:45 AM",
			ErrorCode:        "ERR-302",
			Severity:         "Minor",
			Description:      "high cpu usage alert",
			SourceComponent:  "Monitoring agent",
			ActionTaken:      "analyzing,scaling",
			ResolutionStatus: "ongoing investigation",
			AdditionalNotes:  "deployed additional monitoring tools",
		},
	}

	// Convert error events to JSON
	data, err := json.MarshalIndent(errorEvents, "", "  ")
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(data))

	return data, nil
}

func init() {
	// Add flags for query and element type
	ListErrorsCmd.Flags().String("query", "", "Query name")
	ListErrorsCmd.Flags().String("elementType", "", "Element type")
}
