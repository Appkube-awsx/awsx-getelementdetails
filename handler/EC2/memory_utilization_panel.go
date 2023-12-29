package EC2

import (
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var MemoryUtilizationPanelCmd = &cobra.Command{
	Use:   "memory_utilization_panel",
	Short: "getCpuUtilizationPanel command gets cloudwatch metrics data",
	Long:  `getCpuUtilizationPanel command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		var authFlag, clientAuth, err = authenticate.AuthenticateCommand(cmd)
		if err != nil {
			log.Println("Error during authentication: %v", err)
			cmd.Help()
			return
		}
		if authFlag {
			instanceID := "i-05e4e6757f13da657"
			metricName := "MemoryUtilization"
			namespace := "AWS/EC2"

			startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
			endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

			var startTime, endTime *time.Time

			// Parse start time if provided
			if startTimeStr != "" {
				parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
				if err != nil {
					log.Printf("Error parsing start time: %v", err)
					cmd.Help()
					return
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
					cmd.Help()
					return
				}
				endTime = &parsedEndTime
			} else {
				defaultEndTime := time.Now()
				endTime = &defaultEndTime
			}

			currentUsage, err := getMemoryMetricData(clientAuth, instanceID, metricName, namespace, startTime, endTime, "SampleCount")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Current Memory Usage: %v\n", currentUsage)
			// Get average usage
			averageUsage, err := getMemoryMetricData(clientAuth, instanceID, metricName, namespace, startTime, endTime, "Average")
			if err != nil {
				log.Fatal(err)
			}

			// Get max usage
			maxUsage, err := getMemoryMetricData(clientAuth, instanceID, metricName, namespace, startTime, endTime, "Maximum")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Current Memory Usage: %v\n", currentUsage)
			fmt.Printf("Average Memory Usage: %v\n", averageUsage)
			fmt.Printf("Max Memory Usage: %v\n", maxUsage)
		}

	},
}

func processMetricData(currentUsage, averageUsage, maxUsage *cloudwatch.GetMetricDataOutput) {
	fmt.Printf("Current Memory Usage: %v\n", extractMetricValue(currentUsage))
	fmt.Printf("Average Memory Usage: %v\n", extractMetricValue(averageUsage))
	fmt.Printf("Max Memory Usage: %v\n", extractMetricValue(maxUsage))
}

func extractMetricValue(metricData *cloudwatch.GetMetricDataOutput) float64 {
	if len(metricData.MetricDataResults) > 0 && len(metricData.MetricDataResults[0].Values) > 0 {
		return *metricData.MetricDataResults[0].Values[0]
	}
	return 0.0
}
func getMemoryMetricData(clientAuth *model.Auth, instanceID, metricName, namespace string, startTime, endTime *time.Time, statistic string) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String("i-05e4e6757f13da657"),
							},
						},
						MetricName: aws.String(metricName),
						Namespace:  aws.String(namespace),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
					Unit:   aws.String("Bytes"),
				},
			},
		},
	}
	cloudWatchClient := awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Execute() {
	if err := MemoryUtilizationPanelCmd.Execute(); err != nil {
		log.Println("error executing command: %v", err)
	}
}

func init() {
	MemoryUtilizationPanelCmd.PersistentFlags().String("cloudElementId", "", "cloud element id")
	MemoryUtilizationPanelCmd.PersistentFlags().String("cloudElementApiUrl", "", "cloud element api")
	MemoryUtilizationPanelCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	MemoryUtilizationPanelCmd.PersistentFlags().String("vaultToken", "", "vault token")
	MemoryUtilizationPanelCmd.PersistentFlags().String("accountId", "", "aws account number")
	MemoryUtilizationPanelCmd.PersistentFlags().String("zone", "", "aws region")
	MemoryUtilizationPanelCmd.PersistentFlags().String("accessKey", "", "aws access key")
	MemoryUtilizationPanelCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	MemoryUtilizationPanelCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	MemoryUtilizationPanelCmd.PersistentFlags().String("externalId", "", "aws external id")
	MemoryUtilizationPanelCmd.PersistentFlags().String("elementType", "", "element type")
	MemoryUtilizationPanelCmd.PersistentFlags().String("instanceID", "", "instance id")
	MemoryUtilizationPanelCmd.PersistentFlags().String("query", "", "query")
	MemoryUtilizationPanelCmd.PersistentFlags().String("timeRange", "", "timeRange")
}
