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

type NetworkThroughputData struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NetworkThroughputData"`
}

var AwsxEc2NetworkThroughputCmd = &cobra.Command{
	Use:   "network_throughput_panel",
	Short: "get network throughput metrics data",
	Long:  `command to get network throughput metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkThroughputPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting network throughput data: ", err)
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

func GetNetworkThroughputPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
			return "", nil, err
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

	log.Printf("StartTime: %v, EndTime: %v", startTime, endTime)

	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{}

	// Fetch raw data for NetworkIn
	rawDataIn, err := GetNetworkThroughputMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", "NetworkIn", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data for NetworkIn: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkThroughputData"] = rawDataIn

	// Fetch raw data for NetworkOut
	rawDataOut, err := GetNetworkThroughputMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", "NetworkOut", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data for NetworkOut: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["NetworkThroughputData"] = rawDataOut

	// Combine the raw data for both NetworkIn and NetworkOut
	combinedRawData := combineNetworkThroughputRawData(rawDataIn, rawDataOut)

	result := processNetworkThroughputRawData(combinedRawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNetworkThroughputMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic, metricName string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
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

func combineNetworkThroughputRawData(rawDataIn, rawDataOut *cloudwatch.GetMetricDataOutput) map[string][]struct {
	Timestamp time.Time
	Value     float64
} {
	combinedRawData := make(map[string][]struct {
		Timestamp time.Time
		Value     float64
	})

	// Combine the timestamps and values for NetworkIn and NetworkOut
	for i, timestamp := range rawDataIn.MetricDataResults[0].Timestamps {
		combinedRawData["RawData"] = append(combinedRawData["RawData"], struct {
			Timestamp time.Time
			Value     float64
		}{
			Timestamp: *timestamp,
			Value:     *rawDataIn.MetricDataResults[0].Values[i] + *rawDataOut.MetricDataResults[0].Values[i],
		})
	}

	return combinedRawData
}

func processNetworkThroughputRawData(rawData map[string][]struct {
	Timestamp time.Time
	Value     float64
}) NetworkThroughputData {
	var processedData NetworkThroughputData

	// Assign the combined raw data to the processed data
	processedData.RawData = rawData["NetworkThroughputData"]

	return processedData
}

func init() {
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("query", "", "query")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEc2NetworkThroughputCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
