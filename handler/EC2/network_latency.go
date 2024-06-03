package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	comman_function "github.com/Appkube-awsx/awsx-getelementdetails/comman-function"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type Latency struct {
	Average float64 `json:"average"`
	Maximum float64 `json:"maximum"`
}

var AwsxEc2NetworkLatencyAcrossAllInstanceCmd = &cobra.Command{
	Use:   "total_network_utilization_panel",
	Short: "get network utilization metrics data for all instances",
	Long:  `command to get network utilization metrics data for all instances`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkLatencyAcrossAllInstancesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network utilization: ", err)
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

func GetNetworkLatencyAcrossAllInstancesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementType, _ := cmd.PersistentFlags().GetString("elementType")

	startTime, endTime, err := comman_function.ParseTimes(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing time: %v", err)
	}

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Get Inbound Traffic
	networkInOutput, err := GetTotalNetworkUtilizationMetricData(clientAuth, elementType, startTime, endTime, "Average", "NetworkIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting inbound traffic: ", err)
		return "", nil, err
	}
	//fmt.Println(networkInOutput)
	if len(networkInOutput.MetricDataResults) == 0 || len(networkInOutput.MetricDataResults[0].Values) == 0 {
		return "null", nil, nil
	}

	// Get Outbound Traffic
	networkOutOutput, err := GetTotalNetworkUtilizationMetricData(clientAuth, elementType, startTime, endTime, "Average", "NetworkOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting outbound traffic: ", err)
		return "", nil, err
	}
	if len(networkOutOutput.MetricDataResults) == 0 || len(networkOutOutput.MetricDataResults[0].Values) == 0 {
		return "null", nil, nil
	}

	networkLatency := CalculateNetworkLatency(networkInOutput, networkOutOutput)

	jsondata, err := json.Marshal(networkLatency)
	return string(jsondata), cloudwatchMetricData, err
}

func GetTotalNetworkUtilizationMetricData(clientAuth *model.Auth, elementType string, startTime, endTime *time.Time, statistic string, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data in namespace %s from %v to %v", elementType, startTime, endTime)
	elmType := "AWS/EC2"
	if elementType == "EC2" {
		elmType = "AWS/" + elementType
	}
	input := &cloudwatch.GetMetricDataInput{
		EndTime:   endTime,
		StartTime: startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
				},
			},
		},
	}
	if cloudWatchClient == nil {
		cloudWatchClient = awsclient.GetClient(*clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func CalculateNetworkLatency(networkIn, networkOut *cloudwatch.GetMetricDataOutput) *Latency {
	var totalNetworkIn, totalNetworkOut float64
	var maxNetworkLatency float64
	var dataPoints int

	for _, dataPoint := range networkIn.MetricDataResults[0].Values {
		totalNetworkIn += *dataPoint
		dataPoints++
	}

	for _, dataPoint := range networkOut.MetricDataResults[0].Values {
		totalNetworkOut += *dataPoint
	}

	if dataPoints == 0 {
		return &Latency{Average: 0, Maximum: 0}
	}

	averageNetworkIn := totalNetworkIn / float64(dataPoints)
	averageNetworkOut := totalNetworkOut / float64(dataPoints)
	totalnetwork := totalNetworkIn + totalNetworkOut
	// Placeholder formula for network latency
	networkLatency := (averageNetworkIn + averageNetworkOut) / 2

	// Calculate maximum latency
	if maxNetworkLatency < totalnetwork {
		maxNetworkLatency = totalnetwork
	}

	return &Latency{
		Average: networkLatency,
		Maximum: maxNetworkLatency,
	}
}

func init() {
	comman_function.InitAwsCmdFlags(AwsxEc2NetworkLatencyAcrossAllInstanceCmd)
}
