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

var AwsxEc2CpuUtilizationCmd = &cobra.Command{
	Use:   "cpu_utilization_panel",
	Short: "get cpu utilization metrics data",
	Long:  `command to get cpu utilization metrics data`,

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
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			jsonResp, cloudwatchMetricResp, err := GetCpuUtilizationPanel(cmd, clientAuth)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
				return
			}
			if responseType == "frame" {
				fmt.Println(cloudwatchMetricResp)
			} else {
				// default case. it prints json
				fmt.Println(jsonResp)
			}
		}

	},
}

func GetCpuUtilizationPanel(cmd *cobra.Command, clientAuth *model.Auth) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {

	instanceID, _ := cmd.PersistentFlags().GetString("instanceID")
	namespace, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	// Parse start time if provided
	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
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
			err := cmd.Help()
			if err != nil {
				return "", nil, err
			}
			return "", nil, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}
	//if queryName == "cpu_utilization_panel" {
	currentUsage, err := GetCpuUtilizationMetricData(clientAuth, instanceID, namespace, startTime, endTime, "SampleCount")
	if err != nil {
		log.Println("Error in getting sample count: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["CurrentUsage"] = currentUsage
	// Get average usage
	averageUsage, err := GetCpuUtilizationMetricData(clientAuth, instanceID, namespace, startTime, endTime, "Average")
	if err != nil {
		log.Println("Error in getting average: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["AverageUsage"] = averageUsage
	// Get max usage
	maxUsage, err := GetCpuUtilizationMetricData(clientAuth, instanceID, namespace, startTime, endTime, "Maximum")
	if err != nil {
		log.Println("Error in getting maximum: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["MaxUsage"] = maxUsage
	jsonOutput := map[string]float64{
		"CurrentUsage": *currentUsage.MetricDataResults[0].Values[0],
		"AverageUsage": *averageUsage.MetricDataResults[0].Values[0],
		"MaxUsage":     *maxUsage.MetricDataResults[0].Values[0],
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil

}

func GetCpuUtilizationMetricData(clientAuth *model.Auth, instanceID, namespace string, startTime, endTime *time.Time, statistic string) (*cloudwatch.GetMetricDataOutput, error) {
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
								Value: aws.String(instanceID),
							},
						},
						MetricName: aws.String("CPUUtilization"),
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

func init() {
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("cloudElementId", "", "cloud element id")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("cloudElementApiUrl", "", "cloud element api")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")

	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("instanceID", "", "instance id")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2CpuUtilizationCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
