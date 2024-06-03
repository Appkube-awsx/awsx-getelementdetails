package EC2

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var AwsxEc2CpuUtilizationAcrossAllInstanceCmd = &cobra.Command{
	Use:   "total_cpu_utilization_panel",
	Short: "get cpu utilization metrics data for all instances",
	Long:  `command to get cpu utilization metrics data for all instances`,

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
			jsonResp, cloudwatchMetricResp, err := GetCpuUtilizationPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting cpu utilization: ", err)
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

func GetCpuUtilizationAcrossAllInstancesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("CPUUtilization"),
						Namespace:  aws.String("AWS/" + elementType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String("Average"),
				},
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	cpuUtilization, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return "", nil, err
	}
	// fmt.Println("cpuUtilization", cpuUtilization)

	if len(cpuUtilization.MetricDataResults) > 0 && len(cpuUtilization.MetricDataResults[0].Values) > 0 {
		cloudwatchMetricData["CpuUtilization"] = cpuUtilization
	} else {
		log.Println("No data available for current Usage")
	}

	jsonOutput := make(map[string]float64)
	if len(cpuUtilization.MetricDataResults) > 0 && len(cpuUtilization.MetricDataResults[0].Values) > 0 {
		jsonOutput["CpuUtilization"] = *cpuUtilization.MetricDataResults[0].Values[0]
	}

	jsonString, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2CpuUtilizationAcrossAllInstanceCmd)
}
