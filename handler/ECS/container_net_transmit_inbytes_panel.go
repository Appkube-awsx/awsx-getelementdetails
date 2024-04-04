package ECS

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"

	//"github.com/Appkube-awsx/awsx-common/cmdb"
	"github.com/Appkube-awsx/awsx-common/config"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type ContainerNetTxInBytes struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxECSContainerNetTxInBytesCmd = &cobra.Command{
	Use:   "container_net_txinbytes_panel",
	Short: "get container net transmit inbytes metrics data",
	Long:  `command to get container net transmit inbytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSContainerNetTxInBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting container net transmit inbytes metrics data: ", err)
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

func GetECSContainerNetTxInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId := "cluster-01-02-2024"

	if elementId != "" {
		log.Println("getting cloud-element data from cmdb")
		apiUrl := cmdbApiUrl
		if cmdbApiUrl == "" {
			log.Println("using default cmdb url")
			apiUrl = config.CmdbUrl
		}
		log.Println("cmdb url: " + apiUrl)
		// cmdbData, err := cmdb.GetCloudElementData(apiUrl, elementId)
		// if err != nil {
		// 	return "", nil, err
		// }
		// instanceId = cmdbData.InstanceId

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

	// Fetch raw data
	rawData, err := GetECSContainerNetTxInBytesMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Container_net_transmit_inbytes"] = rawData

	result := processECSContainerNetTxInbytesRawdata(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetECSContainerNetTxInBytesMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "ECS/ContainerInsights"

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
								Name:  aws.String("ClusterName"),
								Value: aws.String(instanceID),
							},
						},
						MetricName: aws.String("NetworkTxBytes"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"), // Assuming you want the sum of network received in bytes
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

func processECSContainerNetTxInbytesRawdata(result *cloudwatch.GetMetricDataOutput) ContainerNetRxInBytes {
	var rawData ContainerNetRxInBytes
	rawData.RawData = make([]struct {
		Timestamp time.Time
		Value     float64
	}, len(result.MetricDataResults[0].Timestamps))

	for i, timestamp := range result.MetricDataResults[0].Timestamps {
		rawData.RawData[i].Timestamp = *timestamp
		rawData.RawData[i].Value = *result.MetricDataResults[0].Values[i]
	}

	return rawData
}

func init() {
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSContainerNetTxInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
