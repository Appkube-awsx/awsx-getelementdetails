package EKS

import (
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/commanFunction"
	"github.com/Appkube-awsx/awsx-getelementdetails/global-function/metricData"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// type nodeFailureResult struct {
// 	RawData []struct {
// 		Timestamp time.Time
// 		Value     float64
// 	} `json:"Node Failures"`
// }

var AwsxEKSNodeFailureCmd = &cobra.Command{
	Use:   "node_failure_panel",
	Short: "Get node failure metrics data",
	Long:  `Command to get node failure metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running from child command")
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
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err := GetNodeFailureData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting node failure data: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetNodeFailureData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	fmt.Println(elementType)

	startTime, endTime, err := commanFunction.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	instanceId, err = commanFunction.GetCmdbData(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error getting instance ID: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data
	rawData, err := metricData.GetMetricClusterData(clientAuth, instanceId, "ContainerInsights", "cluster_failed_node_count", startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Node Failures"] = rawData

	// Process the raw data if needed
	// result := processNodeFailureRawData(rawData)

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	log.Println("Error in marshalling JSON in string: ", err)
	// 	return "", nil, err
	// }

	return "", cloudwatchMetricData, nil
}

// func processNodeFailureRawData(result *cloudwatch.GetMetricDataOutput) nodeFailureResult {
// 	var rawData nodeFailureResult

// 	// Initialize map to store data by date
// 	dateMap := make(map[time.Time]float64)

// 	// Iterate through timestamps and values
// 	for i := 0; i < len(result.MetricDataResults[0].Timestamps); i++ {
// 		timestamp := result.MetricDataResults[0].Timestamps[i]

// 		// Truncate timestamp to start of day
// 		date := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())

// 		// Add value to corresponding date
// 		dateMap[date] += *result.MetricDataResults[0].Values[i]
// 	}

// 	// Convert map to array of struct
// 	for date, value := range dateMap {
// 		rawData.RawData = append(rawData.RawData, struct {
// 			Timestamp time.Time
// 			Value     float64
// 		}{
// 			Timestamp: date,
// 			Value:     value,
// 		})
// 	}

// 	return rawData
// }

// Function to calculate the number of days ago based on the day string
// Function to calculate the number of days ago based on the day string
func daysAgo(day string) int {
	today := time.Now().Weekday()
	targetDay, _ := time.Parse("Monday", day) // Parsing day string to time.Time
	targetDayOfWeek := int(targetDay.Weekday())
	daysAgo := today - time.Weekday(targetDayOfWeek)
	if daysAgo < 0 {
		daysAgo += 7 // Wrap around to previous week
	}
	return int(daysAgo)
}

func init() {
	AwsxEKSNodeFailureCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNodeFailureCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
