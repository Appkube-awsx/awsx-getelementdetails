package Lambda

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

type LatencyGraph struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"latency_graph_panel"`
}

var AwsxLambdaLatencyGraphCmd = &cobra.Command{
	Use:   "Latency_graph_panel",
	Short: "get Latency count graph metrics data",
	Long:  `command to get Latency count graph metrics data`,

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
			jsonResp, cloudwatchMetricResp, err := GetLambdaLatencyGraphData(cmd, clientAuth, nil)
			if err != nil {
				log.Println("Error getting lambda Latency response data: ", err)
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

func GetLambdaLatencyGraphData(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
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

	// Fetch raw data
	LatencyCount, err := GetLambdaLatencyCountMetricValue(clientAuth, instanceId, elementType, startTime, endTime, "Average", cloudWatchClient)
	if err != nil {
		log.Println("Error in getting lambda latency count data: ", err)
		return "", nil, err
	}
	cloudwatchMetricData["Latency"] = LatencyCount

	result := ProcessLambdaLatencyRawData(LatencyCount)

	jsonString, err := json.Marshal(result)
	if err != nil {
		log.Println("Error in marshalling json in string: ", err)
		return "", nil, err
	}
	fmt.Println(jsonString)

	return string(jsonString), cloudwatchMetricData, nil
}

func GetLambdaLatencyCountMetricValue(clientAuth *model.Auth, instanceId string, elementType string, startTime, endTime *time.Time, statistic string, cloudWatchClient *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("latency"),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{

						Namespace:  aws.String("AWS/Lambda"),
						MetricName: aws.String("Duration"),
					},
					Period: aws.Int64(300),
					Stat:   aws.String(statistic),
				},
			},
		},
		StartTime: startTime,
		EndTime:   endTime,
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

func ProcessLambdaLatencyRawData(result *cloudwatch.GetMetricDataOutput) LatencyGraph {
	var rawData LatencyGraph
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
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("elementId", "", "element id")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("elementType", "", "element type")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("query", "", "query")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("cmdbApiUrl", "", "cmdb api")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("instanceId", "", "instance id")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("startTime", "", "start time")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("endTime", "", "end time")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("responseType", "", "response type. json/frame")
	AwsxLambdaLatencyGraphCmd.PersistentFlags().String("ApiName", "", "api name")
}
