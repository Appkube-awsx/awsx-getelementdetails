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

type ContainerNetRxInBytes struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"RawData"`
}

var AwsxECSContainerNetRxInBytesCmd = &cobra.Command{
	Use:   "container_net_rxinbytes_panel",
	Short: "get container net received inbytes metrics data",
	Long:  `command to get container net received inbytes metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetECSContainerNetRxInBytesPanel(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting container net received inbytes metrics data: ", err)
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

func GetECSContainerNetRxInBytesPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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
	rawData, err := GetECSContainerNetRxInBytesMetricData(clientAuth, instanceId, elementType, startTime, endTime, "Sum", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Container_net_received_inbytes"] = rawData

	result := processECSContainerNetRxInbytesRawdata(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetECSContainerNetRxInBytesMetricData(clientAuth *model.Auth, instanceID, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	log.Printf("Getting metric data for instance %s in namespace %s from %v to %v", instanceID, elementType, startTime, endTime)

	elmType := "ECS/ContainerInsights"
	// if elementType == "ECS" {
	// 	elmType = "AWS/" + elementType
	// }
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
						MetricName: aws.String("NetworkRxBytes"),
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

func processECSContainerNetRxInbytesRawdata(result *cloudwatch.GetMetricDataOutput) ContainerNetRxInBytes {
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
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("query", "", "query")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxECSContainerNetRxInBytesCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
