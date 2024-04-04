package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type MetricResults struct {
	TimeSeriesData map[string]string `json:"timeSeriesData"`
}

var AwsxRDSUptimeCmd = &cobra.Command{
	Use:   "rds_uptime_panel",
	Short: "get uptime metrics data for RDS",
	Long:  `command to get uptime metrics data for RDS`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running from child command")
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
			jsonResp, _, err := GetRDSUptimeData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting RDS uptime data: ", err)
				return
			}
			fmt.Println(jsonResp)
		}
	},
}

func GetRDSUptimeData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]string, error) {
	DBInstanceIdentifier := "postgresql"
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", nil, err
		}
		startTime = &parsedStartTime
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = &defaultStartTime
	}

	if endTimeStr != "" {
		parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			log.Printf("Error parsing end time: %v", err)
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	// Debug prints
	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	// Fetch raw data
	totalUptime, totalTime, err := GetDatabaseConnectionsMetricValues(clientAuth, startTime, endTime, DBInstanceIdentifier, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting database connections metric values: ", err)
		return "", nil, err
	}

	// Calculate uptime percentage
	uptimePercentage := (totalUptime / totalTime) * 100
	if uptimePercentage > 100 {
		uptimePercentage = 100
	}

	timeSeriesData := map[string]string{
		"uptimePercentage": fmt.Sprintf("%.2f%%", uptimePercentage),
	}

	log.Printf("Uptime Percentage: %f", uptimePercentage)

	jsonString, err := json.Marshal(MetricResults{TimeSeriesData: timeSeriesData})
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), timeSeriesData, nil
}

func GetDatabaseConnectionsMetricValues(clientAuth *model.Auth, startTime, endTime *time.Time, DBInstanceIdentifier string, cloudWatchClient *cloudwatch.CloudWatch) (float64, float64, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/RDS"),
		MetricName: aws.String("DatabaseConnections"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("DBInstanceIdentifier"),
				Value: aws.String(DBInstanceIdentifier),
			},
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     aws.Int64(300), // 5-minute intervals
		Statistics: []*string{aws.String("Sum")},
		Unit:       aws.String("Count"),
	}

	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, 0, err
	}

	totalUptime := 0.0
	totalTime := endTime.Sub(*startTime).Minutes()

	for _, dp := range result.Datapoints {
		totalUptime += *dp.Sum
	}

	return totalUptime, totalTime, nil
}

func init() {
	AwsxRDSUptimeCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxRDSUptimeCmd.PersistentFlags().String("endTime", "", "end time")
}
