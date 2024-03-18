package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// ServiceStatus represents the status of a service.
type ServiceStatus struct {
	ServiceName  string `json:"serviceName"`
	HealthStatus string `json:"healthStatus"`
	ResponseTime string `json:"responseTime"`
	ErrorRate    string `json:"errorRate"`
	Availability string `json:"availability"`
	Throughput   string `json:"throughput"`
}

// AwsxEc2hostedServicesCmd represents the EC2 command.
var AwsxEc2hostedServicesCmd = &cobra.Command{
	Use:   "EC2",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		queryName, _ := cmd.Flags().GetString("query")
		elementType, _ := cmd.Flags().GetString("elementType")

		if queryName == "hosted_services_overview_panel" && (elementType == "EC2" || elementType == "AWS/EC2") {
			GetHostedServicesData(cmd)
		}
	},
}

func GetHostedServicesData(cmd *cobra.Command) {
	serviceStatus := []ServiceStatus{
		{
			ServiceName:  "WebServer",
			HealthStatus: "Healthy",
			ResponseTime: "100ms",
			ErrorRate:    "0.05%",
			Availability: "99.5%",
			Throughput:   "1000 req/s",
		},
		{
			ServiceName:  "Database",
			HealthStatus: "Degraded",
			ResponseTime: "115ms",
			ErrorRate:    "0.02%",
			Availability: "99.8%",
			Throughput:   "800 req/s",
		},
		{
			ServiceName:  "LoadBalancer",
			HealthStatus: "Degraded",
			ResponseTime: "156ms",
			ErrorRate:    "0.03%",
			Availability: "99.7%",
			Throughput:   "950 req/s",
		},
	}

	jsonData, err := json.Marshal(serviceStatus)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}

func init() {
	// Add flags for query and element type
	AwsxEc2hostedServicesCmd.Flags().String("query", "", "Query name")
	AwsxEc2hostedServicesCmd.Flags().String("elementType", "", "Element type")
	AwsxEc2hostedServicesCmd.Flags().String("serviceName", "", "Service name")
	AwsxEc2hostedServicesCmd.Flags().String("healthStatus", "", "Health status")
	AwsxEc2hostedServicesCmd.Flags().String("responseTime", "", "Response time")
	AwsxEc2hostedServicesCmd.Flags().String("errorRate", "", "Error rate")
	AwsxEc2hostedServicesCmd.Flags().String("availability", "", "Availability")
	AwsxEc2hostedServicesCmd.Flags().String("throughput", "", "Throughput")
}
