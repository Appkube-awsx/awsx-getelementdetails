package EC2

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type Networktrafficbound struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"network_traffic_panel"`
}

var AwsxEc2NetworkTrafficCmd = &cobra.Command{
	Use:   "network_traffic_panel",
	Short: "get network traffic metrics data",
	Long:  `command to get network traffic metrics data`,

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
			_, _, totalNetworkTraffic, err := GetNetworkTrafficPanel(cmd, clientAuth, nil)
			responseType, _ := cmd.PersistentFlags().GetString("responseType")
			if err != nil {
				log.Println("Error getting network in bytes metrics data: ", err)
				return
			}

			totalNetworkTrafficInMB := totalNetworkTraffic / (1024 * 1024)

			if responseType == "frame" {
				// Print only the total network traffic value
				fmt.Printf("NetworkTraffic: %.2f MB\n", totalNetworkTrafficInMB)
			} else {
				// Print the output in JSON format
				formattedTraffic := fmt.Sprintf("%.2f", totalNetworkTrafficInMB)
				fmt.Printf("{\"NetworkTraffic\": %.2f}\n", formattedTraffic)
			}
		}
	},
}

func GetNetworkTrafficPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, float64, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		if err != nil {
			return "", nil, 0, err
		}
		instanceId = cmdbData.InstanceId

	}

	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		parsedStartTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return "", nil, 0, err
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
			return "", nil, 0, err
		}
		endTime = &parsedEndTime
	} else {
		defaultEndTime := time.Now()
		endTime = &defaultEndTime
	}

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data for NetworkIn
	networkIn, err := GetNetworkTrafficMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", "NetworkIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NetworkIn data: ", err)
		return "", nil, 0, err
	}
	cloudwatchMetricData["NetworkIn"] = networkIn

	// Fetch raw data for NetworkOut
	networkOut, err := GetNetworkTrafficMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", "NetworkOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NetworkOut data: ", err)
		return "", nil, 0, err
	}
	cloudwatchMetricData["NetworkOut"] = networkOut

	// Fetch raw data for NetworkPacketsIn
	packetsIn, err := GetNetworkTrafficMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", "NetworkPacketsIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NetworkPacketsIn data: ", err)
		return "", nil, 0, err
	}
	cloudwatchMetricData["NetworkPacketsIn"] = packetsIn

	// Fetch raw data for NetworkPacketsOut
	packetsOut, err := GetNetworkTrafficMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", "NetworkPacketsOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting NetworkPacketsOut data: ", err)
		return "", nil, 0, err
	}
	cloudwatchMetricData["NetworkPacketsOut"] = packetsOut

	// Process raw data for all metrics
	totalNetworkTraffic := calculateTotalNetworkTraffic(networkIn, networkOut, packetsIn, packetsOut)

	// Prepare combined output
	output := map[string]float64{
		"NetworkTraffic": totalNetworkTraffic,
	}

	jsonString, err := json.Marshal(output)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, 0, err
	}

	return string(jsonString), cloudwatchMetricData, totalNetworkTraffic, nil
}

func calculateTotalNetworkTraffic(networkIn, networkOut, packetsIn, packetsOut *cloudwatch.GetMetricDataOutput) float64 {
	sum := calculateSum(networkIn) + calculateSum(networkOut) + calculateSum(packetsIn) + calculateSum(packetsOut)
	return sum
}

func calculateSum(data *cloudwatch.GetMetricDataOutput) float64 {
	var sum float64
	for _, value := range data.MetricDataResults[0].Values {
		sum += *value
	}
	return sum
}

func GetNetworkTrafficMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "AWS/EC2"

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
						MetricName: aws.String(metricName),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
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

func init() {
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkTrafficCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
