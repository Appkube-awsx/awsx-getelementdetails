package EKS

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

type NetworkInOutResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"Network in and Network out"`
}

var AwsxEKSNetworkInOutCmd = &cobra.Command{
	Use:   "Network_in_out_panel",
	Short: "get Network in out graph metrics data",
	Long:  `command to get Network in out graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetNetworkInOutData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting Network in out data: ", err)
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

func GetNetworkInOutData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	elementId, _ := cmd.PersistentFlags().GetString("elementId")
	cmdbApiUrl, _ := cmd.PersistentFlags().GetString("cmdbApiUrl")
	instanceId, _ := cmd.PersistentFlags().GetString("instanceId")
	elementType, _ := cmd.PersistentFlags().GetString("elementType")
	startTimeStr, _ := cmd.PersistentFlags().GetString("startTime")
	endTimeStr, _ := cmd.PersistentFlags().GetString("endTime")

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

	rawData, err := GetNetworkInOutMetricData(clientAuth, instanceId, elementType, startTime, endTime, cloudWatchClient)
	if err != nil {
		log.Println("Error in getting raw data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Network in and Network out"] = rawData

	result := processNetworkInOutRawData(rawData)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}

	return string(jsonString), cloudwatchMetricData, nil
}

func GetNetworkInOutMetricData(clientAuth *model.Auth, instanceId, elementType string, startTime, endTime *time.Time, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	elmType := "ContainerInsights"
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
								Value: aws.String(instanceId),
							},
						},
						MetricName: aws.String("node_network_total_bytes"),
						Namespace:  aws.String(elmType),
					},
					Period: aws.Int64(60),
					Stat:   aws.String("Sum"), // Using Sum for network metrics
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

func processNetworkInOutRawData(result *cloudwatch.GetMetricDataOutput) NetworkInOutResult {
	var rawData NetworkInOutResult
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
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("query", "", "query")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("endTime", "", "endcl time")
	AwsxEKSNetworkInOutCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
}
