package EC2

import (
	"encoding/json"
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

type Result struct {
	CurrentUsage float64 `json:"currentUsage"`
	AverageUsage float64 `json:"averageUsage"`
	MaxUsage     float64 `json:"maxUsage"`
}

var CpuUtilizationPanelCmd = &cobra.Command{
	Use:   "cpu_utilization_panel",
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
			instanceID := "i-5456-646g"
			metricName := "CPUUtilization"
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

			currentUsage, err := getMetricData(clientAuth, instanceID, metricName, namespace, startTime, endTime, "SampleCount")
			if err != nil {
				log.Fatal(err)
			}
			// Get average usage
			averageUsage, err := getMetricData(clientAuth, instanceID, metricName, namespace, startTime, endTime, "Average")
			if err != nil {
				log.Fatal(err)
			}

			// Get max usage
			maxUsage, err := getMetricData(clientAuth, instanceID, metricName, namespace, startTime, endTime, "Maximum")
			if err != nil {
				log.Fatal(err)
			}

			jsonOutput := map[string]float64{
				"CurrentUsage": *currentUsage.MetricDataResults[0].Values[0],
				"AverageUsage": *averageUsage.MetricDataResults[0].Values[0],
				"MaxUsage":     *maxUsage.MetricDataResults[0].Values[0],
			}

			jsonString, err := json.Marshal(jsonOutput)
			if err != nil {
				log.Fatal(err)
			}

			// Print the JSON string
			fmt.Println(string(jsonString))
		}

	},
}

func getMetricData(clientAuth *model.Auth, instanceID, metricName, namespace string, startTime, endTime *time.Time, statistic string) (*cloudwatch.GetMetricDataOutput, error) {
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

func Executed() {
	if err := CpuUtilizationPanelCmd.Execute(); err != nil {
		log.Println("error executing command: %v", err)
	}
}

func init() {
	CpuUtilizationPanelCmd.PersistentFlags().String("cloudElementId", "", "cloud element id")
	CpuUtilizationPanelCmd.PersistentFlags().String("cloudElementApiUrl", "", "cloud element api")
	CpuUtilizationPanelCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	CpuUtilizationPanelCmd.PersistentFlags().String("vaultToken", "", "vault token")
	CpuUtilizationPanelCmd.PersistentFlags().String("accountId", "", "aws account number")
	CpuUtilizationPanelCmd.PersistentFlags().String("zone", "", "aws region")
	CpuUtilizationPanelCmd.PersistentFlags().String("accessKey", "", "aws access key")
	CpuUtilizationPanelCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	CpuUtilizationPanelCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	CpuUtilizationPanelCmd.PersistentFlags().String("externalId", "", "aws external id")
	CpuUtilizationPanelCmd.PersistentFlags().String("elementType", "", "element type")
	CpuUtilizationPanelCmd.PersistentFlags().String("instanceID", "", "instance id")
	CpuUtilizationPanelCmd.PersistentFlags().String("query", "", "query")
	CpuUtilizationPanelCmd.PersistentFlags().String("timeRange", "", "timeRange")
}
